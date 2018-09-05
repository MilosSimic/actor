package main

import (
	"fmt"
	"time"
)

type MyMessage struct{}

func (m MyMessage) Name() string {
	return "MyMessage"
}

func (m MyMessage) Params() map[string][]byte {
	return nil
}

type ChangeState struct{}

func (m ChangeState) Name() string {
	return "MyMessage"
}

func (m ChangeState) Params() map[string][]byte {
	return nil
}

func (n NormalState) Execute(msg interface{}, context *ActorProp) {
	switch conv := msg.(type) {
	case MyMessage:
		fmt.Println("Hello ", conv.Name(), "From ", context.name)
		context.Replay("Replay")
	case string:
		fmt.Println("Killed ", conv)
	default:
		context.Become(MyState{})
		fmt.Println("State Changed!")
	}
}

type MyState struct{}

func (ms MyState) Execute(msg interface{}, context *ActorProp) {
	switch conv := msg.(type) {
	case MyMessage:
		fmt.Println("I'm in my state")
	case string:
		fmt.Println("Killed fuck yeah", conv)
	default:
		fmt.Println("bla")
		context.Become(NormalState{})
	}
}

type MyActor struct{}

func (m MyActor) Receive(msg interface{}, context *ActorProp) {
	switch msg.(type) {
	case ChangeState:
		context.Become(MyState{})
	default:
		context.state.Execute(msg, context)
	}
}

func main() {
	system := NewSystem("TestSystem")

	p1 := system.ActorOf("Parrent_Actor_1", MyActor{})

	p1.Tell(ChangeState{})
	p1.Tell(MyMessage{})

	p1.Tell(nil)
	p1.Tell(MyMessage{})

	// c1 := p1.ActorOf("Child_Actor_1", MyActor{})
	// c2 := p1.ActorOf("Child_Actor_2", MyActor{})
	// p1.Watch(c1, c2)

	// p2 := system.ActorOf("Parrent_Actor_2", MyActor{})
	// c3 := p2.ActorOf("Child_Actor_3", MyActor{})
	// c4 := p2.ActorOf("Child_Actor_4", MyActor{})
	// p2.Watch(c3, c4)

	// system.AllActors()

	// p.Tell(MyMessage{})
	// p1.TellChildren(MyMessage{})
	// fmt.Println(p.Resp())

	// c1.Kill()
	// c2.Kill()

	time.Sleep(time.Second)

	// system.AllActors()
	system.Shutdown()
}
