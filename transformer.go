package piper

import (
	"container/list"
	"fmt"
	"reflect"
)

// transforms are used to wrap a pipefunc with some additional information
// transforms are stored in the lists inside the transformer
type transform struct {
	// uid is the identifier for the transform which can be used to retrieve
	// values from the store located in the transformer
	uid uint
	// this is the flag if there should be the effort to check the storage for the uid or not
	hasStorage bool
	// the function to apply to the element
	f Pipefunc
}

// transformer implements Piper
type transformer struct {
	// the channel to listen on for values
	inputChan reflect.Value
	// the channel to emit value to
	outputChan reflect.Value
	// fChan is the channel you can add new transforms to the transformer
	fChan chan PipeOperator
	// funcs is a list storing pointers to the transforms to apply to a value that is received
	// they are called in succession
	funcs *list.List
	// store saves the storage of the pipefuncs
	// it is accessed via the uid of the transform. Because not every (more like a few) transforms
	// need a storage so we only have an entry if we really need one
	store map[uint]interface{}
	// the index corresponds to the transform that triggered the event.
	// this way upon receiving values from an EventEmitter we know who to call
	eventEmitterTransforms []*list.Element
	// counter to keep track of the amount of transforms that there are
	transformCount uint
}

// Pipe takes a variable amount of Pipefuncs and registers them in the transformer
func (t *transformer) Pipe(ops ...PipeOperator) Piper {
	for _, o := range ops {
		t.fChan <- o
	}

	return t
}

// Get returns the output channel it needs to be asserted like this: .Get().(chan int)
func (t *transformer) Get() interface{} {
	return t.outputChan.Interface()
}

// core does all the event handling of the channels
func core(t *transformer) {
	// Initially there are two channels that we receive data from
	// the fChan to receive new transforms and the input channel
	const numCoreChans int = 2
	cases := make([]reflect.SelectCase, numCoreChans)
	cases[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(t.fChan),
	}
	cases[1] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: t.inputChan,
	}

	// this is used to tell future runs if the input channel has been closed before or not
	inClosed := false

	// this can be used so handleMsg can cut the connection to the
	// input channel
	closeCallback := func() {
		cases[1].Chan = reflect.ValueOf(nil)
		inClosed = true
	}

main:
	for {
		chosen, val, ok := reflect.Select(cases)
		// handle closing of channels
		if !ok {
			switch chosen {
			// if the input channel closes, disable it
			case 1:
				cases[1].Chan = reflect.ValueOf(nil)
				inClosed = true
			// when a event emitter closes remove it and forget about it
			default:
				// throw the event emitter out
				cases = append(cases[:chosen], cases[chosen+1:]...)
				// throw the reference to the transform for that event emiter out
				t.eventEmitterTransforms = append(t.eventEmitterTransforms[:chosen-numCoreChans], t.eventEmitterTransforms[chosen-numCoreChans+1:]...)
				// this is not propagated
				continue main
			}
		}

		// now handle the message
		switch chosen {
		// Receiving from fChan
		case 0:
			n := val.Interface().(PipeOperator)
			handleNewTransform(t, &n, &cases)
		// Receiving from inputChan
		case 1:
			// This happens n times when open and one time when closed
			handleMsg(t, t.funcs.Front(), toPipeResult(&val, ok, false), closeCallback)
		// in every other case it is a registered EventEmitter
		default:
			// Because there are the two core channels we need to subtract them to get the transform
			// that belongs to the EventEmiter and handle the message with it
			handleMsg(t, t.eventEmitterTransforms[chosen-numCoreChans], toPipeResult(&val, true, true), closeCallback)
		}

		// if the input channel is closed and all pipefuncs have been removed after the run
		// we can stop
		if inClosed && t.funcs.Len() == 0 {
			break main
		}
	}

	// Clean up the channels
	close(t.fChan)
	t.outputChan.Close()
}

// used to add transforms to the transformer
func handleNewTransform(t *transformer, n *PipeOperator, cases *[]reflect.SelectCase) {
	// Create the transform
	trans := transform{
		f:          n.F,
		uid:        t.transformCount,
		hasStorage: false,
	}

	// Add the intial storage into the map if there is one
	if n.InitialStorage != nil {
		trans.hasStorage = true
		t.store[trans.uid] = n.InitialStorage
	}

	// Add the transform to execution chain
	element := t.funcs.PushBack(&trans)

	// save the new event emitter
	if n.EventEmitter != nil {
		// save the starting entry for the execution chain when there is an event emitted
		t.eventEmitterTransforms = append(t.eventEmitterTransforms, element)
		// append a new entry to the select cases
		*cases = append(*cases, reflect.SelectCase{Chan: reflect.ValueOf(n.EventEmitter), Dir: reflect.SelectRecv})
	}

	t.transformCount++
}

