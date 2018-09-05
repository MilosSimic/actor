package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

func (s *System) Shutdown() {
	s.cancel()
}

func join(parts ...string) string {
	return strings.Join(parts, "/")
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
	}
}

func (a *ActorProp) ActorOf(name string, val Actor) *ActorProp {
	key := join(a.name, name)
	child := a.Context.ActorOf(key, val)
	child.parrent = a

	return child
}

func (s *System) watch(parrent *ActorProp) {
	go func() {
		for {
			select {
			case path := <-parrent.watch:
				name := filepath.Base(path)
				parrent.Tell(name)
			case <-s.ctx.Done():
				return
			}
		}
	}()
}

func (a *ActorProp) notify() {
	if a.parrent != nil {
		a.parrent.watch <- a.name
	}
}

func (s *System) Terminate(ap *ActorProp) bool {
	if strings.Contains(ap.name, ap.name) {
		delete(s.Actors, ap.name)
		ap.notify()

		return true
	}

	return false
}

func (a *ActorProp) Watch(aprs ...*ActorProp) {
	for _, ap := range aprs {
		ap.parrent = a
		a.Context.watch(a)
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

func (s *System) AllActors() {
	for k, _ := range s.Actors {
		fmt.Println(k)
	}
}
