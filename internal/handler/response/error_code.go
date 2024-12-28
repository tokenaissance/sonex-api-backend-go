package response

const (
	ServerErrorCode        = 10101
	TooManyRequestsCode    = 10102
	ParamBindErrorCode     = 10103
	AuthorizationErrorCode = 10104
	CacheSetErrorCode      = 10105
	CacheGetErrorCode      = 10106
	CacheDelErrorCode      = 10107
	CacheNotExistCode      = 10108
	ResubmitErrorCode      = 10109
	RBACErrorCode          = 10110
	RedisConnectErrorCode  = 10111
	DBConnectErrorCode     = 10112
	DBExecErrorCode        = 10113
	GoVersionErrorCode     = 10114
	SocketConnectErrorCode = 10115
	SocketSendErrorCode    = 10116
	ChainClientErrorCode   = 10117

	CronCreateErrorCode  = 10201
	CronUpdateErrorCode  = 10202
	CronListErrorCode    = 10203
	CronDetailErrorCode  = 10204
	CronExecuteErrorCode = 10205

	// user 2
	AddressErrorCode          = 20101
	SignatureErrorCode        = 20102
	CreateUserErrorCode       = 20103
	UserNotExistsCode         = 20104
	UserNotLoginCode          = 20105
	UpdateRefererErrorCode    = 20106
	JWTMissingHeaderErrorCode = 20107
	JWTInvalidHeaderErrorCode = 20108
	JWTTokenExpiredErrorCode  = 20109

	// liquidity 3

	// farm 4
	FarmErrorCode      = 40101
	FarmParamErrorCode = 40102

	// xpool 5
	XpoolErrorCode      = 50101
	XpoolParamErrorCode = 50102
)

func Text(code int) string {
	return enUSText[code]
}
