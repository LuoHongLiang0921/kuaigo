package constant

const (
	TabbyMode    = "TABBY_MODE"
	TabbyModeDev = "dev"
)

const (
	HeaderFieldRv        = "rv"
	HeaderFieldRt        = "rt"
	HeaderFieldSt        = "st"
	HeaderFieldPk        = "pk"
	HeaderFieldSi        = "si"
	HeaderFieldTi        = "ti"
	HeaderFieldSk        = "sk"
	HeaderFieldOv        = "ov"
	HeaderFieldCs        = "cs"
	HeaderFieldUserAgent = "user-agent"
	HeaderFieldAi        = "ai"
	OvMin                = 23
	//OsIos 设备系统类型ios
	OsIos = "ios"
	//OsAndroid 设备系统类型android
	OsAndroid = "android"
	//OsWeb 设备系统类型web
	OsWeb = "web"
	//OsH5 设备系统类型h5
	OsH5 = "h5"

	HeaderUserAgent        = "user-agent"
	HeaderUserAgentAndroid = "adr"
	HeaderUserAgentIpad    = "ipad"
	HeaderUserAgentIpod    = "ipod"
	HeaderUserAgentIphone  = "iphone"
	OsPc                   = "pc"
	DmH5                   = "H5"
	DmPc                   = "pc"

	SignExpireTime        int64 = 3 * 60 * 1000
	HeaderFieldAdminId          = "adminId"
	HeaderFieldAdminAppId       = "adminAppId"
	HeaderFieldT                = "t"

	//上下文传递的实体key
	KeyHeader      = "tabbyHeader"
	KeyAdminHeader = "adminHeader"
	KeyPParam      = "pparam"
	KeySecret      = "secret"
	KeyAppId       = "appId"

	KeyP = "p"
	KeyS = "s"
	// response 业务 code
	RespCode = "code"

	KeyCurrUser       = "currUser"
	KeyFormData       = "multipart/form-data"
	KeyFormUrlencoded = "application/x-www-form-urlencoded"

	AdminHeaderFieldAid = "aid"
	AdminHeaderFieldAln = "aln"
	AdminHeaderFieldAac = "aac"
)

const (
	CommonKey = "common"
	ParamsKey = "params"
	TypeKey   = "type"
)

const (
	AbId    = "abId"
	AbParam = "abParam"
	FUDid   = "fudid"
)

// apollo 配置文件
const (
	ConfigAppKey = "app"
)

// 游客登录过期时间设置为1个月，单位：分钟
const GuestLoginExpireTime = 43200

const GuestUserKeyPre = "guest:user:"
