package main

import (
	"fmt"
	"testing"
)

func Test_test1(t *testing.T) {
	var e = fmt.Errorf("Pepe2")
	_ = e
}
