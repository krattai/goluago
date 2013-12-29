package main

import (
	"fmt"
	lua "github.com/akavel/goluago/internal"
	"runtime/debug"
	"unsafe"
)

func newdumper() func(L lua.State, pbuf, size uintptr, userdata interface{}) int {
	x := 0
	return func(L lua.State, pbuf, size uintptr, userdata interface{}) int {
		for i := uintptr(0); i < size; i++ {
			print(fmt.Sprintf("%02x ", *((*byte)(unsafe.Pointer(pbuf + i)))))
			if x%8 == 7 {
				print("\n")
			} else if x%4 == 3 {
				print(" ")
			}
			x++
		}
		return 0
	}
}

func main() {
	defer func() {
		if x := recover(); x != nil {
			fmt.Printf("PANIC: %v\n", x)
			fmt.Printf("panic stacktrace:\n%s", string(debug.Stack()))
		}
	}()

	println("hello wrld")
	s := lua.Open()
	t := s.Gettop()
	println("top=", t)
	s.Pushinteger(5)
	println("push 5")
	println("top=", s.Gettop())
	s.Pushinteger(5)
	println("push 5")
	s.Pushinteger(0)
	println("push 0")
	println("top=", s.Gettop())
	println("equal(-1,-2)=", s.Equal(-1, -2))
	println("equal(-2,-3)=", s.Equal(-2, -3))
	println("equal(1,2)=", s.Equal(1, 2))
	println("equal(2,3)=", s.Equal(2, 3))

	/* FIXME: choose correct chunk depending on architecture?
	r := s.Loadbuffer([]byte(
		// "return 2+3", generated by '../gen_chunk.lua'
		"\x1bLuaQ\x00\x01\x04\x04\x04\x08\x00\x01\x00\x00\x00\x00\x00"+
			"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02\x02\x03\x00\x00\x00"+
			"\x01\x00\x00\x00\x1e\x00\x00\x01\x1e\x00\x80\x00\x01\x00\x00"+
			"\x00\x03\x00\x00\x00\x00\x00\x00\x14@\x00\x00\x00\x00\x03\x00"+
			"\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00"+
			"\x00\x00\x00\x00\x00\x00\x00"),
		"chunk 1")
	*/
	r := s.Loadbuffer([]byte("return 2+3"), "chunk 1")
	if r != 0 {
		panic(r)
	}

	println("dump:")
	s.Dump(newdumper(), nil)
	println()

	s.Call(0, 1)
	println("call")
	println("top=", s.Gettop())
	println("equal(-1,1)=", s.Equal(-1, 1)) // expected: yes (i.e. peek()==5)

	s.Pushlstring([]byte("foobar"))
	println("pushlstring")
	println("top=", s.Gettop())
	println("tolstring(-1)=", string(s.Tolstring(-1)))

	println("tolstring(-2)=", string(s.Tolstring(-2)))
	println("tonumber(-2)=", s.Tonumber(-2))

	println()
	prog := "a=8 return a"
	println("load '" + prog + "'")
	r = s.Loadbuffer([]byte(prog), "chunk 2")
	if r != 0 {
		println("err: ", string(s.Tolstring(-1)))
		panic(r)
	}
	println("dump:")
	r = s.Dump(newdumper(), nil)
	println()
	if r != 0 {
		println("err: ", string(s.Tolstring(-1)))
		panic(r)
	}
	s.Call(0, 1)
	println("call")
	println("top=", s.Gettop())
	println("tolstring(-1)=", string(s.Tolstring(-1)))
	println()

	s.Pushgofunction(func(l lua.State) int32 {
		println("hello goluago callback world!")
		return 0
	})
	s.Call(0, 0)
	println()

	//prog = " return {...} "
	prog = " return {...} "
	println("load '" + prog + "'")
	r = s.Loadbuffer([]byte(prog), "chunk x")
	if r != 0 {
		println("err: ", string(s.Tolstring(-1)))
		panic(r)
	}
	println("dump:")
	r = s.Dump(newdumper(), nil)
	println()
	if r != 0 {
		println("err: ", string(s.Tolstring(-1)))
		panic(r)
	}
	s.Call(0, 0)
	println("call")
	println("top=", s.Gettop())
	println("tolstring(-1)=", string(s.Tolstring(-1)))
	println()
	

	println()
	prog = "b=1; local aaa='a'; x=aaa+b"
	/* // same as code above, only compiled with normal lua and dumped
		prog = "\x1b\x4c\x75\x61\x51\x00\x01\x04\x04\x04\x08\x00\x01\x00\x00\x00\x00\x00\x00\x00"+
	"\x00\x00\x00\x00\x00\x00\x00\x02\x02\x07\x00\x00\x00\x01\x40\x00\x00\x07\x00\x00"+
	"\x00\x01\x80\x00\x00\x45\x00\x00\x00\x4c\x40\x00\x00\x47\xc0\x00\x00\x1e\x00\x80"+
	"\x00\x04\x00\x00\x00\x04\x02\x00\x00\x00\x62\x00\x03\x00\x00\x00\x00\x00\x00\xf0"+
	"\x3f\x04\x02\x00\x00\x00\x61\x00\x04\x02\x00\x00\x00\x78\x00\x00\x00\x00\x00\x07"+
	"\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x01"+
	"\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x04\x00\x00\x00\x61"+
	"\x61\x61\x00\x03\x00\x00\x00\x06\x00\x00\x00\x00\x00\x00\x00"
	*/
	println("load '" + prog + "'")
	r = s.Loadbuffer([]byte(prog), "chunk 3")
	if r != 0 {
		println("err: ", string(s.Tolstring(-1)))
		panic(r)
	}
	println("dump:")
	r = s.Dump(newdumper(), nil)
	println()
	if r != 0 {
		println("err: ", string(s.Tolstring(-1)))
		panic(r)
	}
	s.Call(0, 0)
	println("call")
	println("top=", s.Gettop())
	println("tolstring(-1)=", string(s.Tolstring(-1)))
	println()
	
	// PANIC and fprintf(stdio,...) test
	println("lessthan(-2,-1)=", s.Lessthan(-3, -1))
}
