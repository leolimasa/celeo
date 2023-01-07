package cell

import "context"

type AsyncMethod[IN any, OUT any] struct {
	Params IN
	OutChan chan OUT
}

func Call[IN any, OUT any](inChan chan AsyncMethod[IN, OUT], params IN) chan OUT {
	// Buffer a single value so that the channel won't block if there
	// is no immediate listener.
	return CallOut(inChan, params, make(chan OUT, 1))
}

func CallOut[IN any, OUT any](inChan chan AsyncMethod[IN, OUT], params IN, outChan chan OUT) chan OUT {
	inChan <- AsyncMethod[IN, OUT] { 
		Params: params, 
		OutChan: outChan,
	}
	return outChan
}

func LoopWithContext(ctx context.Context, fn func ()) func() {
	return func () {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					fn()
				}
			}
		}()	
	}
}
