package serviceerr

type ServiceErr struct {
	Message string
	Status  int
}

func (s ServiceErr) Error() string {
	return s.Message
}

// TODO: Add validationErr map[string]string field
// NewServiceErr creates a new error that can be used in services and handlers(controllers)
func NewServiceErr(message string, status int) ServiceErr {
	return ServiceErr{
		Message: message,
		Status:  status,
	}
}
