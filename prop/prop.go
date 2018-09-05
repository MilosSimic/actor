package prop

import (
	"context"
)

type Actor interface {
	Receive(context *ActorProp, msg interface{})
}

type ActorProp struct {
	name string
	box  chan interface{}
	resp chan interface{}
	kill chan bool
}

func (a *ActorProp) Start(ctx context.Context, actor Actor) {
	for {
		select {
		case <-a.kill:
			a.clean(true)
			return
		case msg := <-a.box:
			actor.Receive(a, msg)
		case <-ctx.Done():
			a.clean(false)
			return
		}
	}
}

func (a *ActorProp) clean(remove bool) {
	close(a.kill)
	close(a.box)
	close(a.resp)

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

func New(name string) *ActorProp {
	return &ActorProp{
		name: name,
		box:  make(chan interface{}),
		resp: make(chan interface{}),
		kill: make(chan bool),
	}
}
