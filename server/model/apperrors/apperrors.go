package apperrors

const (
	ServerError    = "Something went wrong. Try again later"
	InvalidSession = "Provided session is invalid"
)

const (
	InvalidCredentials  = "Invalid email and password combination"
	DuplicateEmail      = "An account with that email already exists"
	InvalidResetToken   = "Invalid reset token"
	PasswordsDoNotMatch = "Passwords do not match"
)
