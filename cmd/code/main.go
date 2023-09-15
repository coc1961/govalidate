package main

import "fmt"

func main() {
	test2()
}

var e = fmt.Errorf("Pepe1")

func test1() {
	var e = fmt.Errorf("Pepe2")
	_ = e
}

func test2() error {
	return fmt.Errorf("Pepe3")
}
