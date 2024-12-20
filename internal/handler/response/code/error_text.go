package code

var enUSText = map[int]string{
	ServerError:        "Internal server error",
	TooManyRequests:    "Too many requests",
	ParamBindError:     "Parameter error",
	AuthorizationError: "Authorization error",
	CacheSetError:      "Failed to set cache",
	CacheGetError:      "Failed to get cache",
	CacheDelError:      "Failed to del cache",
	CacheNotExist:      "Cache does not exist",
	ResubmitError:      "Please do not submit repeatedly",
	RBACError:          "No access",
	RedisConnectError:  "Failed to connection Redis",
	DBConnectError:     "Failed to connection DB",
	DBExecError:        "SQL execution failed",
	GoVersionError:     "Go Version mismatch",
	SocketConnectError: "Socket not connected",
	SocketSendError:    "Socket message sending failed",

	CronCreateError:  "Failed to create cron",
	CronUpdateError:  "Failed to update menu",
	CronListError:    "Failed to get cron list",
	CronDetailError:  "Failed to get cron detail",
	CronExecuteError: "Failed to execute cron",

	// user
	AddressError:    "Illegal address",
	SignatureError:  "Signature error",
	CreateUserError: "Failed to create user",
	UserNotExists:   "User does not exist",
	UserNotLogin:    "User not logged in",

	JWTMissingHeaderError: "Missing Authorization header",
	JWTInvalidHeaderError: "Invalid Authorization header",
	JWTTokenExpiredError:  "JWT expired",

	// liquidity

	// farm

	// x-pool
	XpoolError: "xpool error",
}
