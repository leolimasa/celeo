package main

import (
	"context"
	"fmt"

	"github.com/leolimasa/celeo/cell"
)

type IncrementIn struct{ Value int }

type IncrementOut struct{ Value int }

type State struct {
	number int
}

type AsyncCounterMethods struct {
	Increment chan cell.AsyncMethod[IncrementIn, IncrementOut]
}

func NewAsyncCounterInputs() AsyncCounterMethods {
	self := AsyncCounterMethods{
		Increment: make(chan cell.AsyncMethod[IncrementIn, IncrementOut], 512),
	}
	return self
}

func AsyncCounterMain(meth AsyncCounterMethods, value int, ctx context.Context) func() {
	state := State{number: value}
	return cell.LoopWithContext(ctx, func() {
		select {
		case call := <-meth.Increment:
			state.number = state.number + call.Params.Value
			call.OutChan <- IncrementOut{Value: state.number}
			close(call.OutChan)
		default:

		}
	})
}

func main() {
	rootCtx := context.Background()
	counterCtx, stopCounter := context.WithCancel(rootCtx)
	counter := NewAsyncCounterInputs()
	AsyncCounterMain(counter, 0, counterCtx)()
	cell.Call(counter.Increment, IncrementIn{Value: 10})
	newVal := <-cell.Call(counter.Increment, IncrementIn{Value: 15})
	fmt.Println(newVal)
	stopCounter()
}
