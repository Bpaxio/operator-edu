package controller

import (
	"bpax.io/ru/cmx/edu/MyApp/pkg/controller/myapp"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, myapp.Add)
}
