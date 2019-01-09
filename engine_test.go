package smart

import (
	"fmt"
	"testing"
)

func TestSmartEngine_StartInstanceById(t *testing.T) {
	e := NewSmartEngine()
	fmt.Println(e.StartInstanceById(1))
}