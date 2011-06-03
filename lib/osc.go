package osc

// #cgo LDFLAGS: -llo
// #include <stdlib.h>
// #include "lo/lo.h"
import "C"

import (
	"unsafe"
	//    "fmt"
	"runtime"
)

/* to implement liblo's high level OSC interface, we need the following:
timetag type
message arg type, can act as int32, int64, float, double, char, unsigned char, uint8[4], and timetag
*/

type OscType byte

type Arg interface {
	GetType() OscType
	GetValue() interface{}
}

type SymbolType string

func (this SymbolType) GetType() OscType {
	return Symbol
}

func (this SymbolType) GetValue() interface{} {
	return this
}

type InfinitumType struct{}

func (this InfinitumType) GetType() OscType {
	return Infinitum
}

func (this InfinitumType) GetValue() interface{} {
	return nil
}

type Int32Type int32

func (this Int32Type) GetType() OscType {
	return Int32
}

func (this Int32Type) GetValue() interface{} {
	return this
}

type FloatType float32

func (this FloatType) GetType() OscType {
	return Float
}

func (this FloatType) GetValue() interface{} {
	return this
}

type Blob []byte

func (this Blob) GetType() OscType {
	return BlobCode
}

func (this Blob) GetValue() interface{} {
	return this
}

type Int64Type int64

func (this Int64Type) GetType() OscType {
	return Int64
}

func (this Int64Type) GetValue() interface{} {
	return this
}

type Timetag struct {
	Sec, Frac C.uint32_t
}

func (this Timetag) GetType() OscType {
	return TimetagCode
}

func (this Timetag) GetValue() interface{} {
	return this
}

type StringType string

func (this StringType) GetType() OscType {
	return String
}

func (this StringType) GetValue() interface{} {
	return this
}

type CharType byte

func (this CharType) GetType() OscType {
	return Char
}

func (this CharType) GetValue() interface{} {
	return this
}

type DoubleType float64

func (this DoubleType) GetType() OscType {
	return Double
}

func (this DoubleType) GetValue() interface{} {
	return this
}

type MidiMsg [4]uint8

func (this MidiMsg) GetType() OscType {
	return MidiMsgCode
}

func (this MidiMsg) GetValue() interface{} {
	var out uint32
	for i, r := range this {
		out = out | (uint32(r) << uint(3-i))
	}
	return out
}

const (
	Udp = iota
	Tcp
	Unix
)

const (
	Int32       OscType = 'i'
	Float       OscType = 'f'
	String      OscType = 's'
	BlobCode    OscType = 'b'
	Int64       OscType = 'h'
	TimetagCode OscType = 't'
	Double      OscType = 'd'
	Symbol      OscType = 'S'
	Char        OscType = 'c'
	MidiMsgCode OscType = 'm'
	True        OscType = 'T'
	False       OscType = 'F'
	Nil         OscType = 'N'
	Infinitum   OscType = 'I'
)

var Now = Timetag{0, 1}

/* opaque address type */
type Address struct {
	lo_address C.lo_address
	dead       bool
}

func NewAddress(host, port *string) (ret *Address) {
	ret = new(Address)
	ret.dead = false
	var chost, cport *C.char = nil, nil
	if host != nil {
		chost = C.CString(*host)
		defer C.free(unsafe.Pointer(chost))
	}
	if port != nil {
		cport = C.CString(*port)
		defer C.free(unsafe.Pointer(cport))
	}
	ret.lo_address = C.lo_address_new(chost, cport)
	runtime.SetFinalizer(ret, (*Address).Free)
	return
}

func NewAddressWithProto(proto int, host, port string) (ret *Address) {
	ret = new(Address)
	ret.dead = false
	chost := C.CString(host)
	defer C.free(unsafe.Pointer(chost))
	cport := C.CString(port)
	defer C.free(unsafe.Pointer(cport))
	ret.lo_address = C.lo_address_new_with_proto(C.int(proto), chost, cport)
	runtime.SetFinalizer(ret, (*Address).Free)
	return
}

func NewAddressFromUrl(url string) (ret *Address) {
	ret = new(Address)
	ret.dead = false
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))
	ret.lo_address = C.lo_address_new_from_url(curl)
	runtime.SetFinalizer(ret, (*Address).Free)
	return
}

func (this *Address) SetTtl(ttl int) {
	if this.dead {
		panic("Method called on dead object")
	}
	C.lo_address_set_ttl(this.lo_address, C.int(ttl))
}

func (this *Address) GetTtl() (ttl int) {
	if this.dead {
		panic("Method called on dead object")
	}
	ttl = int(C.lo_address_get_ttl(this.lo_address))
	return
}

func (this *Address) Errno() int {
	if this.dead {
		panic("Method called on dead object")
	}
	return int(C.lo_address_errno(this.lo_address))
}

