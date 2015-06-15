package zipkin

import (
	"bytes"
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	zkcore "github.com/mattkanwisher/distributedtrace/gen/zipkincore"
)

func convertBinaryAnnotationFromThrift(ann *zkcore.BinaryAnnotation) (string, interface{}, error) {
	buffer := &thrift.TMemoryBuffer{}
	buffer.Buffer = bytes.NewBuffer(ann.Value)
	decoder := thrift.NewTBinaryProtocolTransport(buffer)

	var val interface{}
	var e error
	switch ann.AnnotationType {
	case zkcore.AnnotationType_BOOL:
		val, e = decoder.ReadBool()
	case zkcore.AnnotationType_BYTES:
		val, e = ann.Value, nil
	case zkcore.AnnotationType_I16:
		val, e = decoder.ReadI16()
	case zkcore.AnnotationType_I32:
		val, e = decoder.ReadI32()
	case zkcore.AnnotationType_I64:
		val, e = decoder.ReadI64()
	case zkcore.AnnotationType_DOUBLE:
		val, e = decoder.ReadDouble()
	case zkcore.AnnotationType_STRING:
		// val, e = decoder.ReadString()
		val, e = string(ann.Value), nil
	default:
		e = fmt.Errorf("unrecognized AnnotationType: %#v", ann.AnnotationType)
	}

	if e != nil {
		return "", nil, e
	}

	return ann.Key, val, nil
}

func convertSpanToOutputMap(config *Config, span *zkcore.Span) (OutputMap, error) {
	fields := OutputMap{
		"id":      span.Id,
		"traceId": span.TraceId,
		"name":    span.Name,
	}

	if span.ParentId != nil {
		fields["parentId"] = *span.ParentId
	} else {
		fields["parentId"] = nil
	}

	// flatten annotations into k/v map for the series
	for _, ann := range span.Annotations {
		fields[ann.Value] = ann.Timestamp
		if ann.Duration != nil {
			fields[ann.Value+"-duration"] = *ann.Duration
		}
	}

	for _, ann := range span.BinaryAnnotations {
		// TODO: TMemoryBuffer embeds *bytes.Buffer so we may not even need this Copy.
		key, value, e := convertBinaryAnnotationFromThrift(ann)
		if e != nil {
			return nil, e
		}

		fields[key] = value
	}

	// computes client duration if cs and cr values present.
	if cs, ok := fields.CS(); ok {
		if cr, ok := fields.CR(); ok {
			fields["cd"] = cr - cs
		}
	}

	// computes server duration if ss and sr values present.
	if ss, ok := fields.SS(); ok {
		if sr, ok := fields.SR(); ok {
			fields["sd"] = ss - sr
		}
	}

	return fields, nil
}
