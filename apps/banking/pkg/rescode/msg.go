package rescode

const (
	msgValidationFailed string = "validation_failed"
	msgFailed           string = "failed"
	msgNotFound         string = "not_found"

	msgUserDisabled       string = "user_disabled"
	msgUserVerifyRequired string = "user_verify_required"
	msgEmailAlreadyExists string = "email_already_exists"

	msgVerificationExpired          string = "verification_expired"
	MsfVerificationExceeded         string = "verification_exceeded"
	msgVerificationInvalid          string = "verification_invalid"
	msgInvalidRefreshOrAccessTokens string = "invalid_refresh_or_access_tokens"
	msgInvalidOrExpiredToken        string = "invalid_or_expired_token"
	msgInvalidAccess                string = "invalid_access"
	msgInvalidRefreshToken          string = "invalid_refresh_token"
	msgRequiredVerifyToken          string = "required_verify_token"
	msgExcludedVerifyToken          string = "excluded_verify_token"

	msgUnauthorized      string = "unauthorized"
	msgPermissionDenied  string = "permission_denied"
	msgRecaptchaFailed   string = "recaptcha_failed"
	msgRecaptchaRequired string = "recaptcha_required"

	msgAccountBalanceInsufficient   string = "account_balance_insufficient"
	msgAccountNotAvailable          string = "account_not_available"
	msgToAccountNotAvailable        string = "to_account_not_available"
	msgAccountNotFound              string = "account_not_found"
	msgAccountTransferToSameAccount string = "account_transfer_to_same_account"
	msgAccountCurrencyMismatch      string = "account_currency_mismatch"
)
