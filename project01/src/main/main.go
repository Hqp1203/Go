package main

import (
	"fmt"
	"project01/src/model"
)

func main() {
	var a int = 5
	var p *int = &a
	*p = 27
	fmt.Printf("%v %v %s", a, &p,model.Name)
}
