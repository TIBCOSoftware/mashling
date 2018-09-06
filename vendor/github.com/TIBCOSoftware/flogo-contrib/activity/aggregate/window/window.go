package window

// Window is a basic sample window
type Window interface {
	// AddSample adds a sample to the window
	AddSample(sample interface{}) (bool, interface{})
}

// TimeWindow a time based sample window
type TimeWindow interface {
	Window

	// NextBlock tells the time window to advance
	NextBlock() (bool, interface{})
}
