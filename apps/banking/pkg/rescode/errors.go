package rescode

import "net/http"

var (
	Failed             = New(codeFailed, http.StatusInternalServerError, msgFailed)
	NotFound           = New(codeNotFound, http.StatusNotFound, msgNotFound)
	ValidationFailed   = New(codeValidationFailed, http.StatusUnprocessableEntity, msgValidationFailed)
	UserDisabled       = New(codeUserDisabled, http.StatusForbidden, msgUserDisabled)
	UserVerifyRequired = New(codeUserVerifyRequired, http.StatusForbidden, msgUserVerifyRequired)
	EmailAlreadyExists = New(codeEmailAlreadyExists, http.StatusConflict, msgEmailAlreadyExists, R{
		"isExists": true,
	})
	VerificationExpired = New(codeVerificationExpired, http.StatusForbidden, msgVerificationExpired, R{
		"isExpired": true,
	})
	VerificationExceeded = New(codeVerificationExceeded, http.StatusForbidden, MsfVerificationExceeded, R{
		"isExceeded": true,
	})
	VerificationInvalid = New(codeVerificationInvalid, http.StatusForbidden, msgVerificationInvalid, R{
		"isInvalid": true,
	})
	InvalidRefreshOrAccessTokens = New(codeInvalidRefreshOrAccessTokens, http.StatusForbidden, msgInvalidRefreshOrAccessTokens)
	InvalidOrExpiredToken        = New(codeInvalidOrExpiredToken, http.StatusForbidden, msgInvalidOrExpiredToken, R{
		"isExpired": true,
	})
	InvalidAccess       = New(codeInvalidAccess, http.StatusForbidden, msgInvalidAccess)
	InvalidRefreshToken = New(codeInvalidRefreshToken, http.StatusForbidden, msgInvalidRefreshToken)
	RequiredVerifyToken = New(codeRequiredVerifyToken, http.StatusForbidden, msgRequiredVerifyToken)
	ExcludedVerifyToken = New(codeExcludedVerifyToken, http.StatusForbidden, msgExcludedVerifyToken)
	Unauthorized        = New(codeUnauthorized, http.StatusUnauthorized, msgUnauthorized)
	PermissionDenied    = New(codePermissionDenied, http.StatusForbidden, msgPermissionDenied)
	RecaptchaFailed     = New(codeRecaptchaFailed, http.StatusForbidden, msgRecaptchaFailed)
	RecaptchaRequired   = New(codeRecaptchaRequired, http.StatusForbidden, msgRecaptchaRequired)

	AccountBalanceInsufficient   = New(codeAccountBalanceInsufficient, http.StatusBadRequest, msgAccountBalanceInsufficient)
	AccountNotAvailable          = New(codeAccountNotAvailable, http.StatusBadRequest, msgAccountNotAvailable)
	ToAccountNotAvailable        = New(codeToAccountNotAvailable, http.StatusBadRequest, msgToAccountNotAvailable)
	AccountNotFound              = New(codeAccountNotFound, http.StatusNotFound, msgAccountNotFound)
	AccountTransferToSameAccount = New(codeAccountTransferToSameAccount, http.StatusBadRequest, msgAccountTransferToSameAccount)
	AccountCurrencyMismatch      = New(codeAccountCurrencyMismatch, http.StatusBadRequest, msgAccountCurrencyMismatch)
)
