package rescode

const (
	codeValidationFailed uint64 = 1000
	codeFailed           uint64 = 1001
	codeNotFound         uint64 = 1002

	codeUserDisabled       uint64 = 2000
	codeUserVerifyRequired uint64 = 2001
	codeEmailAlreadyExists uint64 = 2002

	codeVerificationExpired          uint64 = 3000
	codeVerificationExceeded         uint64 = 3001
	codeVerificationInvalid          uint64 = 3002
	codeInvalidRefreshOrAccessTokens uint64 = 3050
	codeInvalidOrExpiredToken        uint64 = 3051
	codeInvalidAccess                uint64 = 3052
	codeInvalidRefreshToken          uint64 = 3053
	codeRequiredVerifyToken          uint64 = 3054
	codeExcludedVerifyToken          uint64 = 3055

	codeUnauthorized      uint64 = 3100
	codePermissionDenied  uint64 = 3101
	codeRecaptchaFailed   uint64 = 3102
	codeRecaptchaRequired uint64 = 3103

	codeAccountBalanceInsufficient   uint64 = 4000
	codeAccountNotAvailable          uint64 = 4001
	codeToAccountNotAvailable        uint64 = 4001
	codeAccountNotFound              uint64 = 4002
	codeAccountTransferToSameAccount uint64 = 4003
	codeAccountCurrencyMismatch      uint64 = 4004
)
