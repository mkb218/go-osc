package osc

// #cgo LDFLAGS: -llo
// #include "lo/lo.h"
import "C"

import (
    "unsafe"
    "runtime"
    )

/* to implement liblo's high level OSC interface, we need the following:
timetag type
message arg type, can act as int32, int64, float, double, char, unsigned char, uint8[4], and timetag
 */

type Timetag struct {
    uint32 Sec
    uint32 Frac
}

type OscType byte

type Arg interface {
    Type() OscType
    GetChar() uint8
    GetDouble() float64
    GetFloat() float32
    GetInt64() int64
    GetInt32() int32
    GetMidiMsg() uint8[4]
    GetSymbol() string
    GetString() string
    GetTimetag() Timetag
}

const (
    Udp = iota
    Tcp
    Unix
    )

const (
    Int32 = 'i'
    Float = 'f'
    String = 's'
    Blob = 'b'
    Int64 = 'h'
    Timetag = 't'
    Double = 'd'
    Symbol = 'S'
    Char = 'c'
    MidiMsg = 'm'
    True = 'T'
    False = 'F'
    Nil = 'N'
    Infinitum = 'I'
)
    
    

var Now = Timetag{0,1}

/* opaque address type */
type Address struct {
    lo_address unsafe.Pointer
    dead bool
}

func NewAddress(host, port string) (ret *Address) {
    ret = new(Address)
    ret.dead = false
    chost := C.CString(host)
    defer C.free(unsafe.Pointer(chost))
    cport := C.CString(port)
    defer C.free(unsafe.Pointer(cport))
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
    ret.lo_address = C.lo_address_new_with_proto(proto, chost, cport)
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
    if (this.dead) {
        panic("Method called on dead object")
    }
    C.lo_address_set_ttl(this.lo_address, ttl)
}

func (this *Address) GetTtl() (ttl int) {
    if (this.dead) {
        panic("Method called on dead object")
    }
    ttl = C.lo_address_get_ttl(this.lo_address)
}

func (this *Address) Errno() int {
    if (this.dead) {
        panic("Method called on dead object")
    }
    return C.lo_address_errno(this.lo_address)
}

func (this *Address) Errstr() string {
    if (this.dead) {
        panic("Method called on dead object")
    }
    return C.GoString(C.lo_address_errstr(this.lo_address))
}

func (this *Address) Free() {
    if (this.dead) {
        panic("Method called on dead object")
    }
    C.lo_address_free(this.lo_address)
    this.lo_address = nil
    this.dead = true
}

/* Why go through a bunch of junk to use the lo blob type? Just make a byte slice and call lo_blob_new when we add it to a message */
type Blob []byte

type Message struct {
    lo_message unsafe.Pointer
    dead bool
}

func NewMessage() (ret *Message) {
    ret.dead = false
    lo_message = C.lo_message_new()
    runtime.SetFinalizer(ret, (*Message).Free)
    return
}

func (this *Message) Free() {
    this.dead = true
    C.lo_message_free(this.lo_message)
}

func (this *Message) Add(arg Arg) int {
    
 
func Send(targ Address, path string, args ...Arg) ret int {
    msg := NewMessage(args)
}    
    
type Bundle struct {
    lo_address unsafe.Pointer
}