package array

import (
	"fmt"
	"testing"
)

var (
	arr1 = []string{}
	arr2 = []string{"hello", "你好"}
)

func TestArrayIndex(t *testing.T) {
	fmt.Println("Start testing array.Index()")
	if Index(arr1, "test") != -1 {
		t.Error("\tShould return -1 for empty array")
	}
	if Index(arr2, "hello") != 0 {
		t.Error("\tShould return 0 for string hello")
	}
	if Index(arr2, "你好") != 1 {
		t.Error("\tShould return 1 for string 你好")
	}
	if Index(arr2, "world") != -1 {
		t.Error("\tShould return -1 for string world")
	}
}

func TestArrayInclude(t *testing.T) {
	fmt.Println("Start testing array.Include()")
	if Include(arr1, "test") {
		t.Error("\nShould return false for empty array")
	}
	if !Include(arr2, "hello") {
		t.Error("\nShould return true for string hello")
	}
	if !Include(arr2, "你好") {
		t.Error("\nShould return true for string 你好")
	}
	if Include(arr2, "world") {
		t.Error("\nShould return false for string world")
	}
}
