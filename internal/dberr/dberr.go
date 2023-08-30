package dberr

// DBError satisfies Error interface but allows constant values for direct comparison.
type DBError string

// Error is required by error interface.
func (s DBError) Error() string {
	return string(s)
}

const (
	// ErrCredentials means credentials like email or captcha must be validated.
	ErrCredentials = DBError("credentials")
	// ErrUserNotFound means the user was not found.
	ErrUserNotFound = DBError("user not found")
	// ErrTopicNotFound means the topic was not found.
	ErrTopicNotFound = DBError("topic not found")
	// ErrNotFound means the object other than user or topic was not found.
	ErrNotFound = DBError("not found")
	// ErrPermissionDenied means the operation is not permitted.
	ErrPermissionDenied = DBError("denied")
)
