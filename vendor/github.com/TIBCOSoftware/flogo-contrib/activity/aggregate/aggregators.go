package aggregate

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-contrib/activity/aggregate/window"
	"github.com/TIBCOSoftware/flogo-contrib/activity/aggregate/window/functions"
)

func NewTumblingWindow(function string, settings *window.Settings) (window.Window, error) {
	switch function {
	case "avg":
		return window.NewTumblingWindow(functions.AddSampleSum, functions.AggregateSingleAvg, settings), nil
	case "sum":
		return window.NewTumblingWindow(functions.AddSampleSum, functions.AggregateSingleNoopFunc, settings), nil
	case "min":
		return window.NewTumblingWindow(functions.AddSampleMin, functions.AggregateSingleNoopFunc, settings), nil
	case "max":
		return window.NewTumblingWindow(functions.AddSampleMax, functions.AggregateSingleNoopFunc, settings), nil
	case "count":
		return window.NewTumblingWindow(functions.AddSampleCount, functions.AggregateSingleNoopFunc, settings), nil
	case "accumulate":
		return window.NewTumblingWindow(functions.AddSampleAccum, functions.AggregateSingleNoopFunc, settings), nil
	default:
		return nil, fmt.Errorf("unsupported function: %s", function)
	}
}

// NewTumblingTimeWindow creates a new tumbling time window, all time windows are managed
// externally and are progressed using the NextBlock() method
func NewTumblingTimeWindow(function string, settings *window.Settings) (window.TimeWindow, error) {
	switch function {
	case "avg":
		return window.NewTumblingTimeWindow(functions.AddSampleSum, functions.AggregateSingleAvg, settings), nil
	case "sum":
		return window.NewTumblingTimeWindow(functions.AddSampleSum, functions.AggregateSingleNoopFunc, settings), nil
	case "min":
		return window.NewTumblingTimeWindow(functions.AddSampleMin, functions.AggregateSingleNoopFunc, settings), nil
	case "max":
		return window.NewTumblingTimeWindow(functions.AddSampleMax, functions.AggregateSingleNoopFunc, settings), nil
	case "count":
		return window.NewTumblingTimeWindow(functions.AddSampleCount, functions.AggregateSingleNoopFunc, settings), nil
	case "accumulate":
		return window.NewTumblingTimeWindow(functions.AddSampleAccum, functions.AggregateSingleNoopFunc, settings), nil
	default:
		return nil, fmt.Errorf("unsupported function: %s", function)
	}
}

func NewSlidingWindow(function string, settings *window.Settings) (window.Window, error) {
	switch function {
	case "avg":
		return window.NewSlidingWindow(functions.AggregateBlocksAvg, settings), nil
	case "sum":
		return window.NewSlidingWindow(functions.AggregateBlocksSum, settings), nil
	case "min":
		return window.NewSlidingWindow(functions.AggregateBlocksMin, settings), nil
	case "max":
		return window.NewSlidingWindow(functions.AggregateBlocksMax, settings), nil
	case "count":
		return window.NewSlidingWindow(functions.AggregateBlocksCount, settings), nil
	case "accumulate":
		return window.NewSlidingWindow(functions.AggregateBlocksAccumulate, settings), nil
	default:
		return nil, fmt.Errorf("unsupported function: %s", function)
	}
}

// NewSlidingTimeWindow creates a new sliding time window, all time windows are managed
// externally and are progressed using the NextBlock() method
func NewSlidingTimeWindow(function string, settings *window.Settings) (window.TimeWindow, error) {
	switch function {
	case "avg":
		return window.NewSlidingTimeWindow(functions.AddSampleSum, functions.AggregateBlocksAvg, settings), nil
	case "sum":
		return window.NewSlidingTimeWindow(functions.AddSampleSum, functions.AggregateBlocksSum, settings), nil
	case "min":
		return window.NewSlidingTimeWindow(functions.AddSampleMin, functions.AggregateBlocksMin, settings), nil
	case "max":
		return window.NewSlidingTimeWindow(functions.AddSampleMax, functions.AggregateBlocksMax, settings), nil
	case "count":
		return window.NewSlidingTimeWindow(functions.AddSampleCount, functions.AggregateBlocksSum, settings), nil
	default:
		return nil, fmt.Errorf("unsupported function: %s", function)
	}
}
