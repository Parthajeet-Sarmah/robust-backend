package custom_errors

import "net/http"

func ClientIdNonExistentError(cause error) *AppError {
	return New("CLIENT_ID_NON_EXISTENT", "This client ID does not exist", http.StatusForbidden, cause)
}
func ClientIdMismatchError(cause error) *AppError {
	return New("CLIENT_ID_MISMATCH", "Client ID does not match", http.StatusForbidden, cause)
}
func RedirectURIProtocolMismatch(cause error) *AppError {
	return New("REDIRECT_URI_PROTOCOL_MISMATCH", "Redirect URI does not have a valid HTTP protocol", http.StatusForbidden, cause)
}
func RedirectURIMismatchError(cause error) *AppError {
	return New("REDIRECT_URI_MISMATCH", "Redirect URI does not match", http.StatusForbidden, cause)
}
func ClientNotFoundError(cause error) *AppError {
	return New("CLIENT_NOT_FOUND", "This client was not found", http.StatusNotFound, cause)
}
