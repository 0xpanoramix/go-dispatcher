# go-dispatcher

A single-file dispatcher written in Golang.

## Description

Inspired by the JSON-RPC module of [geth](https://github.com/ethereum/go-ethereum), this package provides a way to 
register methods from a struct pointer.

It can be used to bootstrap a server or create a worker pool in your project !

## Installation

Run the following command :
````shell
go get github.com/PtitLuca/go-dispatcher
````

## Quick Start

Here's an example of the dispatcher's usage :
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

## Features

## Authors

- [PtitLuca](https://github.com/PtitLuca)