package utils

type UserNotLoggedInError struct{}
type UserScopeDeniedError struct{}
type CouldNotConnectToDatabaseError struct{}

type RedisCouldNotCreateClient struct{}
type RedisGetHashError struct{}
type RedisSetHasError struct{}
type RedisGetHashNoResourceFoundError struct{}

type ClientIdNonExistentError struct{}
type RedirectURIMismatchError struct{}

type UserNotFoundError struct{}
type ClientNotFoundError struct{}

type UnknownError struct{}

func (e *UnknownError) Error() string {
	return "Some unknown error occured"
}

func (e *UserNotFoundError) Error() string {
	return "The requested user was not found! Please check your credentials"
}

func (e *ClientNotFoundError) Error() string {
	return "The requested client was not found!"
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

func (e *RedisSetHasError) Error() string {
	return "Error while setting resource to Redis hash!"
}

func (e *ClientIdNonExistentError) Error() string {
	return "The client id provided does not exist!"
}

func (e *RedirectURIMismatchError) Error() string {
	return "The redirect URI provided is invalid!"
}
