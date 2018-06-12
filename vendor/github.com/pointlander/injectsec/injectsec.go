// Copyright 2018 The InjectSec Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package injectsec

import (
	"bytes"

	"github.com/pointlander/injectsec/gru"
)

// DetectorMaker makes SQL injection attack detectors
type DetectorMaker struct {
	*gru.DetectorMaker
}

// NewDetectorMaker creates a new detector maker
func NewDetectorMaker() *DetectorMaker {
	maker := gru.NewDetectorMaker()
	weights, err := ReadFile("weights.w")
	if err != nil {
		panic(err)
	}
	err = maker.Read(bytes.NewBuffer(weights))
	if err != nil {
		panic(err)
	}
	return &DetectorMaker{
		DetectorMaker: maker,
	}
}
