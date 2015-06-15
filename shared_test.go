package zipkin_test

func noError(e error) {
	if e != nil {
		panic(e)
	}
}
