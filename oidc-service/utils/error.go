package utils

type UserAlreadyExistsError struct{}
type UserNotLoggedInError struct{}
type UserScopeDeniedError struct{}

type CouldNotConnectToDatabaseError struct{}
type CouldNotFetchAuthCode struct{}

type RedisCouldNotCreateClient struct{}
type RedisGetHashError struct{}
type RedisSetHasError struct{}
type RedisGetHashNoResourceFoundError struct{}

type ClientIdNonExistentError struct{}

type ClientIdMismatchError struct{}
type RedirectURIMismatchError struct{}

type UserNotFoundError struct{}
type ClientNotFoundError struct{}

type ExpiredAuthCodeError struct{}
type ExpiredRefreshTokenError struct{}
type AuthCodeUsedUpdateError struct{}
type CodeChallengeDoesNotMatchError struct{}
type UnknownError struct{}
type InvalidGrantType struct{}
type TokenParsingError struct{}
type NoAccessTokenFoundError struct{}

type RefreshTokenNotFoundError struct {
	Status int
	Msg    string
}

type MalformedRequest struct {
	Status int
	Msg    string
}

func (e *MalformedRequest) Error() string {
	return e.Msg
}

func (e *UnknownError) Error() string {
	return "Some unknown error occured"
}

func (e *UserAlreadyExistsError) Error() string {
	return "A user already exists with this credential"
}

func (e *UserNotFoundError) Error() string {
	return "The requested user was not found! Please check your credentials"
}

func (e *ClientNotFoundError) Error() string {
	return "The requested client was not found!"
}

func (e *ExpiredAuthCodeError) Error() string {
	return "The auth code has expired!"
}

func (e *ExpiredRefreshTokenError) Error() string {
	return "The refresh token has expired!"
}

func (e *AuthCodeUsedUpdateError) Error() string {
	return "The authorization code could not be updated!"
}

func (e *CouldNotFetchAuthCode) Error() string {
	return "The authorization code could not be fetched!"
}

func (e *InvalidGrantType) Error() string {
	return "The grant type provided is not valid!"
}

func (e *TokenParsingError) Error() string {
	return "The token could not be parsed!"
}

func (e *NoAccessTokenFoundError) Error() string {
	return "The token could not be parsed!"
}

func (e *RefreshTokenNotFoundError) Error() string {
	return e.Msg
}

func (e *CouldNotConnectToDatabaseError) Error() string {
	return "Could not connect to Database"
}

func (e *UserNotLoggedInError) Error() string {
	return "The user is not logged in!"
}

func (e *UserScopeDeniedError) Error() string {
	return "The user denied permission for the scope!"
}

func (e *RedisCouldNotCreateClient) Error() string {
	return "Could not create Redis client"
}

func (e *RedisGetHashError) Error() string {
	return "Error while getting resource from Redis hash!"
}

func (e *RedisGetHashNoResourceFoundError) Error() string {
	return "No resource returned while getting Redis hash!"
}

func (e *CodeChallengeDoesNotMatchError) Error() string {
	return "The code challenge does not match!"
}

func (e *RedisSetHasError) Error() string {
	return "Error while setting resource to Redis hash!"
}

func (e *ClientIdNonExistentError) Error() string {
	return "The client id provided does not exist!"
}

func (e *ClientIdMismatchError) Error() string {
	return "The client id does not match"
}

func (e *RedirectURIMismatchError) Error() string {
	return "The redirect URI provided is invalid!"
}
