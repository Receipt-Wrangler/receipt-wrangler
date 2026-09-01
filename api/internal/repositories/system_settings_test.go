package repositories

import (
	"strings"
	"sync"
	"testing"

	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"

	"gorm.io/gorm"
)

func intPtr(value int) *int {
	return &value
}

// buildSettingsCommand returns a command that passes validation, with both
// lifetimes omitted by default so callers can set only the one under test.
func buildSettingsCommand() commands.UpsertSystemSettingsCommand {
	queueConfigs := make([]commands.UpsertTaskQueueConfigurationCommand, 0)
	for _, config := range models.GetAllDefaultQueueConfigurations() {
		queueConfigs = append(queueConfigs, commands.UpsertTaskQueueConfigurationCommand{
			Name:     config.Name,
			Priority: 1,
		})
	}

	return commands.UpsertSystemSettingsCommand{
		CurrencyDisplay:              "$",
		CurrencySymbolPosition:       models.START,
		CurrencyThousandthsSeparator: models.COMMA,
		CurrencyDecimalSeparator:     models.DOT,
		TaskConcurrency:              1,
		EmailPollingInterval:         60,
		TaskQueueConfigurations:      queueConfigs,
	}
}

func seedLifetimes(t *testing.T, appHours int, mcpHours int) {
	t.Helper()

	if err := GetDB().Create(&models.SystemSettings{}).Error; err != nil {
		t.Fatalf("failed to create system settings: %v", err)
	}

	err := GetDB().Model(&models.SystemSettings{}).
		Where("id = ?", 1).
		Updates(map[string]interface{}{
			"refresh_token_valid_for_hours":     appHours,
			"mcp_refresh_token_valid_for_hours": mcpHours,
		}).Error
	if err != nil {
		t.Fatalf("failed to seed lifetimes: %v", err)
	}
}

// A lifetime the request omitted must be left out of the UPDATE statement
// altogether. Copying the stored value onto the row instead would still write
// the column, which is what lets a concurrent update clobber it. Asserted on the
// generated SQL because that is the precise property the concurrency safety
// rests on.
func TestUpdateSystemSettingsOmitsUnsentLifetimeColumns(t *testing.T) {
	defer TruncateTestDb()
	seedLifetimes(t, 24, 24)

	settings := models.SystemSettings{}
	settings.ID = 1

	tests := map[string]struct {
		command       commands.UpsertSystemSettingsCommand
		expectWritten []string
		expectSkipped []string
	}{
		"both omitted": {
			command:       buildSettingsCommand(),
			expectSkipped: []string{"`refresh_token_valid_for_hours`", "`mcp_refresh_token_valid_for_hours`"},
		},
		"app sent, mcp omitted": {
			command: func() commands.UpsertSystemSettingsCommand {
				cmd := buildSettingsCommand()
				cmd.RefreshTokenValidForHours = intPtr(720)
				return cmd
			}(),
			expectWritten: []string{"`refresh_token_valid_for_hours`"},
			expectSkipped: []string{"`mcp_refresh_token_valid_for_hours`"},
		},
		"mcp sent, app omitted": {
			command: func() commands.UpsertSystemSettingsCommand {
				cmd := buildSettingsCommand()
				cmd.McpRefreshTokenValidForHours = intPtr(6)
				return cmd
			}(),
			expectWritten: []string{"`mcp_refresh_token_valid_for_hours`"},
			expectSkipped: []string{"`refresh_token_valid_for_hours`"},
		},
		"both sent": {
			command: func() commands.UpsertSystemSettingsCommand {
				cmd := buildSettingsCommand()
				cmd.RefreshTokenValidForHours = intPtr(720)
				cmd.McpRefreshTokenValidForHours = intPtr(6)
				return cmd
			}(),
			expectWritten: []string{"`refresh_token_valid_for_hours`", "`mcp_refresh_token_valid_for_hours`"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			omitted := append([]string{"TaskQueueConfigurations"}, test.command.OmittedColumns()...)
			sql := GetDB().Session(&gorm.Session{DryRun: true}).
				Model(&settings).Select("*").Omit(omitted...).
				Where("id = ?", 1).Updates(&settings).Statement.SQL.String()

			for _, column := range test.expectWritten {
				if !strings.Contains(sql, column+"=") {
					t.Errorf("expected %s to be written, SQL: %s", column, sql)
				}
			}

			for _, column := range test.expectSkipped {
				if strings.Contains(sql, column+"=") {
					t.Errorf("expected %s to be skipped, SQL: %s", column, sql)
				}
			}
		})
	}
}

// Two admins saving at once, each setting one lifetime and omitting the other.
// Both explicit values must survive: with the columns skipped there is nothing
// for either write to clobber. Repeated because the losing interleaving is
// timing dependent -- a single pass could pass by luck against a broken
// implementation.
func TestUpdateSystemSettingsConcurrentLifetimeUpdatesBothSurvive(t *testing.T) {
	defer TruncateTestDb()

	const attempts = 15

	for attempt := 0; attempt < attempts; attempt++ {
		TruncateTestDb()
		seedLifetimes(t, 24, 24)

		appCommand := buildSettingsCommand()
		appCommand.RefreshTokenValidForHours = intPtr(720)

		mcpCommand := buildSettingsCommand()
		mcpCommand.McpRefreshTokenValidForHours = intPtr(6)

		var start sync.WaitGroup
		var done sync.WaitGroup
		start.Add(1)

		for _, command := range []commands.UpsertSystemSettingsCommand{appCommand, mcpCommand} {
			done.Add(1)
			go func(c commands.UpsertSystemSettingsCommand) {
				defer done.Done()
				start.Wait()
				NewSystemSettingsRepository(nil).UpdateSystemSettings(c)
			}(command)
		}

		start.Done()
		done.Wait()

		settings, err := NewSystemSettingsRepository(nil).GetSystemSettings()
		if err != nil {
			t.Fatalf("failed to read system settings: %v", err)
		}

		if settings.RefreshTokenValidForHours != 720 {
			t.Fatalf("attempt %d: app lifetime = %d, expected 720 (clobbered by the concurrent MCP update)", attempt, settings.RefreshTokenValidForHours)
		}

		if settings.McpRefreshTokenValidForHours != 6 {
			t.Fatalf("attempt %d: mcp lifetime = %d, expected 6 (clobbered by the concurrent app update)", attempt, settings.McpRefreshTokenValidForHours)
		}
	}
}