func handleMsg(t *transformer, e *list.Element, result PipeResult, closeInChan func()) {
chain:
	// Pump it through the chain
	for e != nil {
		// get the transform that belongs to the element
		trans := e.Value.(*transform)

		// call the pipefunc and save it's returned storage
		nResult, nStorage := trans.f(result, getStorage(t, trans))
		saveStorage(t, trans, nStorage)

		switch nResult.State {
		// When the channel is open but a IsValue=false is returned, equals the Break behaviour
		case Open:
			if !nResult.IsValue {
				result = nResult
				break chain
			}
		// Break stops the execution chain
		case Break:
			result = nResult
			break chain
		// When a transformer is closed without the predecessor being closed results in a panic.
		// If it was then remove the transformer from the chain
		case Closed:
			// When a pipefunc closes in the middle of an open chain we stop listening to the
			// input channel and issue a close call to all predecessors before removing them
			if result.State != Closed && !result.IsEvent {
				closeInChan()
				killPredecessors(t, e)
			}
			result = nResult
			// Remove the transform and it's storage
			delete(t.store, trans.uid)
			toDelete := e
			e = e.Next()
			t.funcs.Remove(toDelete)
			continue chain
		}

		// take the result to the next run which cannot be an event anymore
		result = nResult
		result.IsEvent = false
		e = e.Next()
	}

	// if the final PipeResult contains a value and it was not broken out of the chain then emit
	if result.State != Break && result.IsValue {
		// We can only sent a value on a channel if the type is asignable to the type of the channel
		// if this is possible then send the value otherwise panic. Go would panic anyways but this way
		// there is a more detailed cause of it specified.
		if reflect.TypeOf(result.Value).AssignableTo(t.outputChan.Type().Elem()) {
			t.outputChan.Send(reflect.ValueOf(result.Value))
		} else {
			panic(fmt.Sprintf("You tried to send a value of type %v on a channel of type %s. The value sent must be assignable to the type of the channel",
				reflect.TypeOf(result.Value), t.outputChan.Type().Elem()))
		}
	}
}

// Given the variables from the channel return a PipeResult to act as the
// initial PipeResult for the others.
func toPipeResult(val *reflect.Value, ok, isEvent bool) PipeResult {
	// Determine state
	var state PipeState
	if ok {
		state = Open
	} else {
		state = Closed
	}

	// !ok && isEvent => can never happen

	return PipeResult{
		Value:   val.Interface(),
		State:   state,
		IsEvent: isEvent,
		IsValue: ok,
	}
}

// Returns the storage of a transformer or nil if there is none
func getStorage(t *transformer, trans *transform) interface{} {
	// get the storage if there is one
	var storage interface{} = nil
	if trans.hasStorage {
		storage = t.store[trans.uid]
	}
	return storage
}

// saves the newStorage of a transformer
func saveStorage(t *transformer, trans *transform, newStorage interface{}) {
	// if the new storage is not nil, save it
	if newStorage != nil {
		trans.hasStorage = true
		t.store[trans.uid] = newStorage
	}
	// if it had storage before but is now nil then we should delete the storage entry
	if trans.hasStorage && newStorage == nil {
		trans.hasStorage = false
		delete(t.store, trans.uid)
	}
}

// Calls all predecessors with state closed so they can shutdown their EventEmtiters
// their output will be ignored because they are not important for the chain anymore
func killPredecessors(t *transformer, of *list.Element) {
	e := of.Prev()
	for e != nil {
		// call the pipefunc with closed so it can close the emitters but
		// ignore it's result
		trans := e.Value.(*transform)
		trans.f(PipeResult{IsValue: false, State: Closed}, getStorage(t, trans))
		// delete the storage of the pipefunc
		delete(t.store, trans.uid)
		// go back while deleting the entry
		toDelete := e
		e = e.Prev()
		t.funcs.Remove(toDelete)
	}
}
