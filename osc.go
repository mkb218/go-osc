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

type MidiMsg uint8[4]

type OscType byte

type Arg interface {
    GetType() OscType
    GetValue() interface{}
}

type SymbolType string

func (this *SymbolType) GetType() OscType {
    return Symbol
}

func (this *SymbolType) GetValue interface{} {
    return *this
}    

type InfinitumType interface{}

func (this *InfinitumType) GetType() OscType {
    return Infinitum
}

func (this *InfinitumType) GetValue interface{} {
    return nil
}    

type goTypeWrapper struct {
    Arg
    value interface{} 
}

func (this *goTypeWrapper) GetType() o OscType {
    switch i := (*this).(type) {
    case int32:
        o = Int32
    case float32:
        o = Float32
    case string:
        o = String
    case Blob:
        o = BlobCode
    case int64:
        o = Int64
    case Timetag:
        o = TimetagCode
    case float64:
        o = Double
        // no way to automatically detect symbols!
    case byte:
        o = Char
    case MidiMsg:
        o = MidiMsgCode
    case bool:
        if bool(*this) {
            o = True
        } else {
            o = False
        }
    case nil:
        o = Nil
        // no way to detect infinitum
    }
}

const (
    Udp = iota
    Tcp
    Unix
    )

const (
    Int32 OscType = 'i'
    Float OscType = 'f'
    String OscType = 's'
    BlobCode OscType = 'b'
    Int64 OscType = 'h'
    TimetagCode OscType = 't'
    Double OscType = 'd'
    Symbol OscType = 'S'
    Char OscType = 'c'
    MidiMsgCode OscType = 'm'
    True OscType = 'T'
    False OscType = 'F'
    Nil OscType = 'N'
    Infinitum OscType = 'I'
)

const Now = Timetag{0,1}

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
        return
    }
    C.lo_address_free(this.lo_address)
    this.lo_address = nil
    this.dead = true
}

/* Why go through a bunch of junk to use the lo blob type? Just make a byte slice and call lo_blob_new when we add it to a message */
type Blob []byte

/* why use the message type before it's time to send it? */
type Message []Arg
 
func (this *Message) Send(targ Address, path string) ret int {
    ret = this.SendTimestamped(targ, Now, path)
}

func (this *Message) SendTimestamped(targ Address, time Timetag, path string) ret int {
    // build a new lo_message
    m := C.lo_message_new()
    defer C.lo_message_free(unsafe.Pointer(m))
    for _, arg := range m {
        switch arg.GetType() {
        case Int32:
            C.lo_message_add_int32(m, int32(arg.GetValue()))
        case Float:
            C.lo_message_add_float(m, float32(arg.GetValue()))
        case Blob:
            b := C.lo_message_blob_new(len(Blob(arg)), unsafe.Pointer(arg))
            defer C.lo_blob_free(unsafe.Pointer(b))
            C.lo_message_add_blob(m, b)
        case:
            
    }
}
    
type Bundle struct {
    lo_address unsafe.Pointer
}