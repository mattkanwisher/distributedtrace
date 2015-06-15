package zipkin_test

import (
	"net"
	"testing"
	"time"

	. "github.com/mattkanwisher/distributedtrace"
	a "github.com/stretchr/testify/assert"
)

//func DialTimeout(network, address string, timeout time.Duration) (Conn, error)
const TestAddress = "0.0.0.0:8123"

func TestServer_Serve(t *testing.T) {
	server, e := NewServer(&Config{ListenAddress: TestAddress})
	a.NoError(t, e)

	done := make(chan bool)
	go func() {
		a.NoError(t, server.Start())
		done <- true
	}()

	go func() {
		// since server.Start() blocks (no way to signal successful start)
		// we'll just have to wait here.
		time.Sleep(1 * time.Second)

		conn, e := net.DialTimeout("tcp", TestAddress, 1*time.Second)
		a.NoError(t, e)

		a.NoError(t, conn.Close())
		a.NoError(t, server.Stop())
		done <- true
	}()

	<-done
	<-done
}
