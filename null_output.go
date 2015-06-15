package zipkin

type nullOutput struct{}

// NullOutput() returns an Output that consumes ZipKin spans but otherwise does nothing.
// This is only useful for testing.
func NullOutput() Output {
	return (*nullOutput)(nil)
}

func (*nullOutput) Write(result OutputMap) error {
	return nil
}
