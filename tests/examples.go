//go:generate goderive
package examples

// derive-set
type Int = int

// derive-set:Rename=intOrderSet;Order=Append
type Int2 = int
