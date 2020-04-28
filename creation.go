package piper

import (
	"container/list"
	"reflect"
)

// Clone is used to create a Piper from a given channel where the output
// channel is of the exact same type as the input channel
// TODO check if read only
func Clone(in interface{}) Piper {
	// Check if the value that was received is indeed a channel
	inVal := reflect.ValueOf(in)
	if inVal.Kind() != reflect.Chan {
		return nil
	}

	// Create the transformer
	t := transformer{
		inputChan:              inVal,
		outputChan:             reflect.MakeChan(inVal.Type(), inVal.Cap()),
		fChan:                  make(chan PipeOperator),
		funcs:                  list.New(),
		store:                  make(map[uint]interface{}),
		eventEmitterTransforms: make([]*list.Element, 0),
		transformCount:         0,
	}

	// Start the core which does the work
	go core(&t)

	return &t
}

// FromTo takes two channels and returns a piper that
// reads from inChan and outputs to outChan
func FromTo(inChan, outChan interface{}) Piper {
	// Check if the value that was received is indeed a channel
	inVal := reflect.ValueOf(inChan)
	outVal := reflect.ValueOf(outChan)
	if inVal.Kind() != reflect.Chan || outVal.Kind() != reflect.Chan {
		return nil
	}

	// Create the transformer
	t := transformer{
		inputChan:              inVal,
		outputChan:             outVal,
		fChan:                  make(chan PipeOperator),
		funcs:                  list.New(),
		store:                  make(map[uint]interface{}),
		eventEmitterTransforms: make([]*list.Element, 0),
		transformCount:         0,
	}

	// Start the core which does the work
	go core(&t)

	return &t
}
