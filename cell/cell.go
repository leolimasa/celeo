package cell

import (
	"sync"
)

type ICell interface {
	StopChan() chan bool
	Parent() ICell
	AddChild(child ICell) ICell
	Destroy()
	setParent(parent ICell)
	deleteChild(child ICell)
}

type Cell[I any, O any] struct {
	Input     chan I
	Output    chan O
	Stop      chan bool
	children  []ICell
	parent    ICell
	childLock sync.Mutex
}

func NewCell[I any, O any](outputChan chan O) *Cell[I, O] {
	c := &Cell[I, O]{
		Input:  make(chan I, 512),
		Output: outputChan,
		Stop:   make(chan bool),
	}
	return c
}

func (self *Cell[I, O]) AddChild(child ICell) ICell {
	self.childLock.Lock()
	defer self.childLock.Unlock()
	self.children = append(self.children, child)
	child.setParent(self)
	return child
}

func (self *Cell[I, O]) setParent(parent ICell) {
	self.parent = parent
}

func (self *Cell[I, O]) Parent() ICell {
	return self.parent
}

func (self *Cell[I, O]) StopChan() chan bool {
	return self.Stop
}

func (self *Cell[I, O]) deleteChild(child ICell) {
	self.ForEachChild(func(i int, curChild ICell) bool {
		if curChild == child {
			self.children = append(self.children[:i], self.children[i+1:]...)
			child.setParent(nil)
			return false
		}
		return true
	})
}

func (self *Cell[I, O]) Destroy() {
	self.Stop <- true
	self.ForEachChild(func(i int, child ICell) bool {
		child.Destroy()
		return true
	})
	if self.parent != nil {
		self.parent.deleteChild(self)
	}
}

// ForEachChild executes the specified function for each child.
// The function should return whether to continue the loop or not.
func (self *Cell[I, O]) ForEachChild(fn func(i int, child ICell) bool) {
	self.childLock.Lock()
	defer self.childLock.Unlock()
	for i, child := range self.children {
		if !fn(i, child) {
			break
		}
	}
}

// ----------------------------------------
// Cell Collection
// ----------------------------------------

func NewCellCollection[I any, O any]() *Cell[I, O] {
	self := NewCell[I](make(chan O))
	go func() {
		defer self.Destroy()
		for {
			select {
			case <-self.Stop:
				return

			// Redirect all input messages to the children
			case msg := <-self.Input:
				self.ForEachChild(func(i int, child ICell) bool {
					child.(*Cell[I,O]).Input <- msg
					return true

				})
			}
		}
	}()
	return self
}
