package code

const (
	ServerError        = 10101
	TooManyRequests    = 10102
	ParamBindError     = 10103
	AuthorizationError = 10104
	CacheSetError      = 10105
	CacheGetError      = 10106
	CacheDelError      = 10107
	CacheNotExist      = 10108
	ResubmitError      = 10109
	RBACError          = 10110
	RedisConnectError  = 10111
	DBConnectError     = 10112
	DBExecError        = 10113
	GoVersionError     = 10114
	SocketConnectError = 10115
	SocketSendError    = 10116

	CronCreateError  = 10201
	CronUpdateError  = 10202
	CronListError    = 10203
	CronDetailError  = 10204
	CronExecuteError = 10205

	// user 2
	AddressError          = 20101
	SignatureError        = 20102
	CreateUserError       = 20103
	UserNotExists         = 20104
	UserNotLogin          = 20105
	UpdateRefererError    = 20106
	JWTMissingHeaderError = 20107
	JWTInvalidHeaderError = 20108
	JWTTokenExpiredError  = 20109

	// liquidity 3

	// farm 4

	// xpool 5
	XpoolError = 50101
)

func Text(code int) string {
	return enUSText[code]
}
