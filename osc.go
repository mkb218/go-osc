package osc

// #cgo LDFLAGS: -llo
// #include <stdlib.h>
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
    Sec, Frac _Ctypedef_uint32_t
}

type MidiMsg [4]uint8

type OscType byte

type Arg interface {
    GetType() OscType
    GetValue() interface{}
}

type SymbolType string

func (this *SymbolType) GetType() OscType {
    return Symbol
}

func (this *SymbolType) GetValue() (interface{}) {
    return *this
}    

type InfinitumType struct{}

func (this *InfinitumType) GetType() OscType {
    return Infinitum
}

func (this *InfinitumType) GetValue() interface{} {
    return nil
}    

type goTypeWrapper struct {
    Arg
    value interface{} 
}

func (this *goTypeWrapper) GetType() (o OscType) {
    switch i := this.value.(type) {
    case int32:
        o = Int32
    case float32:
        o = Float
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
        if i, _ := this.value.(bool); i {
            o = True
        } else {
            o = False
        }
    case nil:
        o = Nil
        // no way to detect infinitum
    }
    return
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

var Now = Timetag{0,1}

/* opaque address type */
type Address struct {
    lo_address _Ctypedef_lo_address
    dead bool
}

func NewAddress(host, port *string) (ret *Address) {
    ret = new(Address)
    ret.dead = false
    chost := C.CString(*host)
    defer C.free(unsafe.Pointer(chost))
    cport := C.CString(*port)
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
    if (this.dead) {
        panic("Method called on dead object")
    }
    C.lo_address_set_ttl(this.lo_address, C.int(ttl))
}

func (this *Address) GetTtl() (ttl int) {
    if (this.dead) {
        panic("Method called on dead object")
    }
    ttl = int(C.lo_address_get_ttl(this.lo_address))
    return
}

func (this *Address) Errno() int {
    if (this.dead) {
        panic("Method called on dead object")
    }
    return int(C.lo_address_errno(this.lo_address))
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
 
func (this *Message) Send(targ *Address, path string) (ret int) {
    ret = this.SendTimestamped(targ, Now, path)
    return
}

func (this *Message) SendTimestamped(targ *Address, time Timetag, path string) (ret int) {
    // build a new lo_message
    m, ret := this.build_lo_message()
    if m == nil {
        return
    }
    if ret < 0 {
        C.lo_message_free(m)
        return
    }
    defer C.lo_message_free(m)
    s := C.CString(path)
    defer C.free(unsafe.Pointer(s))
    ret = int(C.lo_send_message(targ.lo_address, s, m))
    return
}

func (this *Message) build_lo_message() (m C.lo_message, ret int) {
    m = C.lo_message_new()
    if m == nil { 
        ret = -1
        return
    }
    for _, arg := range *this {
        switch arg.GetType() {
        case Int32:
            ret = int(C.lo_message_add_int32(m, _Ctypedef_int32_t(arg.GetValue().(int32))))
        case Float:
            ret = int(C.lo_message_add_float(m, _Ctype_float(arg.GetValue().(float32))))
        case BlobCode:
            a, i := arg.GetValue().(Blob)
            if !i {
                ret = -1
                break
            }
            
            b := C.lo_blob_new(_Ctypedef_int32_t(len(a)), unsafe.Pointer(&a))
            if b == nil {
                ret = -1
                break
            }
            defer C.lo_blob_free(_Ctypedef_lo_blob(b))
            ret = int(C.lo_message_add_blob(m, b))
        case Int64:
            ret = int(C.lo_message_add_int64(m, _Ctypedef_int64_t(arg.GetValue().(int64))))
        case TimetagCode:
            ret = int(C.lo_message_add_timetag(m, C.lo_timetag{arg.GetValue().(Timetag).Sec, arg.GetValue().(Timetag).Frac}))
        case Double:
            ret = int(C.lo_message_add_double(m, _Ctype_double(arg.GetValue().(float64))))
        case Symbol: 
            s := C.CString(arg.GetValue().(string))
            if s == nil {
                ret = -1
                break
            }
            defer C.free(unsafe.Pointer(s))
            ret = int(C.lo_message_add_symbol(m, s))
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
            ret = -1
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
    Msg *Message
}
    
type Bundle struct {
    Time Timetag
    MsgPaths []MsgPath
}

func (this *Bundle) Send(targ Address) (ret int) {
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