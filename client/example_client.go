package main



import (
	"osc"
	"flag"
	"fmt"
	"unsafe"
)

var testdata string = "ABCDE"

func main() {
	btest := osc.Blob([]byte(testdata))
	btest = append(btest, '\000')

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
	    m := make(osc.Message, 0)
	    m = append(m, 0.12345678)
	    C.lo_message_add_float(m, _Ctype_float(23.0))
		if C.lo_send_message(t, C.CString("/foo/bar"), m) == -1 {
			fmt.Printf("OSC error %d: %s\n", C.lo_address_errno(t), C.GoString(C.lo_address_errstr(t)))
		}
		C.lo_message_free(m)

	    m = C.lo_message_new()
		/* send a message to /a/b/c/d with a mixtrure of float and string
		 * arguments */
		 C.lo_message_add_string(m, C.CString("one"))
 	    C.lo_message_add_float(m, _Ctype_float(0.12345678))
		 C.lo_message_add_string(m, C.CString("three"))
 	    C.lo_message_add_float(m, _Ctype_float(0.00000023001))
 	    C.lo_message_add_float(m, _Ctype_float(1.0))
		C.lo_send_message(t, C.CString("/a/b/c/d"), m)
		C.lo_message_free(m)

	    m = C.lo_message_new()
		/* send a 'blob' object to /a/b/c/d */
		C.lo_message_add_blob(m, btest)
		C.lo_send_message(t, C.CString("/a/b/c/d"), m)
		C.lo_message_free(m)

		/* send a jamin scene change instruction with a 32bit integer argument */
	    m = C.lo_message_new()
	    C.lo_message_add_int32(m, _Ctypedef_int32_t(2))
		C.lo_send_message(t, C.CString("/jamin/scene"), m)
		C.lo_message_free(m)
	}
}
