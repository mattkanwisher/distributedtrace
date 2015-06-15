package zipkin

import (
	"fmt"
	"time"

	zkcore "github.com/mattkanwisher/distributedtrace/gen/zipkincore"
)

// TODO: Turn this into an interface.
type Combiner struct {
	config *Config
	inputs map[int64]chan *zkcore.Span
	output chan OutputMap

	stopped   bool
	semaphore chan bool
}

func NewCombiner(config *Config) *Combiner {
	combiner := &Combiner{
		config: config,
		inputs: map[int64]chan *zkcore.Span{},
		output: make(chan OutputMap, config.OutputBufferSize),

		stopped:   false,
		semaphore: make(chan bool, config.MaxConcurrentTraces),
	}

	// fill up the semaphore
	for i := 0; i < config.MaxConcurrentTraces; i++ {
		combiner.semaphore <- true
	}

	return combiner
}

// TODO: More graceful/waitable exit.
func (c *Combiner) Stop() {
	c.stopped = true
	for i := 0; i < c.config.MaxConcurrentTraces; i++ {
		<-c.semaphore
	}

	for traceId, channel := range c.inputs {
		close(channel)
		delete(c.inputs, traceId)
	}

	close(c.output)
}

func (c *Combiner) Send(span *zkcore.Span) error {
	if c.stopped {
		return fmt.Errorf("combiner stopped.")
	}

	if channel, ok := c.inputs[span.TraceId]; ok {
		channel <- span
		return nil
	}

	channel := make(chan *zkcore.Span, c.config.InputBufferSize)
	c.inputs[span.TraceId] = channel
	go func() {
		defer close(channel)

		spans := []*zkcore.Span{}
		timeout := time.After(c.config.TraceTimeout)

	for_loop:
		for {
			select {
			case s := <-channel:
				spans = append(spans, s)
			case <-timeout:
				break for_loop
			}
		}

		if outputs := c.combine(spans); len(outputs) > 0 {
			for _, output := range outputs {
				c.output <- output
			}
		}
	}()

	return c.Send(span)
}

func (c *Combiner) Receive() <-chan OutputMap {
	return c.output
}

func (c *Combiner) combine(spans []*zkcore.Span) []OutputMap {
	printf := c.config.Logger.Printf
	defer func() {
		if r := recover(); r != nil {
			printf("Combine(): failed to combine %d span(s): %s", len(spans), r)
		} else {
			printf("Combine(): combined %d span(s) into output.", len(spans))
		}
	}()

	const nameSep = "|"

	// since nodes may not come in exact parent-child order, we need a lookup here first so
	// we can reconstruct the tree properly.
	lookup := map[int64]*tree{}
	for _, span := range spans {
		outputMap, e := convertSpanToOutputMap(c.config, span)
		noError(e)

		node := &tree{
			id:     span.Id,
			parent: nil,
			name:   span.Name,

			span:      span,
			outputMap: outputMap,
		}

		lookup[node.id] = node
	}

	// rebuild the tree
	// TODO: Tree building should be a separate step. Then Combiner would be something more
	//   like a Flattener, or this could be decided at the Output level.
	// TODO: panic() less.
	var root *tree = nil
	for _, node := range lookup {
		switch {
		case node.span.ParentId != nil:
			if parentNode, exists := lookup[*node.span.ParentId]; exists {
				parentNode.children = append(parentNode.children, node)
				node.parent = parentNode

			} else { // !exists
				panic(fmt.Errorf("trace has orphaned spans (missing parent)"))

			}

		case root == nil:
			root = node
		default: // root != nil
			panic(fmt.Errorf("trace has multiple root spans (parentId == nil)!"))
		}
	}

	// compute relative names and times for each node.
	root.visitByBreadth(func(node *tree) bool {
		name := node.span.Name
		if node.parent != nil {
			name = node.parent.name + nameSep + name
		}

		node.name = name
		if d, ok := node.outputMap.SD(); ok {
			node.absTime = d
		} else if d, ok = node.outputMap.CD(); ok {
			node.absTime = d
		} else {
			c.config.Logger.Printf("Combine(): no stat to record for step: %s", name)
			node.absTime = 0
		}

		node.relTime = node.absTime
		if node.parent != nil {
			node.parent.relTime -= node.absTime
		}

		return true
	})

	// TODO: adds other non-stats fields (leaves overriding parents, right overrides left.)?
	results := []OutputMap{}
	root.visitByBreadth(func(node *tree) bool {
		result := OutputMap{}
		result["name"] = node.name
		result["absTime"] = node.absTime
		result["relTime"] = node.relTime
		results = append(results, result)
		return true
	})

	return results
}