func (this *Address) Errstr() string {
	if this.dead {
		panic("Method called on dead object")
	}
	return C.GoString(C.lo_address_errstr(this.lo_address))
}

func (this *Address) Free() {
	if this.dead {
		return
	}
	C.lo_address_free(this.lo_address)
	this.lo_address = nil
	this.dead = true
}

/* why use the message type before it's time to send it? */
type Message []Arg

func (this Message) Send(targ *Address, path string) (ret int) {
	//    fmt.Printf("send %v\n", this)
	ret = this.SendTimestamped(targ, Now, path)
	return
}

func (this Message) SendTimestamped(targ *Address, time Timetag, path string) (ret int) {
	// build a new lo_message
	//    fmt.Printf("send timestamped %v %v\n", time, this)
	m, ret := this.build_lo_message()
	if m == nil {
		//        fmt.Printf("m nil ret %d\n", ret)
		return
	}
	if ret < 0 {
		C.lo_message_free(m)
		//        fmt.Printf("ret negative %d\n", ret)
		return
	}
	defer C.lo_message_free(m)
	s := C.CString(path)
	defer C.free(unsafe.Pointer(s))
	ret = int(C.lo_send_message(targ.lo_address, s, m))
	//    fmt.Printf("ret %d\n", ret)
	return
}

func (this Message) build_lo_message() (m C.lo_message, ret int) {
	m = C.lo_message_new()
	if m == nil {
		ret = -1
		return
	}
	for _, arg := range this {
		//        fmt.Printf("%d %s %v\n", n, arg.GetType(), arg.GetValue())
		switch arg.GetType() {
		case Int32:
			ret = int(C.lo_message_add_int32(m, C.int32_t(arg.GetValue().(Int32Type))))
		case Float:
			ret = int(C.lo_message_add_float(m, C.float(arg.GetValue().(FloatType))))
		case BlobCode:
			a, i := arg.GetValue().(Blob)
			if !i {
				ret = -2
				break
			}

			b := C.lo_blob_new(C.int32_t(len(a)), unsafe.Pointer(&(a[0])))
			if b == nil {
				ret = -3
				break
			}
			defer C.lo_blob_free(C.lo_blob(b))
			ret = int(C.lo_message_add_blob(m, b))
		case Int64:
			ret = int(C.lo_message_add_int64(m, C.int64_t(arg.GetValue().(Int64Type))))
		case TimetagCode:
			ret = int(C.lo_message_add_timetag(m, C.lo_timetag{arg.GetValue().(Timetag).Sec, arg.GetValue().(Timetag).Frac}))
		case Double:
			ret = int(C.lo_message_add_double(m, C.double(arg.GetValue().(DoubleType))))
		case Symbol:
			s := C.CString(string(arg.GetValue().(SymbolType)))
			if s == nil {
				ret = -4
				break
			}
			defer C.free(unsafe.Pointer(s))
			ret = int(C.lo_message_add_symbol(m, s))
		case String:
			s := C.CString(string(arg.GetValue().(StringType)))
			if s == nil {
				ret = -4
				break
			}
			defer C.free(unsafe.Pointer(s))
			ret = int(C.lo_message_add_string(m, s))
		case Char:
			ret = int(C.lo_message_add_char(m, arg.GetValue().(C.char)))
		case MidiMsgCode:
			mm := arg.GetValue().(MidiMsg)
			ret = int(C.lo_message_add_midi(m, (*C.uint8_t)(unsafe.Pointer(&mm))))
		case True:
			ret = int(C.lo_message_add_true(m))
		case False:
			ret = int(C.lo_message_add_false(m))
		case Nil:
			ret = int(C.lo_message_add_nil(m))
		case Infinitum:
			ret = int(C.lo_message_add_infinitum(m))
		default:
			ret = -5
		}
		if ret < 0 {
			C.lo_message_free(m)
			m = nil
			return
		}
	}
	return
}

type MsgPath struct {
	Path string
	Msg  Message
}

type Bundle struct {
	Time     Timetag
	MsgPaths []MsgPath
}

func (this *Bundle) Send(targ *Address) (ret int) {
	ret = 0
	b := C.lo_bundle_new(C.lo_timetag{this.Time.Sec, this.Time.Frac})
	if b == nil {
		ret = -1
		return
	}
	defer C.lo_bundle_free_messages(b)
	for _, r := range this.MsgPaths {
		p := C.CString(r.Path)
		defer C.free(unsafe.Pointer(p))
		m, ret := r.Msg.build_lo_message()
		if ret < 0 {
			return
		}

		ret = int(C.lo_bundle_add_message(b, p, m))
		if ret < 0 {
			return
		}
	}
	ret = int(C.lo_send_bundle(targ.lo_address, b))
	return
}
