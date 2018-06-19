// Copyright 2018 The InjectSec Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package injectsec

import (
	"bytes"
	"io"

	"github.com/pointlander/injectsec/gru"
)

// DetectorMaker makes SQL injection attack detectors
type DetectorMaker struct {
	*gru.DetectorMaker
}

// NewDetectorMakerWithWeights creates a new detector maker using weights
func NewDetectorMakerWithWeights(weights io.Reader) (*DetectorMaker, error) {
	maker := gru.NewDetectorMaker()
	err := maker.Read(weights)
	if err != nil {
		return nil, err
	}

	return &DetectorMaker{
		DetectorMaker: maker,
	}, nil
}

// NewDetectorMaker creates a new detector maker
func NewDetectorMaker() (*DetectorMaker, error) {
	weights, err := ReadFile("weights.w")
	if err != nil {
		return nil, err
	}
	return NewDetectorMakerWithWeights(bytes.NewBuffer(weights))
}
