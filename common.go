package main

import (
	"context"
	"strings"
)

func (s *System) Terminate() {
	s.cancel()
}

func join(name string) string {
	system := "system"

	return strings.Join([]string{system, name}, "/")
}

func (a *ActorProp) clean(remove bool) {
	close(a.kill)
	close(a.box)

	if remove {
		delete(a.Context.Actors, a.name)
	}
}

func (a *ActorProp) Kill() {
	a.kill <- true
}

func (a *ActorProp) Tell(msg interface{}) {
	a.box <- msg
}

func (a *ActorProp) Replay(msg interface{}) {
	a.resp <- msg
}

func (a *ActorProp) Resp() interface{} {
	return <-a.resp
}

func newProp(name string, system *System) *ActorProp {
	return &ActorProp{
		name:    name,
		box:     make(chan interface{}),
		resp:    make(chan interface{}),
		kill:    make(chan bool),
		Context: system,
	}
}

func NewSystem(name string) *System {
	ctx, cancel := context.WithCancel(context.Background())
	return &System{
		Name:   name,
		Actors: map[string]*ActorProp{},
		ctx:    ctx,
		cancel: cancel,
	}
}
