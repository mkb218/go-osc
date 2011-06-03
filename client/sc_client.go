package main

import "osc"
import "fmt"

func main() {
	i := "127.0.0.1"
	p := "57120"
	a := osc.NewAddress(&i, &p)

	m := make(osc.Message, 0)
	m = append(m, osc.StringType("drum"))
	m = append(m, osc.Int32Type(442))
	m = append(m, osc.DoubleType(1.0))
	fmt.Printf("send returned %d\n", m.Send(a, "/note"))
}
