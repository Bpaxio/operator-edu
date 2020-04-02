package controller

import (
	"bpax.io/ru/cmx/edu/MyOperators/backuper/pkg/controller/backuper"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, backuper.Add)
}
