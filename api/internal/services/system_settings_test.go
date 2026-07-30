package services

import (
	"net/url"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"testing"
)

func setLoginQrSettings(t *testing.T, show bool, mobileServerUrl string) {
	t.Helper()

	settings, err := repositories.NewSystemSettingsRepository(nil).GetSystemSettings()
	if err != nil {
		t.Fatalf("failed to load system settings: %v", err)
	}

	err = repositories.GetDB().
		Model(&models.SystemSettings{}).
		Where("id = ?", settings.ID).
		Updates(map[string]interface{}{"show_login_qr": show, "mobile_server_url": mobileServerUrl}).Error
	if err != nil {
		t.Fatalf("failed to set login qr system settings: %v", err)
	}
}

func TestBuildLoginQrUrl(t *testing.T) {
	tests := []struct {
		name      string
		serverUrl string
		want      string
	}{
		{
			name:      "encodes the server url into the fragment",
			serverUrl: "https://demo.receiptwrangler.io/api",
			want:      "https://receiptwrangler.io/app/setup#url=https%3A%2F%2Fdemo.receiptwrangler.io%2Fapi",
		},
		{
			name:      "encodes a lan http url with a port",
			serverUrl: "http://192.168.1.50:8081/api",
			want:      "https://receiptwrangler.io/app/setup#url=http%3A%2F%2F192.168.1.50%3A8081%2Fapi",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := BuildLoginQrUrl(test.serverUrl); got != test.want {
				t.Errorf("BuildLoginQrUrl(%q) = %q, want %q", test.serverUrl, got, test.want)
			}
		})
	}
}

// TestBuildLoginQrUrlRoundTrips proves the encoded fragment decodes back to the
// exact server url the mobile app must connect to.
func TestBuildLoginQrUrlRoundTrips(t *testing.T) {
	serverUrl := "https://demo.receiptwrangler.io/api"

	built := BuildLoginQrUrl(serverUrl)

	parsed, err := url.Parse(built)
	if err != nil {
		t.Fatalf("failed to parse built login qr url: %v", err)
	}

	decoded := parsed.Fragment
	if len(decoded) < len("url=") || decoded[:len("url=")] != "url=" {
		t.Fatalf("fragment %q does not start with url=", decoded)
	}

	got, err := url.QueryUnescape(decoded[len("url="):])
	if err != nil {
		t.Fatalf("failed to unescape fragment: %v", err)
	}

	if got != serverUrl {
		t.Errorf("round-tripped server url = %q, want %q", got, serverUrl)
	}
}

func TestGetFeatureConfigLoginQrUrl(t *testing.T) {
	service := NewSystemSettingsService(nil)

	tests := []struct {
		name            string
		show            bool
		mobileServerUrl string
		want            string
	}{
		{
			name:            "disabled with no url is empty",
			show:            false,
			mobileServerUrl: "",
			want:            "",
		},
		{
			name:            "disabled with a url is still empty",
			show:            false,
			mobileServerUrl: "https://demo.receiptwrangler.io/api",
			want:            "",
		},
		{
			name:            "enabled with an empty url is empty",
			show:            true,
			mobileServerUrl: "",
			want:            "",
		},
		{
			name:            "enabled with a url is the composed deep link",
			show:            true,
			mobileServerUrl: "https://demo.receiptwrangler.io/api",
			want:            BuildLoginQrUrl("https://demo.receiptwrangler.io/api"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer repositories.TruncateTestDb()
			setLoginQrSettings(t, test.show, test.mobileServerUrl)

			featureConfig, err := service.GetFeatureConfig()
			if err != nil {
				t.Fatalf("GetFeatureConfig() returned an error: %v", err)
			}

			if featureConfig.LoginQrUrl != test.want {
				t.Errorf("LoginQrUrl = %q, want %q", featureConfig.LoginQrUrl, test.want)
			}
		})
	}
}
