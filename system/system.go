package system

import (
	"context"
	"github.com/milossimic/actor/prop"
)

type System struct {
	Name   string
	Actors map[string]*prop.ActorProp
	ctx    context.Context
	cancel context.CancelFunc
}

func New(name string) *System {
	ctx, cancel := context.WithCancel(context.Background())
	return &System{
		Name:   name,
		Actors: map[string]*prop.ActorProp{},
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *System) Terminate() {
	s.cancel()
}

func (s *System) ActorOf(name string, val prop.Actor) *prop.ActorProp {
	_, ok := interface{}(val).(prop.Actor) // test does val implement Actor interfce
	if ok {
		key := join(name)
		prop := prop.New(name, s)
		go prop.Start(s.ctx, val)
		s.Actors[key] = prop

		return prop
	}

	return nil
}
