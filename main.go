package main

import (
	"fmt"
	"github.com/PtitLuca/go-dispatcher/dispatcher"
	"log"
)

type T struct {
}

func (t *T) ExampleVariadic(a int, b ...string) int {
	return a + len(b)
}

func main() {
	d := dispatcher.New()
	err := d.Register("Test", &T{})
	if err != nil {
		log.Fatalln(err)
	}

	output, err := d.Run("Test", "ExampleVariadic", 1, "This", "Are", "Variadic", "Arguments")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(output[0].Int())
}
