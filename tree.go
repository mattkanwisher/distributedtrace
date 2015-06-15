package zipkin

import (
	zkcore "github.com/mattkanwisher/distributedtrace/gen/zipkincore"
)

type tree struct {
	id     int64
	parent *tree
	name   string

	absTime int64
	relTime int64

	children  []*tree
	span      *zkcore.Span
	outputMap OutputMap
}

func (t *tree) visitByBreadth(visitor func(*tree) bool) bool {
	var next *tree
	queue := []*tree{t}
	for len(queue) > 0 {
		next, queue = queue[0], queue[1:]
		queue = append(queue, next.children...)

		if !visitor(next) {
			return false
		}
	}

	return true
}

func (t *tree) visitByDepth(visitor func(*tree) bool) bool {
	for _, child := range t.children {
		if !child.visitByDepth(visitor) {
			return false
		} else if !visitor(child) {
			return false
		}
	}

	return visitor(t)
}

func (t *tree) childWithId(id int64) *tree {
	for _, node := range t.children {
		if node.id == id {
			return node
		} else if child := node.childWithId(id); child != nil {
			return child
		}
	}

	return nil
}
