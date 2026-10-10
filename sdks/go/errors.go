package alldebrid

import (
	"errors"
	"fmt"
	"net/http"
)

// Error is an error returned by the API in a response with status "error".
// Compare its code with errors.Is:
//
//	if errors.Is(err, alldebrid.CodeMagnetInvalidID) { ... }
func (e *Error) Error() string {
	return fmt.Sprintf("alldebrid: %s: %s", e.Code, e.Message)
}

// Is reports whether the error has the given ErrorCode.
func (e *Error) Is(target error) bool {
	var code ErrorCode
	if errors.As(target, &code) {
		return ErrorCode(e.Code) == code
	}
	return false
}

// ErrorCode is the code of an API error. It implements error so it can be
// used as the target of errors.Is.
type ErrorCode string

func (c ErrorCode) Error() string { return string(c) }

// Codes the API is known to return. The list comes from the ErrorCode schema
// of openapi.yaml; the API may return others.
const (
	CodeGeneric                  ErrorCode = "GENERIC"
	CodeNotFound                 ErrorCode = "404"
	CodeMaintenance              ErrorCode = "MAINTENANCE"
	CodeDiscontinued             ErrorCode = "DISCONTINUED"
	CodeAuthMissingAPIKey        ErrorCode = "AUTH_MISSING_APIKEY"
	CodeAuthBadAPIKey            ErrorCode = "AUTH_BAD_APIKEY"
	CodeAuthBlocked              ErrorCode = "AUTH_BLOCKED"
	CodeAuthUserBanned           ErrorCode = "AUTH_USER_BANNED"
	CodeAlreadySent              ErrorCode = "ALREADY_SENT"
	CodeNoServer                 ErrorCode = "NO_SERVER"
	CodeLinkIsMissing            ErrorCode = "LINK_IS_MISSING"
	CodeBadLink                  ErrorCode = "BAD_LINK"
	CodeLinkHostNotSupported     ErrorCode = "LINK_HOST_NOT_SUPPORTED"
	CodeLinkDown                 ErrorCode = "LINK_DOWN"
	CodeLinkPassProtected        ErrorCode = "LINK_PASS_PROTECTED"
	CodeLinkHostUnavailable      ErrorCode = "LINK_HOST_UNAVAILABLE"
	CodeLinkTooManyDownloads     ErrorCode = "LINK_TOO_MANY_DOWNLOADS"
	CodeLinkHostFull             ErrorCode = "LINK_HOST_FULL"
	CodeLinkHostLimitReached     ErrorCode = "LINK_HOST_LIMIT_REACHED"
	CodeLinkError                ErrorCode = "LINK_ERROR"
	CodeLinkTemporaryUnavailable ErrorCode = "LINK_TEMPORARY_UNAVAILABLE"
	CodeLinkNotSupported         ErrorCode = "LINK_NOT_SUPPORTED"
	CodeRedirectorNotSupported   ErrorCode = "REDIRECTOR_NOT_SUPPORTED"
	CodeRedirectorError          ErrorCode = "REDIRECTOR_ERROR"
	CodeRedirectorRateLimited    ErrorCode = "REDIRECTOR_RATE_LIMITED"
	CodeStreamInvalidGenID       ErrorCode = "STREAM_INVALID_GEN_ID"
	CodeStreamInvalidStreamID    ErrorCode = "STREAM_INVALID_STREAM_ID"
	CodeDelayedInvalidID         ErrorCode = "DELAYED_INVALID_ID"
	CodeFreeTrialLimitReached    ErrorCode = "FREE_TRIAL_LIMIT_REACHED"
	CodeMustBePremium            ErrorCode = "MUST_BE_PREMIUM"
	CodeMagnetInvalidID          ErrorCode = "MAGNET_INVALID_ID"
	CodeMagnetInvalidURI         ErrorCode = "MAGNET_INVALID_URI"
	CodeMagnetInvalidFile        ErrorCode = "MAGNET_INVALID_FILE"
	CodeMagnetFileUploadFailed   ErrorCode = "MAGNET_FILE_UPLOAD_FAILED"
	CodeMagnetNoURI              ErrorCode = "MAGNET_NO_URI"
	CodeMagnetProcessing         ErrorCode = "MAGNET_PROCESSING"
	CodeMagnetTooManyActive      ErrorCode = "MAGNET_TOO_MANY_ACTIVE"
	CodeMagnetTooMany            ErrorCode = "MAGNET_TOO_MANY"
	CodeMagnetMustBePremium      ErrorCode = "MAGNET_MUST_BE_PREMIUM"
	CodeMagnetNoServer           ErrorCode = "MAGNET_NO_SERVER"
	CodeMagnetTooLarge           ErrorCode = "MAGNET_TOO_LARGE"
	CodeMagnetUploadFailed       ErrorCode = "MAGNET_UPLOAD_FAILED"
	CodeMagnetInternalError      ErrorCode = "MAGNET_INTERNAL_ERROR"
	CodeMagnetCantBootstrap      ErrorCode = "MAGNET_CANT_BOOTSTRAP"
	CodeMagnetMagnetTooBig       ErrorCode = "MAGNET_MAGNET_TOO_BIG"
	CodeMagnetTookTooLong        ErrorCode = "MAGNET_TOOK_TOO_LONG"
	CodeMagnetLinksRemoved       ErrorCode = "MAGNET_LINKS_REMOVED"
	CodeMagnetProcessingFailed   ErrorCode = "MAGNET_PROCESSING_FAILED"
	CodePINAlreadyAuthed         ErrorCode = "PIN_ALREADY_AUTHED"
	CodePINExpired               ErrorCode = "PIN_EXPIRED"
	CodePINInvalid               ErrorCode = "PIN_INVALID"
	CodeUserLinkMissing          ErrorCode = "USER_LINK_MISSING"
	CodeUserLinkInvalid          ErrorCode = "USER_LINK_INVALID"
	CodeMissingNotifEndpoint     ErrorCode = "MISSING_NOTIF_ENDPOINT"
	CodeVoucherDurationInvalid   ErrorCode = "VOUCHER_DURATION_INVALID"
	CodeVoucherNbInvalid         ErrorCode = "VOUCHER_NB_INVALID"
	CodeNoMoreVoucher            ErrorCode = "NO_MORE_VOUCHER"
	CodeInsufficientBalance      ErrorCode = "INSUFFICIENT_BALANCE"
	CodeDownloadFailed           ErrorCode = "DOWNLOAD_FAILED"
	CodeAccountInvalid           ErrorCode = "ACCOUNT_INVALID"
	CodeNoJSONParam              ErrorCode = "NO_JSON_PARAM"
	CodeJSONInvalid              ErrorCode = "JSON_INVALID"
	CodeFreedaysInvalidCountry   ErrorCode = "FREEDAYS_INVALID_COUNTRY"
	CodeFreedaysInvalidPhone     ErrorCode = "FREEDAYS_INVALID_PHONE"
	CodeFreedaysInvalidProvider  ErrorCode = "FREEDAYS_INVALID_PROVIDER"
	CodeFreedaysUsedPhone        ErrorCode = "FREEDAYS_USED_PHONE"
	CodeFreedaysAlreadyRequested ErrorCode = "FREEDAYS_ALREADY_REQUESTED"
	CodeFreedaysInvalidStatus    ErrorCode = "FREEDAYS_INVALID_STATUS"
	CodeFreedaysTooMuchRetries   ErrorCode = "FREEDAYS_TOO_MUCH_RETRIES"
	CodeFreedaysInvalidCode      ErrorCode = "FREEDAYS_INVALID_CODE"
	CodeFreedaysDenied           ErrorCode = "FREEDAYS_DENIED"
	CodeFreedaysErrorSending     ErrorCode = "FREEDAYS_ERROR_SENDING"
)

// HTTPError is returned when the API answers with an HTTP status other than
// 200, such as 429 when the rate limit is hit.
type HTTPError struct {
	StatusCode int
	// Err is the API error from the response body, when it has one.
	Err *Error
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("alldebrid: HTTP %d: %s: %s", e.StatusCode, e.Err.Code, e.Err.Message)
	}
	return fmt.Sprintf("alldebrid: HTTP %d %s", e.StatusCode, http.StatusText(e.StatusCode))
}

// Unwrap returns the API error, so errors.Is works on codes.
func (e *HTTPError) Unwrap() error {
	if e.Err == nil {
		return nil
	}
	return e.Err
}

// IsRateLimited reports whether err is a rate limiting answer (HTTP 429 or 503).
func IsRateLimited(err error) bool {
	var h *HTTPError
	return errors.As(err, &h) && (h.StatusCode == http.StatusTooManyRequests || h.StatusCode == http.StatusServiceUnavailable)
}
