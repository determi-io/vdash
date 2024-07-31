package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	go_agda_wrapper "github.com/determi-io/go-agda-wrapper"
)

/////////////////////////////////////////
// Types

type FragmentValue struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type FragmentDep struct {
	CtorArgs []FragmentValue `json:"CtorArgs"`
}

type Fragment struct {
	Name  string        `json:"Name"`
	Ctors []FragmentDep `json:"Ctors"`
}

type FragmentList struct {
	Fragments []Fragment `json:"Fragments"`
}

/////////////////////////////////////////
// Functions

func myfunc() {
	go_agda_wrapper.CompileAndRun("test.agda", "./src", "./build", "main")
}

func LoadFragments(file string) (FragmentList, error) {

	// create temporary directory (source)
	sourcedir, err := os.MkdirTemp("", "go_agda_wrapper_source_")
	if err != nil {
		log.Printf("Could not create temp dir: %s", err)
		return FragmentList{}, err
	}

	// create temporary directory (target)
	targetdir, err := os.MkdirTemp("", "go_agda_wrapper_target_")
	if err != nil {
		log.Printf("Could not create temp dir: %s", err)
		return FragmentList{}, err
	}

	// copy our source file to target directory
	_, err = exec.Command("cp", file, sourcedir).Output()
	if err != nil {
		return FragmentList{}, err
	}

	// call agda to evaluate file
	output, err := go_agda_wrapper.CompileAndRun(sourcedir, filepath.Base(file), targetdir, "output")
	if err != nil {
		return FragmentList{}, err
	}

	// logging
	println("got js output: ", output)

	// unmarshalling
	var fragmentList FragmentList
	err = json.Unmarshal([]byte(output), &fragmentList)
	if err != nil {
		return FragmentList{}, err
	}

	fmt.Println("marshalled: ", fragmentList)

	return fragmentList, nil
}
