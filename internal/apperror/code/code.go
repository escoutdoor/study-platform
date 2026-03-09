package code

type Code string

const (
	EmailAlreadyExists     Code = "EMAIL_ALREADY_EXISTS"
	StudentAlreadyEnrolled Code = "STUDENT_ALREADY_ENROLLED"
	TeacherAlreadyExists   Code = "TEACHER_ALREADY_EXISTS"

	PermissionDenied Code = "PERMISSION_DENIED"

	UserNotFound    Code = "USER_NOT_FOUND"
	StudentNotFound Code = "STUDENT_NOT_FOUND"
	TeacherNotFound Code = "TEACHER_NOT_FOUND"
	CourseNotFound  Code = "COURSE_NOT_FOUND"

	StudentNotEnrolled Code = "STUDENT_NOT_ENROLLED"

	InvalidJson      Code = "INVALID_JSON"
	ValidationFailed Code = "VALIDATION_FAILED"

	IncorrectCredentials Code = "INCORRECT_CREDENTIALS"
	JwtTokenExpired      Code = "JWT_TOKEN_EXPIRED"
	InvalidJwtToken      Code = "INVALID_JWT_TOKEN"
)
