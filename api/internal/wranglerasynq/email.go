package wranglerasynq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"os"
	"os/exec"
	"path/filepath"
	config "receipt-wrangler/api/internal/env"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
)

func StartEmailPolling() error {
	systemSettingsRepository := repositories.NewSystemSettingsRepository(nil)
	systemSettings, err := systemSettingsRepository.GetSystemSettings()
	if err != nil {
		return err
	}

	inspector, err := GetAsynqInspector()
	if err != nil {
		return nil
	}
	defer inspector.Close()

	_, err = inspector.DeleteAllScheduledTasks(string(models.EmailPollingQueue))
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_INFO, err.Error())
	}

	_, err = inspector.DeleteAllScheduledTasks(string(models.EmailReceiptImageCleanupQueue))
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_INFO, err.Error())
	}

	payload := EmailPollTaskPayload{
		PollAllGroups: true,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	pollTask := asynq.NewTask(EmailPoll, payloadBytes)
	_, err = RegisterTask(GetPollTimeString(systemSettings.EmailPollingInterval), pollTask, models.EmailPollingQueue, 0)
	if err != nil {
		return err
	}

	cleanUpTask := asynq.NewTask(EmailProcessImageCleanUp, nil)
	_, err = RegisterTask(GetPollTimeString(systemSettings.EmailPollingInterval*2), cleanUpTask, models.EmailReceiptImageCleanupQueue, 0)
	if err != nil {
		return err
	}

	return nil
}

func GetPollTimeString(pollingInterval int) string {
	return fmt.Sprintf("@every %ds", pollingInterval)
}

func CallClient(pollAllGroups bool, groupIds []string) error {
	groupSettingsRepository := repositories.NewGroupSettingsRepository(nil)
	var groupSettings []models.GroupSettings

	if pollAllGroups {
		allGroupSettings, err := groupSettingsRepository.GetAllGroupSettings("email_integration_enabled = ?", true)
		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
			return err
		}
		groupSettings = allGroupSettings
	} else {
		someGroupSettings, err := groupSettingsRepository.GetAllGroupSettings("email_integration_enabled = ? AND group_id IN ?", true, groupIds)
		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
			return err
		}
		groupSettings = someGroupSettings
	}

	if len(groupSettings) == 0 {
		logging.LogStd(logging.LOG_LEVEL_INFO, "No group settings enabled for email polling")
		return nil
	}

	err := pollEmailForGroupSettings(groupSettings)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		return err
	}
	return nil
}

func pollEmailForGroupSettings(groupSettings []models.GroupSettings) error {
	basePath := config.GetBasePath()
	groupSettingsWithPassword := make([]models.GroupSettingsWithSystemEmailPassword, len(groupSettings))

	// TODO: Could be more efficient by only decrypting the passwords once for each email
	for i := range groupSettings {
		cleartextPassword, err := utils.DecryptB64EncodedData(config.GetEncryptionKey(), groupSettings[i].SystemEmail.Password)
		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
			return err
		}

		var groupSettingWithPassword models.GroupSettingsWithSystemEmailPassword
		groupSettingWithPassword.BaseModel = groupSettings[i].BaseModel
		groupSettingWithPassword.GroupSettings = groupSettings[i]
		groupSettingWithPassword.SystemEmail = models.SystemEmailWithPassword{
			BaseModel:   groupSettings[i].SystemEmail.BaseModel,
			SystemEmail: groupSettings[i].SystemEmail,
			Password:    cleartextPassword,
		}

		groupSettingsWithPassword[i] = groupSettingWithPassword
	}

	bytesArr, err := json.Marshal(groupSettingsWithPassword)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		return err
	}

	var out bytes.Buffer
	cmd := exec.Command("python3", filepath.Join(basePath, "imap-client", "client.py"))
	cmd.Stdout = &out
	cmd.Stdin = bytes.NewReader(bytesArr)
	cmd.Env = os.Environ()

	err = cmd.Run()
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		return err
	}

	var result []structs.EmailMetadata

	err = json.Unmarshal(out.Bytes(), &result)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		return err
	}

	for _, m := range result {
		logging.LogStd(logging.LOG_LEVEL_INFO, fmt.Sprintf(
			"Email metadata captured: subject=%q from=%q attachments=%d bodyLen=%d bodyHtmlLen=%d groups=%v",
			m.Subject, m.FromEmail, len(m.Attachments), len(m.Body), len(m.BodyHtml), m.GroupSettingsIds,
		))
	}
	err = enqueueEmailProcessTasks(result)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		return err
	}

	return nil
}

