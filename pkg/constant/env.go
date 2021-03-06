// @description 常量类 环境变量

package constant

const (
	// EnvKeySentinelLogDir ...
	EnvKeySentinelLogDir = "SENTINEL_LOG_DIR"
	// EnvKeySentinelAppName ...
	EnvKeySentinelAppName = "SENTINEL_APP_NAME"
)

const (
	// EnvAppName ...
	EnvAppName = "APP_NAME"
	// EnvDeployment ...
	EnvDeployment = "APP_DEPLOYMENT"

	EnvAppLogDir   = "APP_LOG_DIR"
	EnvAppMode     = "APP_MODE"
	EnvAppRegion   = "APP_REGION"
	EnvAppZone     = "APP_ZONE"
	EnvAppHost     = "APP_HOST"
	EnvAppInstance = "APP_INSTANCE" // application unique instance id.
	Env
)

const (
	// DefaultDeployment ...
	DefaultDeployment = ""
	// DefaultRegion ...
	DefaultRegion = ""
	// DefaultZone ...
	DefaultZone = ""
)

const (
	// KeyBalanceGroup ...
	KeyBalanceGroup = "__group"

	// DefaultBalanceGroup ...
	DefaultBalanceGroup = "default"
)

const (
	// OsWindows windows os
	OsWindows = "windows"
	// OsMac mac os
	OsMac = "darwin"
	// OsLinux linux os
	OsLinux = "linux"
	NL       = "\n"
)
