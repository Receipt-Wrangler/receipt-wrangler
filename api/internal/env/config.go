package env

import (
	"flag"
	"net/url"
	"os"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strconv"
	"strings"
)

var basePath string
var env string

func GetSecretKey() string {
	if len(os.Getenv(string(constants.SecretKey))) == 0 && env != "test" {
		logging.LogStd(logging.LOG_LEVEL_FATAL, constants.EmptySecretKeyError)
	}

	return os.Getenv(string(constants.SecretKey))
}

func GetDatabaseConfig() (structs.DatabaseConfig, error) {
	dbEngine := os.Getenv(string(constants.DbEngine))
	port := os.Getenv(string(constants.DbPort))
	portToUse := 0

	if dbEngine == "postgresql" || dbEngine == "mariadb" || dbEngine == "mysql" {
		parsedPort, err := utils.StringToInt(port)
		if err != nil {
			return structs.DatabaseConfig{}, err
		}

		portToUse = parsedPort
	}

	return structs.DatabaseConfig{
		User:     os.Getenv(string(constants.DbUser)),
		Password: os.Getenv(string(constants.DbPassword)),
		Name:     os.Getenv(string(constants.DbName)),
		Host:     os.Getenv(string(constants.DbHost)),
		Port:     portToUse,
		Engine:   os.Getenv(string(constants.DbEngine)),
		Filename: os.Getenv(string(constants.DbFileName)),
	}, nil
}

func GetBasePath() string {
	envBase := os.Getenv(string(constants.BasePath))
	if len(envBase) == 0 {
		return basePath
	}

	return envBase
}

func GetEncryptionKey() string {
	if len(os.Getenv(string(constants.EncryptionKey))) == 0 && env != "test" {
		logging.LogStd(logging.LOG_LEVEL_FATAL, constants.EmptyEncryptionKeyError)
	}

	return os.Getenv(string(constants.EncryptionKey))
}

func CheckRequiredEnvironmentVariables() {
	GetEncryptionKey()
	GetSecretKey()
}

func GetDeployEnv() string {
	return env
}

func GetChromiumPath() string {
	path := os.Getenv(string(constants.ChromiumBinaryPath))
	if len(path) == 0 {
		return "/usr/bin/chromium"
	}
	return path
}

// GetChromiumSandboxEnabled reports whether chromium should run with its
// process sandbox enabled. Defaults to false because the supported docker
// images run as root and the chromium sandbox refuses to start in that
// situation. Operators running the API as a non-root user can opt back in
// by setting CHROMIUM_SANDBOX to a truthy value (1, t, true, etc.). An
// unparseable value (e.g. yes/on/enabled) is logged at INFO and treated
// as the default (disabled) so misconfigurations are visible rather than
// silently ignored.
func GetChromiumSandboxEnabled() bool {
	return parseBoolEnv(constants.ChromiumSandbox, false)
}

// GetChromiumAllowExternalResources reports whether chromium should be
// allowed to load network resources (remote images, CSS, fonts) referenced
// from rendered HTML. Defaults to false: secure-by-default, no SSRF /
// tracking-pixel exposure, no remote-asset latency. Operators who need
// remote logos or product imagery in rendered receipts can opt in via
// CHROMIUM_ALLOW_EXTERNAL_RESOURCES=true.
func GetChromiumAllowExternalResources() bool {
	return parseBoolEnv(constants.ChromiumAllowExternalResources, false)
}

// parseBoolEnv reads a boolean env var with a default. Unparseable values
// are logged at INFO and treated as the default so misconfigurations
// surface in logs rather than being silently swapped to one side.
func parseBoolEnv(name constants.EnvironmentVariable, defaultValue bool) bool {
	raw := os.Getenv(string(name))
	if raw == "" {
		return defaultValue
	}
	enabled, err := strconv.ParseBool(raw)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_INFO,
			string(name)+" has unparseable value (use 1/0/true/false), using default: ",
			raw,
		)
		return defaultValue
	}
	return enabled
}

// GetMcpEnabled reports whether the MCP server and its OAuth 2.1
// authorization endpoints should be mounted. Defaults to false so existing
// deployments are unaffected until an operator explicitly opts in via
// MCP_ENABLED=true.
func GetMcpEnabled() bool {
	return parseBoolEnv(constants.McpEnabled, false)
}

// GetMcpPublicUrl returns the externally reachable origin (scheme + host,
// no trailing slash) used to build the OAuth issuer, resource, metadata, and
// redirect URLs advertised to MCP clients such as Claude. BASE_PATH is a
// filesystem path prefix, not an origin, so it cannot be reused here. In dev
// it defaults to the API's local address.
func GetMcpPublicUrl() string {
	const defaultUrl = "http://localhost:8081"

	raw := os.Getenv(string(constants.McpPublicUrl))
	if len(raw) == 0 {
		return defaultUrl
	}

	// MCP_PUBLIC_URL must be a bare origin (scheme + host). Any path, query, or
	// fragment would corrupt the issuer/metadata/redirect URLs derived from it,
	// so normalize to scheme://host and fall back to the default on anything
	// unparseable or missing a scheme/host.
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		logging.LogStd(logging.LOG_LEVEL_ERROR, string(constants.McpPublicUrl)+
			" must be an absolute origin like https://host (got "+raw+"); falling back to "+defaultUrl)
		return defaultUrl
	}

	if parsed.Path != "" && parsed.Path != "/" || parsed.RawQuery != "" || parsed.Fragment != "" {
		logging.LogStd(logging.LOG_LEVEL_INFO, string(constants.McpPublicUrl)+
			" should be a bare origin; ignoring path/query/fragment in "+raw)
	}

	return parsed.Scheme + "://" + parsed.Host
}

func SetConfigs() error {
	setEnv()
	setBasePath()

	return nil
}

func setEnv() {
	envFlag := flag.String("env", "dev", "set runtime environment")
	flag.Parse()

	env = *envFlag
	os.Setenv(string(constants.Env), env)
}

func setBasePath() {
	cwd, _ := os.Getwd()
	result := ""
	paths := strings.Split(cwd, "/")

	for i := 0; i < len(paths); i++ {
		result += "/" + paths[i]
		if paths[i] == "receipt-wrangler-api" {
			basePath = result
			return
		}
	}
}
