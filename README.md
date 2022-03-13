# go-dispatcher

A single-file dispatcher written in Golang.

## Description

Inspired by the JSON-RPC module of [polygon-edge](https://github.com/0xPolygon/polygon-edge) and [geth](https://github.com/ethereum/go-ethereum), this package provides a way to register methods from a struct pointer.

It can be used to bootstrap a server or create a worker pool in your project !

## Installation

Run the following command :
````shell
go get github.com/PtitLuca/go-dispatcher@v1.0.2
````

## Quick Start

<details>
<summary>A very simple example</summary>

````go
package main

import (
	"fmt"
	"github.com/PtitLuca/go-dispatcher/dispatcher"
	"log"
)

type T struct {
}

func (t *T) Example(a, b int) int {
	return a + b
}

func main() {
	d := dispatcher.New()
	err := d.Register("Test", &T{})
	if err != nil {
		log.Fatalln(err)
	}

	output, err := d.Run("Test", "Example", 1, 2)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(output[0].Int())
}
````

After running, you should get this in your terminal :
````shell
3
````

</details>

<details>
<summary>Multi-service registration</summary>

````go
package main

import (
	"fmt"
	"github.com/PtitLuca/go-dispatcher/dispatcher"
	"log"
)

type X struct {
}

func (x *X) Example2(a, b string) string {
	return a + b
}

type T struct {
}

func (t *T) Example(a, b int) int {
	return a + b
}

func main() {
	d := dispatcher.New()
	err := d.Register("Test", &T{})
	if err != nil {
		log.Fatalln(err)
	}

	err = d.Register("TestX", &X{})
	if err != nil {
		log.Fatalln(err)
	}

	output, err := d.Run("Test", "Example", 1, 2)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(output[0].Int())

	output, err = d.Run("TestX", "Example2", "Hello", "World")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(output[0].String())
}
````

After running, you should get this in your terminal :
````shell
3
HelloWorld
````

</details>

<details>
<summary>Variadic arguments</summary>

````go
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

	output, err := d.Run("Test", "ExampleVariadic", 1, "These", "Are", "Variadic", "Arguments")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(output[0].Int())
}
````

After running, you should get this in your terminal :
````shell
5
````

</details>

## Features

This package has support for :

- Multi-service registration
- Methods registration - exported only
- Variadic arguments

## Authors

- [PtitLuca](https://github.com/PtitLuca)
