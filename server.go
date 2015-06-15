package zipkin

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	zkcol "github.com/mattkanwisher/distributedtrace/gen/zipkincollector"
)

// Server defines the basic ZipKin server interface.
type Server interface {
	Start() error
	Stop() error
}

type serverImpl struct {
	*thrift.TSimpleServer
	stop chan bool

	collector *Collector
	combiner  *Combiner
	output    Output
}

// NewServer() creates a new server with the specified configuration, or the default
// configuration if nil is given.
func NewServer(config *Config, output Output) (Server, error) {
	config = fillDefaultConfig(config)
	transport, e := thrift.NewTServerSocket(config.ListenAddress)
	if e != nil {
		return nil, e
	}

	var (
		collector = NewCollector(config)
		combiner  = NewCombiner(config)
		processor = zkcol.NewZipkinCollectorProcessor(collector)
		factory   = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		protocol  = thrift.NewTBinaryProtocolFactoryDefault()
		server    = thrift.NewTSimpleServer4(processor, transport, factory, protocol)
		impl      = &serverImpl{server, nil, collector, combiner, output}
	)

	return impl, nil
}

func (s *serverImpl) Start() error {
	s.stop = make(chan bool)
	go s.spanPump()
	go s.combinePump()

	return s.TSimpleServer.Serve()
}

func (s *serverImpl) Stop() error {
	e := s.TSimpleServer.Stop()
	if e != nil {
		return e
	}

	// stop 2 pumps
	s.stop <- true
	s.stop <- true
	return nil
}

func (s *serverImpl) spanPump() {
	for {
		select {
		case <-s.stop:
			return
		case span := <-s.collector.Receive():
			s.combiner.Send(span)
		}
	}
}

func (s *serverImpl) combinePump() {
	for {
		select {
		case <-s.stop:
			return
		case result := <-s.combiner.Receive():
			s.output.Write(result)
		}
	}
}
