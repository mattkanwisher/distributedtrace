package zipkin

import (
	"encoding/base64"
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/mattkanwisher/distributedtrace/gen/scribe"
	zkcore "github.com/mattkanwisher/distributedtrace/gen/zipkincore"
	zkdep "github.com/mattkanwisher/distributedtrace/gen/zipkindependencies"
)

type Collector struct {
	*Config
	buffer chan *zkcore.Span
}

func NewCollector(config *Config) *Collector {
	return &Collector{
		Config: config,
		buffer: make(chan *zkcore.Span),
	}
}

func (c *Collector) Receive() <-chan *zkcore.Span {
	return c.buffer
}

func (c *Collector) Log(entries []*scribe.LogEntry) (code scribe.ResultCode, e error) {
	spans := make([]*zkcore.Span, 0, len(entries))
	for _, entry := range entries {
		if span, e := spanFromEntry(entry); e != nil {
			return scribe.ResultCode_TRY_LATER, e
		} else {
			spans = append(spans, span)
		}
	}

	c.Logger.Printf("Log(): %d span(s) received.", len(spans))
	for _, span := range spans {
		if span.ParentId != nil {
			c.Logger.Printf("span: t:%d i:%d p:%d n:%s", span.TraceId, span.Id, *span.ParentId, span.Name)
		} else {
			c.Logger.Printf("span: t:%d i:%d p:%d n:%s", span.TraceId, span.Id, 0, span.Name)
		}

		c.buffer <- span
	}

	return scribe.ResultCode_OK, nil
}

func (c *Collector) StoreDependencies(dependencies *zkdep.Dependencies) error {
	c.Logger.Print("StoreDependencies()")
	return nil
}

func (c *Collector) StoreTopAnnotations(serviceName string, annotations []string) error {
	c.Logger.Print("StoreTopAnnotations()")
	return nil
}

func (c *Collector) StoreTopKeyValueAnnotations(serviceName string, annotations []string) error {
	c.Logger.Print("StoreTopKeyValueAnnotations()")
	return nil
}

func spanFromEntry(entry *scribe.LogEntry) (span *zkcore.Span, e error) {
	var buffer []byte // TODO: Reuse buffer
	if buffer, e = base64.StdEncoding.DecodeString(entry.Message); e != nil {
		return nil, e
	}

	thriftBuffer := thrift.NewTMemoryBuffer()
	if n, e := thriftBuffer.Write(buffer); e != nil {
		return nil, e
	} else if n != len(buffer) {
		return nil, fmt.Errorf("buffer copy failure.")
	}

	protocol := thrift.NewTBinaryProtocol(thriftBuffer, true, true)
	span = &zkcore.Span{}
	if e := span.Read(protocol); e != nil {
		return nil, e
	}

	return span, nil
}
