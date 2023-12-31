package statemachine

import (
	"errors"
	"fmt"
	"sync"
)

// ErrEventRejected is the error returned when the state machine cannot process
// an event in the state thta it is in.
var ErrEventRejected = errors.New("event rejected")

const (
	// Default represents the default state of the system.
	Default StateType = ""
	
	// NoOp represents the no-op event.
	NoOp EventType = "NoOp"
)

// StateType represents an extensible state type in the state machine
type StateType string

// EventType represents an extensible event type in the state machine
type EventType string

// EventContext represents the context to be passed to the action implementation
type EventContext interface{}

// Action represents the action to be executed in the given state.
type Action interface {
	Execute(eventCtx EventContext) EventType
}

// Events represents a mapping of events and states.
type Events map[EventType]StateType

// State binds a state with an action and a set of events it can handle.
type State struct {
	Action Action
	Events Events
}

// State represents a mapping of states and their implementations
type States map[StateType]State

// StateMachine represents the state machine
type StateMachine struct {

	// Previous represents the previous state.
	Previous StateType

	// Current represents the current state.
	Current StateType

	// States hold the configuration of states and events handled by state machine.
	States States

	// mutex ensures that only 1 event is processed by the state machine at any given time
	mutex sync.Mutex
}


// getNextState returns the next state for the event given the machine's current
// state, or an error if the event can't be handled in the given state.
func (s *StateMachine) getNextState(event EventType) (StateType, error) {

	if state, ok := s.States[s.Current]; ok {
		if state.Events != nil {
			if next, ok := state.Events[event]; ok {
				return next, nil
			}
		}
	}

	return Default, ErrEventRejected
}

// SendEvent sends an event to the state machine
func (s *StateMachine) SendEvent(event EventType, eventCtx EventContext) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for {

		// Determine the next state for the event given the machine's currnet state
		nextState, err := s.getNextState(event)
		if err != nil {
			return ErrEventRejected
		}

		// Identify the state defination for the next state
		state, ok := s.States[nextState]
		if !ok || state.Action == nil {
			// config error
		}

		// Transition over to the next state.
		s.Previous = s.Current
		s.Current = nextState

		// Execute the next state's action and loop over again if event returned
		nextEvent := state.Action.Execute(eventCtx)
		if nextEvent == NoOp {
			return nil
		}
		event = nextEvent
	}
}


const (
	Off StateType = "Off"
	On StateType = "On"

	SwitchOff EventType = "SwitchOff"
	SwitchOn EventType = "SwitchOn"
)


// OffAction represents the action executed on entering the off state
type OffAction struct{}

func (a *OffAction) Execute(eventCtx EventContext) EventType {
	fmt.Println("The light has been switched off")
	return NoOp
}

// OnAction represents the action executed on entering the On state
type OnAction struct{}

func (a *OnAction) Execute(eventCtx EventContext) EventType {
	fmt.Println("The light has been switched on")
	return NoOp
}

func newLightSwitchFSM() *StateMachine {
  return &StateMachine{
		States: States{
			Default: State{
				Events: Events{
					SwitchOff: Off,
				},
			},
			Off: State{
				Action: &OffAction{},
				Events: Events{
					SwitchOn: On,
				},
			},
			On: State{
				Action: &OnAction{},
				Events: Events{
					SwitchOff: Off,
				},
			},
		},
	}
}

