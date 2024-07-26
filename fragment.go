package main

import (
	go_agda_wrapper "github.com/determi-io/go-agda-wrapper"
)

func myfunc() {
	go_agda_wrapper.CompileAndRun("test.agda", "./src", "./build", "main")
}
