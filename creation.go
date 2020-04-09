package piper

import (
	"container/list"
	"reflect"
)

// From is used to create a Piper from a given channel
func From(in interface{}) Piper {
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
