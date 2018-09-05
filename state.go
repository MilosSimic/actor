package main

type State interface {
	Execute(msg interface{}, context *ActorProp)
}

type NormalState struct{}

// func (n NormalState) Execute(msg interface{}, context *ActorProp) {

// }
