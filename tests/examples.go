//go:generate goderive
package tests

import (
	"net/http"
	t "time"

	"github.com/nextzhou/goderive/plugin"
)

// derive-set
// derive-slice
type Int = int

// derive-set:Rename=intOrderSet;Order=Append
type Int2 = int

// derive-set:Order=Key
type Int3 = int

// unexported type, from imported package
// derive-set
// derive-slice
type h = http.Handler

// from renamed imported package
// derive-set
type T = t.Time

// derive-set:Order=Key
type S = string

// from remote package
// derive-set: Export
type p = plugin.Plugin

// from this package
// derive-set: !Export
type MyType struct {
	Field1 string
	field2 bool
}

// derive-access: Receiver=c
type c struct {
	abc          string                                              // base type
	Def          *int                                                // pointer type
	hi, jk, lmn  struct{ a string }                                  // anonymous struct is unsupported
	a            t.Time                                              // selector expr
	b            []string                                            // slice type
	bb           [3]string                                           // array type
	c            map[int]string                                      // map type
	d            chan int                                            // channel type
	e            chan<- int                                          // write-only channel type
	f            <-chan int                                          // read-only channel type
	http.Request                                                     // anonymous field
	ff           func(int, string, c, d t.Time) (a, b bool, e error) // complex function type
}

// derive-access: Receiver=b
type b struct {
	c *c
	C *c
}

// derive-access: Receiver=a
type AA struct {
	b *b
	B *b
	c *c
	C *c
}
