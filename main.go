package main

import (
	"fmt"
	"github.com/milossimic/actor/system"
	"time"
)

// type Actor interface {
// 	Receive(context *ActorProp, msg interface{})
// }

// type ActorProp struct {
// 	Context *System
// 	name    string
// 	box     chan interface{}
// 	resp    chan interface{}
// 	kill    chan bool
// }

// func (a *ActorProp) start(ctx context.Context, actor Actor) {
// 	for {
// 		select {
// 		case <-a.kill:
// 			a.clean(true)
// 			return
// 		case msg := <-a.box:
// 			actor.Receive(a, msg)
// 		case <-ctx.Done():
// 			a.clean(false)
// 			return
// 		}
// 	}
// }

type Message interface {
	Name() string
	Params() map[string][]byte
}

// type System struct {
// 	Name   string
// 	Actors map[string]*ActorProp
// 	ctx    context.Context
// 	cancel context.CancelFunc
// }

func (s *System) ActorOf(name string, val Actor) *ActorProp {
	_, ok := interface{}(val).(Actor) // test does val implement Actor interfce

	if ok {
		key := join(name)
		prop := newProp(name, s)
		go prop.start(s.ctx, val)
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

func (m MyActor) Receive(context *ActorProp, msg interface{}) {
	switch conv := msg.(type) {
	case MyMessage:
		fmt.Println("Hello ", conv.Name())
		context.Replay("Replay")
	default:
		fmt.Println("bla")
	}
}

func main() {
	sys := system.New("TestSystem")
	prop := sys.ActorOf("MyActor", MyActor{})

	prop.Tell(MyMessage{})
	fmt.Println("Resp:", prop.Resp())

	time.Sleep(time.Second)

	system.Terminate()

	fmt.Println("hello world")
}