func enqueueEmailProcessTasks(metadataList []structs.EmailMetadata) error {
	fileRepository := repositories.NewFileRepository(nil)

	// Build a lookup map of group settings to check EmailBodyProcessingEnabled
	groupSettingsLookup, err := buildGroupSettingsLookup(metadataList)
	if err != nil {
		return err
	}

	for _, metadata := range metadataList {
		// Normalize untrusted attachment filenames up front (GHSA-h3pr-mhcr-8phg):
		// they arrive from the raw MIME Content-Disposition header via the IMAP
		// client and must never escape the temp/data directories. Doing it here —
		// before any path is built or the task payload is populated — means every
		// downstream use (the temp path builders, the task payload, and the
		// persisted FileData.Name that drives the data-dir write) sees a safe
		// basename.
		for i := range metadata.Attachments {
			metadata.Attachments[i].Filename = utils.SanitizeFileName(metadata.Attachments[i].Filename)
		}

		for _, attachment := range metadata.Attachments {
			if attachment.Filename == "" {
				logging.LogStd(logging.LOG_LEVEL_ERROR, "Skipping email attachment with an unsafe or empty filename")
				continue
			}

			tempFilePath := buildTempEmailFilePath(attachment.Filename)
			imageForOcrPath := buildTempEmailOcrFilePath(attachment.Filename)

			fileBytes, err := utils.ReadFile(tempFilePath)
			if err != nil {
				return err
			}

			ocrBytes, err := fileRepository.GetBytesFromImageBytes(fileBytes)
			if err != nil {
				return err
			}

			err = utils.WriteFile(imageForOcrPath, ocrBytes)
			if err != nil {
				return err
			}

			for _, groupSettingsId := range metadata.GroupSettingsIds {
				taskMetadata := metadata
				if gs, ok := groupSettingsLookup[groupSettingsId]; !ok || !gs.EmailBodyProcessingEnabled {
					taskMetadata.Body = ""
					taskMetadata.BodyHtml = ""
				}

				payload := EmailProcessTaskPayload{
					GroupSettingsId: groupSettingsId,
					ImageForOcrPath: imageForOcrPath,
					TempFilePath:    tempFilePath,
					Metadata:        taskMetadata,
					Attachment:      attachment,
					RenderBodyPdf:   shouldRenderEmailBodyPdfForGroup(metadata, groupSettingsId, groupSettingsLookup),
				}
				payloadBytes, err := json.Marshal(payload)
				if err != nil {
					return err
				}

				task := asynq.NewTask(EmailProcess, payloadBytes)
				_, err = EnqueueTask(task, models.EmailReceiptProcessingQueue)
				if err != nil {
					return err
				}
			}
		}

		// Handle body-only emails (no attachments but has any body content)
		if len(metadata.Attachments) == 0 && (len(metadata.Body) > 0 || len(metadata.BodyHtml) > 0) {
			for _, groupSettingsId := range metadata.GroupSettingsIds {
				// Skip body-only emails for groups that don't have body processing enabled
				if gs, ok := groupSettingsLookup[groupSettingsId]; !ok || !gs.EmailBodyProcessingEnabled {
					continue
				}

				payload := EmailProcessTaskPayload{
					GroupSettingsId: groupSettingsId,
					Metadata:        metadata,
					RenderBodyPdf:   shouldRenderEmailBodyPdfForGroup(metadata, groupSettingsId, groupSettingsLookup),
				}
				payloadBytes, err := json.Marshal(payload)
				if err != nil {
					return err
				}

				task := asynq.NewTask(EmailProcess, payloadBytes)
				_, err = EnqueueTask(task, models.EmailReceiptProcessingQueue)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// shouldRenderEmailBodyPdfForGroup returns true when the email has an HTML
// body with meaningful text content AND the named group has email body
// processing enabled. Extracted so the gating logic is testable without
// launching chromium and so a group missing from the lookup defaults to
// "do not render" (defensive).
//
// We also require the stripped-text Body to be non-whitespace: many email
// clients (Gmail's web composer in particular) auto-wrap every message in
// a boilerplate HTML shell (<html><body><div style="..."></div></body>),
// so BodyHtml alone is not a good proxy for "there's a body worth
// rendering." Rendering those wrappers produces a near-blank PDF that
// just adds a second, useless image to the LLM call.
func shouldRenderEmailBodyPdfForGroup(
	metadata structs.EmailMetadata,
	groupSettingsId uint,
	groupSettingsLookup map[uint]models.GroupSettings,
) bool {
	if len(metadata.BodyHtml) == 0 {
		return false
	}
	if len(strings.TrimSpace(metadata.Body)) == 0 {
		return false
	}
	gs, ok := groupSettingsLookup[groupSettingsId]
	if !ok {
		return false
	}
	return gs.EmailBodyProcessingEnabled
}

func buildGroupSettingsLookup(metadataList []structs.EmailMetadata) (map[uint]models.GroupSettings, error) {
	uniqueIds := make(map[uint]struct{})
	for _, metadata := range metadataList {
		for _, id := range metadata.GroupSettingsIds {
			uniqueIds[id] = struct{}{}
		}
	}

	ids := make([]uint, 0, len(uniqueIds))
	for id := range uniqueIds {
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return make(map[uint]models.GroupSettings), nil
	}

	groupSettingsRepository := repositories.NewGroupSettingsRepository(nil)
	allSettings, err := groupSettingsRepository.GetAllGroupSettings("id IN ?", ids)
	if err != nil {
		return nil, err
	}

	lookup := make(map[uint]models.GroupSettings, len(allSettings))
	for _, gs := range allSettings {
		lookup[gs.ID] = gs
	}

	return lookup, nil
}

// buildTempEmailFilePath builds the temp-dir path for an email attachment. The
// caller must pass a filename already reduced to a safe basename via
// utils.SanitizeFileName (done once in enqueueEmailProcessTasks); this join does
// not sanitize, so an unsanitized name could escape the temp dir (GHSA-h3pr-mhcr-8phg).
func buildTempEmailFilePath(attachmentFileName string) string {
	fileRepository := repositories.NewFileRepository(nil)
	return filepath.Join(fileRepository.GetTempDirectoryPath(), attachmentFileName)
}

// buildTempEmailOcrFilePath builds the temp-dir path for the OCR-preprocessed
// copy of an attachment. Same contract as buildTempEmailFilePath: the caller
// passes a filename already sanitized to a safe basename.
func buildTempEmailOcrFilePath(attachmentFileName string) string {
	fileRepository := repositories.NewFileRepository(nil)
	return filepath.Join(fileRepository.GetTempDirectoryPath(), "image-"+attachmentFileName)
}
