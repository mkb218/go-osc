package main

import (
	"osc"
	"flag"
	"fmt"
)

var testdata string = "ABCDE"

func main() {
	btest := osc.Blob([]byte(testdata))

	p := "7770"
	t := osc.NewAddress(nil, &p)

	quit := flag.Bool("q", false, "")

	flag.Parse()

	if *quit {
		m := make(osc.Message, 0)
		if m.Send(t, "/quit") == -1 {
			fmt.Printf("OSC error %d: %s\n", t.Errno(), t.Errstr())
		}
	} else {
		m := make(osc.Message, 0, 2)
		m = append(m, osc.FloatType(0.12345678))
		m = append(m, osc.FloatType(23.0))
		if m.Send(t, "/foo/bar") == -1 {
			fmt.Printf("OSC error %d: %s\n", t.Errno(), t.Errstr())
		}

		m = make(osc.Message, 0, 5)
		/* send a message to /a/b/c/d with a mixture of float and string
		 * arguments */
		m = append(m, osc.StringType("one"))
		m = append(m, osc.DoubleType(0.12345678))
		m = append(m, osc.StringType("three"))
		m = append(m, osc.FloatType(0.00000023001))
		m = append(m, osc.FloatType(1.0))
		m.Send(t, "/a/b/c/d")

		m = make(osc.Message, 0, 1)
		/* send a 'blob' object to /a/b/c/d */
		m = append(m, btest)
		m.Send(t, "/a/b/c/d")

		/* send a jamin scene change instruction with a 32bit integer argument */
		m = make(osc.Message, 0, 1)
		m = append(m, osc.Int32Type(2))
		m.Send(t, "/jamin/scene")
	}
}
