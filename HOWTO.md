# ZIPKIN CUSTOM SERVER IN GO

This post will explain how to build a custom [Twitter ZipKin][0] server that pipes data
into [Redis][5]. The protocol part is easily done due to the fact that ZipKin uses [Apache
Thrift][1] which is a common protocol standard with tools for generating clients and
servers automatically from Thrift files. Since ZipKin is open source, the Thrift protocol
file is also avialable on its GitHub repo.

Before getting started, you will need to install the Thrift's tool executable first. If
you are on Mac OS X, just go `brew install thrift` or if you are on linux, use whatever
package manager you have to install it.

First, we will need the Thrift protocol files which we can download from the
[zipkin-thrift][2] folder on the [ZipKin GitHub repo][3]. Let's also put them into a
folder.

```sh
$ mkdir -p zk
$ cd zk
$ wget -nv \
  "https://raw.githubusercontent.com/twitter/zipkin/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/scribe.thrift" \
  "https://raw.githubusercontent.com/twitter/zipkin/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/zipkinCollector.thrift" \
  "https://raw.githubusercontent.com/twitter/zipkin/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/zipkinDependencies.thrift" \
  "https://raw.githubusercontent.com/twitter/zipkin/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/zipkinCore.thrift" \
  "https://raw.githubusercontent.com/twitter/zipkin/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin/zipkinQuery.thrift"
```

From this we can start generating service and model code from the specification. Run the
`thrift` tool to do this.

```sh
$ thrift -v -r -out zk -gen go:package_prefix=github.com/mattkanwisher/distributedtrace/zk/,package_name=zk zk/zipkinCollector.thrift
$ thrift -v -r -out zk -gen go:package_prefix=github.com/mattkanwisher/distributedtrace/zk/,package_name=zk zk/zipkinCore.thrift
```

* The `-v` flag specifies verbose mode.
* The `-r` flag makes thrift also generates code for included dependencies as well.
* The `-out zk` part specifies the generated code should goes into the `zk` folder.
* The `-gen go` puts Thrift in code generation mode.
* The `:package_prefix=` part adds specific language-specific metadata to the generator.
  In this case we're specifying the package path and name of the generated code.
* The `zk/zipkinCollector.thrift` is the thrift specificiation we're generating from.

We will need the core ZipKin types as well as the collector interface thus the two files.
After running the generator command, additional Go packages will now shows up under the
`zk/` folder.

Let's now create a Thrift server that talks the ZipKin protocol. First we'll need to
import Apache Thrift's base package and the generated service specification code first.

```go
package main

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	zkcol "github.com/mattkanwisher/distributedtrace/zk/zipkincollector"
)
```

Now to create a server, we'll first need to create a "Processor" type first which
basically processes incoming requests according to the service specification. If you check
the generated `zipkincollector` package, you will find a `ZipKinCollector` interface. This
interface list the methods that must be available according to the Thrift file. Let's make
a type that implement this interface first:

```go
// import "github.com/mattkanwisher/distributedtrace/zk/scribe"
// import "fmt"
type Collector struct { }

var _ zkcol.ZipKinCollector = &Collector{}

func (c *Collector) Log(entries []*scribe.LogEntry) (code scribe.ResultCode, e error) {
	fmt.Println("Log")
  return scribe.ResultCode_OK, nil
}

func (c *Collector) StoreDependencies(dependencies *zkdep.Dependencies) error {
  return nil
}

func (c *Collector) StoreTopAnnotations(serviceName string, annotations []string) error {
  return nil
}

func (c *Collector) StoreTopKeyValueAnnotations(serviceName string, annotations []string) error {
  return nil
}
```

Now that we have the service implemented, we can use Thrift's types to help us construct
a server from the service, first make a "Processor" from our service implementation that
we have just made:

```go
processor := zkcol.NewZipKinCollectorProcessor(&Collector{})
```

From the processor we can now construct a generic Thrift server that pass requests to our
server. First we'll need to allocate the server socket.

```go
transport, e := thrift.NewTServerSocket(":8080")
if e != nil {
  panic(e)
}
```

Don't forget to check the error here in case port 8080 is taken up. Next we'll need to
create a "transport factory" and "protocol factory" that will help with encoding and
decoding things that go through our servers.

```go
var (
	factory   = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocol  = thrift.NewTBinaryProtocolFactoryDefault()
)
```

Now we can create our server from the all the parts above:

```go
server := thrift.NewTSimpleServer4(processor, transport, factory, protocol)
if e := server.Serve(); e != nil {
	panic(e)
}
```

To test this out, you can use the [testclient.go file][4] I have implemented which
basically sends a flood of zipkin spans using the scribe collector to the usual zipkin
port on the local machine.

## EXTRACTING SCRIBE ENTRIES

Now we don't actually get a vanilla ZipKin span in our `Collector` type. Instead the span
is wrapped inside a `*scribe.LogEntry`. This is because our test client sends data via a
a scribe logging protocol. The data is in binary and is base64-encoded as a string in the
log's message.

It's easy to extract the bytes needed from the message by using Go's built-in
`encoding/base64` package:

```go
// import "encoding/base64"

func (c *Collector) Log(entries []*scribe.LogEntry) (code scribe.ResultCode, e error) {
	for _, entry := range entries {
		var buffer []byte
		if buffer, e = base64.StdEncoding.DecodeString(entry.Message); e != nil {
			return scribe.ResultCode_TRY_LATER, e
		}
	}
}
```

Now that we have a `[]byte` buffer for each entry we can read out Thrift data model from
it by creating the Thrift's buffer type, copying data to it and reading it out using
Thrift's binary protocol reader:

```go
// import "https://github.com/mattkanwisher/distributedtrace/blob/master/testclient/main.go"

thriftBuffer := thrift.NewTMemoryBuffer()
if n, e := thriftBuffer.Write(buffer); e != nil {
	return scribe.ResultCode_TRY_LATER, e
} else if n != len(buffer) {
	return scribe.ResultCode_TRY_LATER, fmt.Errorf("buffer copy failure.")
}

protocol := thrift.NewTBinaryProtocol(thriftBuffer, true, true)
span := &zkcore.Span{}
if e := span.Read(protocol); e != nil {
	return scribe.ResultCode_TRY_LATER, e
}

// valid zkcore.Span by this point.
```

Now that we have a ZipKin span instance decoded for us, we can throw this into [Redis][5]
as a temporary storage. First let's create a Redis client:

```go
// import "gopkg.in/redis.v3"
client := redis.NewClient(&redis.Options{
	Addr:     "0.0.0.0",
	Password: "",
	DB:       0,
})
```

Since our Redis driver only accepts JSON, let's encode this
into a JSON string before saving it.

```go
// import "encoding/json"

if buffer, e = json.Marshal(span); e != nil {
	return scribe.ResultCode_TRY_LATER, e
}

if _, e = client.RPush("logs", string(buffer)).Result(); e != nil {
	return scribe.ResultCode_TRY_LATER, e
}

return scribe.ResultCode_OK, nil
```

And that concludes our quick and dirty ZipKin server re-implementation in Go.

[0]: https://github.com/twitter/zipkin
[1]: https://thrift.apache.org
[2]: https://github.com/twitter/zipkin/tree/master/zipkin-thrift/src/main/thrift/com/twitter/zipkin
[3]: https://github.com/twitter/zipkin
[4]: https://github.com/mattkanwisher/distributedtrace/blob/master/testclient/main.go
[5]: http://redis.io

