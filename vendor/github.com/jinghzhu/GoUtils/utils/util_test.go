package utils

import (
	"fmt"
	"testing"
)

func TestUtilStruct2String(t *testing.T) {
	fmt.Println("Start testing util.UtilStruct2String()")
	arr1 := []string{"hello", "bonjour"}
	s := Struct2String(arr1)
	if s != "[\"hello\",\"bonjour\"]" {
		t.Error("\tShould return [\"hello\",\"bonjour\"]")
	}
	arr2 := []int{1, 2, 3}
	s = Struct2String(arr2)
	if s != "[1,2,3]" {
		t.Error("\tShould return [1,2,3]")
	}
	arr3 := []float32{1.01, 2.0, -3.45}
	s = Struct2String(arr3)
	if s != "[1.01,2,-3.45]" {
		t.Error("\tShould return [1.01,2,-3.45]")
	}

	m := make(map[string]interface{})
	m["k1"] = arr1
	m["k2"] = arr2
	s = Struct2String(m)
	r := "{\"k1\":[\"hello\",\"bonjour\"],\"k2\":[1,2,3]}"
	if s != r {
		t.Error("\tShould return " + r)
	}

}
