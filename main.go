package main

import (
	"fmt"
	"time"
)

type TrafficState interface {
	Enter()
	Exit()
	Update(l *TrafficStateMachine)
}

type TrafficStateMachine struct {
	currentState TrafficState
	states       map[string]TrafficState
}

func (sm *TrafficStateMachine) setState(s TrafficState) {
	sm.currentState = s
	sm.currentState.Enter()
}

func (sm *TrafficStateMachine) Transition() {
	sm.currentState.Update(sm)
}

type RedLight struct{}

func (g RedLight) Enter() {
	fmt.Println("Red light is on. Stop driving")
	time.Sleep(time.Second * 5)
}

func (g RedLight) Exit() {}
func (g RedLight) Update(l *TrafficStateMachine) {
	l.setState(&GreenLight{})
}

type GreenLight struct{}

func (g GreenLight) Enter() {
	fmt.Println("Green light is on. You can drive.")
	time.Sleep(time.Second * 5)
}
func (g GreenLight) Exit() {}
func (g GreenLight) Update(l *TrafficStateMachine) {
	l.setState(&YellowLight{})
}

type YellowLight struct{}

func (g YellowLight) Enter() {
	fmt.Println("Yellow light is on. Prepare to stop.")
	time.Sleep(time.Second * 5)
}
func (g YellowLight) Exit() {}
func (g YellowLight) Update(l *TrafficStateMachine) {
	l.setState(&RedLight{})
}

func NewStateMachine(initialState TrafficState) *TrafficStateMachine {
	sm := &TrafficStateMachine{
		currentState: initialState,
		states:       make(map[string]TrafficState),
	}

	sm.currentState.Enter()
	return sm
}

func main() {
  sm := NewStateMachine(&RedLight{})

	for {
		sm.Transition()
	}
}
