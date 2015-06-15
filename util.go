package zipkin

import (
	"fmt"
)

func autoRecover(e *error) {
	if r := recover(); r != nil {
		switch r_ := r.(type) {
		case error:
			*e = r_
		case string:
			*e = fmt.Errorf(r_)
		default:
			*e = fmt.Errorf("%s", r_)
		}
	}
}

func noError(e error) {
	if e != nil {
		panic(e)
	}
}
