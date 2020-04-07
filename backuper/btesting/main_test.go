package btesting

import (
	"testing"

	optest "github.com/operator-framework/operator-sdk/pkg/test"
)

func TestMain(m *testing.M) {
	optest.MainEntry(m)
}
