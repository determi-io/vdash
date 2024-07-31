package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func createAndWrite(filename string, content string) (*os.File, error) {
	// create file
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	// write content to file
	_, err = file.WriteString(content)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func TestAll(t *testing.T) {
	// example fragment
	osf := Fragment{
		Name:  "OS",
		Ctors: []FragmentDep{},
	}

	list := FragmentList{
		Fragments: []Fragment{osf},
	}

	data, err := json.Marshal(list)
	if err != nil {
		fmt.Println("Error serializing:", err)
		return
	}

	fmt.Printf("Test: %s", data)
	t.Error()
}

const testmodule = `
module Test where

open import Agda.Builtin.String

output : String
output = "{\"Fragments\":[{\"Name\":\"OS\",\"Ctors\":[]}]}"
`

func TestLoad(t *testing.T) {

	// get current dir
	ex, err := os.Executable()
	if err != nil {
		t.Error(err)
		return
	}
	curDir := filepath.Dir(ex)

	// create source file (in current dir)
	sourceFile := filepath.Join(curDir, "Test.agda")
	_, err = createAndWrite(sourceFile, testmodule)
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	// Try to load fragments
	fs, err := LoadFragments(sourceFile)
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	// print
	fmt.Println("loaded: ", fs)
	t.Error()
}
