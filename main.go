package main

import (
	"context"
	"fmt"
	"time"
)

type Actor interface {
	Receive(msg interface{}, context *ActorProp)
}

type ActorProp struct {
	Context *System
	name    string
	box     chan interface{}
	resp    chan interface{}
	kill    chan bool
	watch   chan string
	parrent *ActorProp
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

type Message interface {
	Name() string
	Params() map[string][]byte
}

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

type MyActor struct{}

type MyMessage struct{}

func (m MyMessage) Name() string {
	return "MyMessage"
}

func (m MyMessage) Params() map[string][]byte {
	return nil
}

func (m MyActor) Receive(msg interface{}, context *ActorProp) {
	switch conv := msg.(type) {
	case MyMessage:
		fmt.Println("Hello ", conv.Name(), "From ", context.name)
		context.Replay("Replay")
	case string:
		fmt.Println("Killed ", conv)
	default:
		fmt.Println("bla")
	}
}

func main() {
	system := NewSystem("TestSystem")

	p1 := system.ActorOf("Parrent_Actor_1", MyActor{})
	c1 := p1.ActorOf("Child_Actor_1", MyActor{})
	c2 := p1.ActorOf("Child_Actor_2", MyActor{})
	p1.Watch(c1, c2)

	p2 := system.ActorOf("Parrent_Actor_2", MyActor{})
	c3 := p2.ActorOf("Child_Actor_3", MyActor{})
	c4 := p2.ActorOf("Child_Actor_4", MyActor{})
	p2.Watch(c3, c4)

	system.AllActors()

	// p.Tell(MyMessage{})
	p1.TellChildren(MyMessage{})
	// fmt.Println(p.Resp())

	// c1.Kill()
	// c2.Kill()

	time.Sleep(time.Second)

	// system.AllActors()
	system.Shutdown()
}
