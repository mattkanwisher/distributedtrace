package zipkin

import (
	zkcore "github.com/mattkanwisher/distributedtrace/gen/zipkincore"
)

// The Output interface defines a handlers for handling incoming ZipKin spans. The
// Configure() method will be called once during initialization. Implementation should
// ensure that all the output configuration works (i.e. sends a PING) before returning
// from Configure()
type Output interface {
	Write(result OutputMap) error
}

type OutputMap map[string]interface{}

// Convenience for some interesting fields.
func (m OutputMap) CS() (int64, bool) { return m.get(zkcore.CLIENT_SEND) }
func (m OutputMap) CR() (int64, bool) { return m.get(zkcore.CLIENT_RECV) }
func (m OutputMap) SS() (int64, bool) { return m.get(zkcore.SERVER_SEND) }
func (m OutputMap) SR() (int64, bool) { return m.get(zkcore.SERVER_RECV) }

func (m OutputMap) CD() (int64, bool) { return m.get("cd") }
func (m OutputMap) SD() (int64, bool) { return m.get("sd") }

func (m OutputMap) get(key string) (int64, bool) {
	if result, ok := m[key]; ok {
		return result.(int64), ok
	} else {
		return 0, false
	}
}
