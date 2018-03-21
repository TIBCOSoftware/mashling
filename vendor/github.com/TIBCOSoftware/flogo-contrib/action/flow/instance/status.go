package instance

// Status is value that indicates the status of a Flow Instance
type Status int

const (
	// StatusNotStarted indicates that the FlowInstance has not started
	StatusNotStarted Status = 0

	// StatusActive indicates that the FlowInstance is active
	StatusActive Status = 100

	// StatusCompleted indicates that the FlowInstance has been completed
	StatusCompleted Status = 500

	// StatusCancelled indicates that the FlowInstance has been cancelled
	StatusCancelled Status = 600

	// StatusFailed indicates that the FlowInstance has failed
	StatusFailed Status = 700
)
