package response

var enUSText = map[int]string{
	ServerErrorCode:        "Internal server error",
	TooManyRequestsCode:    "Too many requests",
	ParamBindErrorCode:     "Parameter error",
	AuthorizationErrorCode: "Authorization error",
	CacheSetErrorCode:      "Failed to set cache",
	CacheGetErrorCode:      "Failed to get cache",
	CacheDelErrorCode:      "Failed to del cache",
	CacheNotExistCode:      "Cache does not exist",
	ResubmitErrorCode:      "Please do not submit repeatedly",
	RBACErrorCode:          "No access",
	RedisConnectErrorCode:  "Failed to connection Redis",
	DBConnectErrorCode:     "Failed to connection DB",
	DBExecErrorCode:        "SQL execution failed",
	GoVersionErrorCode:     "Go Version mismatch",
	SocketConnectErrorCode: "Socket not connected",
	SocketSendErrorCode:    "Socket message sending failed",

	CronCreateErrorCode:  "Failed to create cron",
	CronUpdateErrorCode:  "Failed to update menu",
	CronListErrorCode:    "Failed to get cron list",
	CronDetailErrorCode:  "Failed to get cron detail",
	CronExecuteErrorCode: "Failed to execute cron",

	// user
	AddressErrorCode:    "Illegal address",
	SignatureErrorCode:  "Signature error",
	CreateUserErrorCode: "Failed to create user",
	UserNotExistsCode:   "User does not exist",
	UserNotLoginCode:    "User not logged in",

	JWTMissingHeaderErrorCode: "Missing Authorization header",
	JWTInvalidHeaderErrorCode: "Invalid Authorization header",
	JWTTokenExpiredErrorCode:  "JWT expired",

	// liquidity

	// farm
	FarmErrorCode:      "farm error",
	FarmParamErrorCode: "farm parameter error",

	// x-pool
	XpoolErrorCode:      "xpool error",
	XpoolParamErrorCode: "xpool parameter error",
}
