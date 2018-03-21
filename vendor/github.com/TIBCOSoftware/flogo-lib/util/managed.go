package util

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// Managed is an interface that is implemented by an object that needs to be
// managed via start/stop
type Managed interface {

	// Start starts the managed object
	Start() error

	// Stop stops the manged object
	Stop() error
}

// startManaged starts a "Managed" object
func startManaged(managed Managed) error {

	defer func() error {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}

			return err
		}

		return nil
	}()

	return managed.Start()
}

// stopManaged stops a "Managed" object
func stopManaged(managed Managed) error {

	defer func() error {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}

			return err
		}

		return nil
	}()

	return managed.Stop()
}

// StartManaged starts a Managed object, handles panics and logs details
func StartManaged(name string, managed Managed) error {

	logger.Debugf("%s: Starting...", name)
	err := managed.Start()

	if err != nil {
		logger.Errorf("%s: Error Starting", name)
		return err
	}

	logger.Debugf("%s: Started", name)
	return nil
}

// StopManaged stops a Managed object, handles panics and logs details
func StopManaged(name string, managed Managed) error {

	logger.Debugf("%s: Stopping...", name)

	err := stopManaged(managed)

	if err != nil {
		logger.Errorf("Error stopping '%s': %s", name, err.Error())
		return err
	}

	logger.Debugf("%s: Stopped", name)
	return nil
}
