//go:generate goderive
package tests

import (
	"net/http"
	t "time"
)

// derive-set
type Int = int

// derive-set:Rename=intOrderSet;Order=Append
type Int2 = int

// derive-set:Order=Key
type Int3 = int

// unexported type, from imported package
// derive-set
type h = http.Handler

// from rename imported package
// derive-set
type T = t.Time

// from this package
// derive-set
type A struct{ s string }
