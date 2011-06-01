package main



import (
	"osc"
	"flag"
	"fmt"
)

var testdata string = "ABCDE"

func main() {
	btest := osc.Blob([]byte(testdata))
	btest = append(btest, '\000')

	p := "7770"
	t := osc.NewAddress(nil, &p)
	fmt.Printf("%s %v\n", p, t)

	quit := flag.Bool("q", false, "")

	flag.Parse()

	if *quit {
	    m := make(osc.Message, 0)
		if m.Send(t, "/quit") == -1 {
			fmt.Printf("OSC error %d: %s\n", t.Errno(), t.Errstr())
		}
	} else {
	    m := new(osc.Message)
	    *m = make(osc.Message, 0)
	    m = m.Add(0.12345678)
	    m = m.Add(23.0)
		if m.Send(t, "/foo/bar") == -1 {
			fmt.Printf("OSC error %d: %s\n", t.Errno(), t.Errstr())
		}

	    *m = make(osc.Message, 0)
		/* send a message to /a/b/c/d with a mixtrure of float and string
		 * arguments */
		m = m.Add("one")
    	m = m.Add(0.12345678)
		m = m.Add("three")
 	    m = m.Add(0.00000023001)
 	    m = m.Add(1.0)
		m.Send(t, "/a/b/c/d")

	    *m = make(osc.Message, 0)
		/* send a 'blob' object to /a/b/c/d */
		m = m.Add(btest)
		m.Send(t, "/a/b/c/d")

		/* send a jamin scene change instruction with a 32bit integer argument */
	    *m = make(osc.Message, 0)
	    m = m.Add(2)
		m.Send(t, "/jamin/scene")
	}
}
