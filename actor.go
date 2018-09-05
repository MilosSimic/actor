package main

import (
	"context"
	"strings"
)

type ActorProp struct {
	Context *System
	name    string
	box     chan interface{}
	resp    chan interface{}
	kill    chan bool
	watch   chan string
	parrent *ActorProp
	state   State
}

func (a *ActorProp) start(ctx context.Context, actor Actor) {
	go func() {
		for {
			select {
			case <-a.kill:
				a.notify()
				a.clean(true)
				return
			case msg := <-a.box:
				actor.Receive(msg, a)
			case <-ctx.Done():
				a.notify()
				a.clean(false)
				return
			}
		}
	}()
}

func (a *ActorProp) clean(remove bool) {
	close(a.kill)
	close(a.box)
	close(a.resp)
	close(a.watch)

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

func (a *ActorProp) TellChildren(msg interface{}) {
	key := join(a.Context.Name, a.name)
	for k, v := range a.Context.Actors {
		if strings.Contains(k, key) {
			v.Tell(msg)
		}
	}
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
		watch:   make(chan string),
		Context: system,
		state:   NormalState{},
	}
}

func (a *ActorProp) Become(state State) {
	a.state = state
}

func (a *ActorProp) ActorOf(name string, val Actor) *ActorProp {
	key := join(a.name, name)
	child := a.Context.ActorOf(key, val)
	child.parrent = a

	return child
}

func (a *ActorProp) notify() {
	if a.parrent != nil {
		a.parrent.watch <- a.name
	}
}

func (a *ActorProp) Watch(aprs ...*ActorProp) {
	for _, ap := range aprs {
		ap.parrent = a
		a.Context.watch(a)
	}
}
