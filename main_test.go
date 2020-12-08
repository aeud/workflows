package main

import (
	"log"
	"testing"
)

func TestDAG(t *testing.T) {
	if err := exec("./dag.yaml"); err != nil {
		log.Println(err)
		t.Fail()
	}
}
