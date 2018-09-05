package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

type System struct {
	Name   string
	Actors map[string]*ActorProp
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *System) ActorOf(name string, val Actor) *ActorProp {
	_, ok := interface{}(val).(Actor) // test does val implement Actor interfce

	if ok {
		key := join(s.Name, name)
		prop := newProp(key, s)
		prop.start(s.ctx, val)
		s.Actors[key] = prop

		return prop
	}

	return nil
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

func (s *System) Terminate(ap *ActorProp) bool {
	if strings.Contains(ap.name, ap.name) {
		delete(s.Actors, ap.name)
		ap.notify()

		return true
	}

	return false
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

func (s *System) Shutdown() {
	s.cancel()
}
