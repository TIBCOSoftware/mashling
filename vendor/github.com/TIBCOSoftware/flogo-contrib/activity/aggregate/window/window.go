package window

type Window interface {
	AddSample(sample interface{}) (bool, interface{})
}

type TimeWindow interface {
	Window
	NextBlock() (bool, interface{})
}
