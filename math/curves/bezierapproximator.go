package curves

import (
	"github.com/Wieku/gosu-pp/math/vector"
	"sync"
)

const BezierQuantization = 0.5
const BezierQuantizationSq = BezierQuantization * BezierQuantization

type ItemStack struct {
	items [][]vector.Vector2f
	lock  sync.RWMutex
}

func NewStack() *ItemStack {
	return &ItemStack{items: make([][]vector.Vector2f, 0)}
}

// Push adds an Item to the top of the stack
func (s *ItemStack) Push(t []vector.Vector2f) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Pop removes an Item from the top of the stack
func (s *ItemStack) Pop() []vector.Vector2f {
	s.lock.Lock()
	item := s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]
	s.lock.Unlock()
	return item
}

func (s *ItemStack) Count() int {
	return len(s.items)
}

type BezierApproximator struct {
	count              int
	controlPoints      []vector.Vector2f
	subdivisionBuffer1 []vector.Vector2f
	subdivisionBuffer2 []vector.Vector2f
}

func NewBezierApproximator(controlPoints []vector.Vector2f) *BezierApproximator {
	return &BezierApproximator{count: len(controlPoints), controlPoints: controlPoints, subdivisionBuffer1: make([]vector.Vector2f, len(controlPoints)), subdivisionBuffer2: make([]vector.Vector2f, len(controlPoints)*2-1)}
}

func IsFlatEnough(controlPoints []vector.Vector2f) bool {
	for i := 1; i < len(controlPoints)-1; i++ {
		if controlPoints[i-1].Sub(controlPoints[i].Scl(2)).Add(controlPoints[i+1]).LenSq() > BezierQuantizationSq {
			return false
		}
	}

	return true
}

func (approximator *BezierApproximator) Subdivide(controlPoints, l, r []vector.Vector2f) {
	midpoints := approximator.subdivisionBuffer1

	for i := 0; i < approximator.count; i++ {
		midpoints[i] = controlPoints[i]
	}

	for i := 0; i < approximator.count; i++ {
		l[i] = midpoints[0]
		r[approximator.count-i-1] = midpoints[approximator.count-i-1]

		for j := 0; j < approximator.count-i-1; j++ {
			midpoints[j] = (midpoints[j].Add(midpoints[j+1])).Scl(0.5)

		}
	}
}

func (approximator *BezierApproximator) Approximate(controlPoints []vector.Vector2f, output *[]vector.Vector2f) {
	l := approximator.subdivisionBuffer2
	r := approximator.subdivisionBuffer1

	approximator.Subdivide(controlPoints, l, r)

	for i := 0; i < approximator.count-1; i++ {
		l[approximator.count+i] = r[i+1]
	}

	*output = append(*output, controlPoints[0])

	for i := 1; i < approximator.count-1; i++ {
		index := 2 * i
		p := (l[index-1].Add(l[index].Scl(2.0)).Add(l[index+1])).Scl(0.25)
		*output = append(*output, p)
	}
}

func (approximator *BezierApproximator) CreateBezier() []vector.Vector2f {
	output := make([]vector.Vector2f, 0)

	if approximator.count == 0 {
		return output
	}

	toFlatten := NewStack()
	freeBuffers := NewStack()

	// "toFlatten" contains all the curves which are not yet approximated well enough.
	// We use a stack to emulate recursion without the risk of running into a stack overflow.
	// (More specifically, we iteratively and adaptively refine our curve with a
	// <a href="https://en.wikipedia.org/wiki/Depth-first_search">Depth-first search</a>
	// over the tree resulting from the subdivisions we make.)

	nCP := make([]vector.Vector2f, len(approximator.controlPoints))

	copy(nCP, approximator.controlPoints)

	toFlatten.Push(nCP)

	leftChild := approximator.subdivisionBuffer2

	for toFlatten.Count() > 0 {
		parent := toFlatten.Pop()
		if IsFlatEnough(parent) {
			// If the control points we currently operate on are sufficiently "flat", we use
			// an extension to De Casteljau's algorithm to obtain a piecewise-linear approximation
			// of the bezier curve represented by our control points, consisting of the same amount
			// of points as there are control points.
			approximator.Approximate(parent, &output)
			freeBuffers.Push(parent)
			continue
		}

		// If we do not yet have a sufficiently "flat" (in other words, detailed) approximation we keep
		// subdividing the curve we are currently operating on.
		var rightChild []vector.Vector2f = nil
		if freeBuffers.Count() > 0 {
			rightChild = freeBuffers.Pop()
		} else {
			rightChild = make([]vector.Vector2f, approximator.count)
		}

		approximator.Subdivide(parent, leftChild, rightChild)

		// We re-use the buffer of the parent for one of the children, so that we save one allocation per iteration.
		for i := 0; i < approximator.count; i++ {
			parent[i] = leftChild[i]
		}

		toFlatten.Push(rightChild)
		toFlatten.Push(parent)
	}

	output = append(output, approximator.controlPoints[approximator.count-1])
	return output
}
