package piper

import "fmt"

/*
FACTS:
- The Piper will run until all pipefuncs have closed
- Returning IsValue as false while staying open is the same as returing break.
	This is usefull because this way you can always close with the predecessor without checking
	it's state while also being able to break the chain.
- If a typed function like MapInt does not receive a value that can be type asserted
	to a Integer this will also Break because this is not usable.
- If you return Closed when the predecessor is still open means that we wont read from the
	input channel anymore therefore only all following pipefuncs will be called which might
	end in the piper to end and the output channel to be closed
- Sending a value that is not assignable to the channeltype results in a panic
- Dont close the output channel
- Dont send into the output channel. It wont be read.
*/

// PipeOperator is used to register a new operation to the Piper
type PipeOperator struct {
	F              Pipefunc
	InitialStorage interface{}
	EventEmitter   chan interface{}
}

// PipeState tells what state the pipefunc is called in, see below
type PipeState uint8

const (
	// Open tells the pipefunc that it's predecessor is open. This is also the case for Events.
	// When the predecessor is closed you should not close
	Open PipeState = iota
	// Closed is received when the predecessor closed and now it's a good time to close aswell together with
	// any registered EventEmitters. But you can also stay open or break. But be aware that there wont be
	// any input to that pipefunc anymore except through their own EventEmitters. If you dont close then
	// the Piper will run forever
	Closed
	// Break is used to tell the Piper to exit out of the pipe execution chain
	Break
)

// PipeResult is the Result of a Pipefunc which will be either
// emitted or passed on to the next Pipefunc
type PipeResult struct {
	Value   interface{}
	IsValue bool
	IsEvent bool
	State   PipeState
}

// Pipefunc is a function that when a value and the current state of the Piper
// decides what the new value and new state will be. This can be archieved with
// it's own storage.
// receives: 	- the PipeResult of the predecessor
//				- the storage of the pipefunc or nil if there is none
// returns: 	- the PipeResult of the current Pipefunc
//				- the new storage of the Pipefunc
type Pipefunc func(PipeResult, interface{}) (PipeResult, interface{})

// Piper is an interface for ...
type Piper interface {
	// Pipe takes n pipefuncs what all initially have no storage
	// Note the event emitter has to be closed by the pipefunc otherwise the piper runs forever
	Pipe(...PipeOperator) Piper
	// Get yields the output channel
	// TODO this is too unsafe
	Get() interface{}
}

func (r PipeResult) String() string {
	return fmt.Sprintf("Value: %v | IsValue: %v | IsEvent: %v | State: %v", r.Value, r.IsValue, r.IsEvent, r.State)
}
