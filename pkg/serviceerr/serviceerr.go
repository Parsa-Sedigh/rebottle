package serviceerr

type ServiceErr struct {
	Message string
	Status  int
}

func (s ServiceErr) Error() string {
	return s.Message
}

// NewServiceErr creates a new error that can be used in services and handlers(controllers)
func NewServiceErr(message string, status int) ServiceErr {
	return ServiceErr{
		Message: message,
		Status:  status,
	}
}
