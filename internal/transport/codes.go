package transport

const (
	// Authorization
	Code_Unauthorized int = 10001
	Code_Forbidden    int = 10003

	// Request
	Code_Validation     int = 20000
	Code_RequestInvalid int = 29999

	// Resource
	Code_NotFound int = 30004

	// Internal
	Code_Unknown  int = 99998
	Code_Internal int = 99999
)
