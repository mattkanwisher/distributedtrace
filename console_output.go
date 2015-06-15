package zipkin

import (
	"encoding/json"
	"os"
)

type consoleOutput struct{}

// ConsoleOutput() returns an Output that encods ZipKin spans to JSON and writes it to
// standard output.
func ConsoleOutput() Output {
	return (*consoleOutput)(nil)
}

func (*consoleOutput) Write(result OutputMap) (e error) {
	defer autoRecover(&e)

	buffer, e := json.Marshal(result)
	noError(e)
	_, e = os.Stdout.Write(buffer)
	noError(e)
	_, e = os.Stdout.Write([]byte("\n"))
	noError(e)
	return nil
}
