package utils

import (
        "testing"
)


// To-do: add test cases for non-English
func TestCompare(t *testing.T) {
        s1 := ""
        s2 := "Shanghai"
        s3 := "Paris"
        if Compare(s1, s2) != -1 {
            t.Error("error1")
        }
        if Compare(s2, "Shanghai") != 0 {
            t.Error("error2")
        }
        if Compare(s3, "paris") != -1 {
            t.Error("error3")
        }
}


func TestIndexRune(t *testing.T) {
        s1 := ""
        s2 := "abcd"
        var r rune = 99 // c
        if IndexRune(s1, r) != -1 {
            t.Error("error1")
        }
        if IndexRune(s2, r) != 2 {
            t.Error("error2")
        }
}


func TestIsAnagram(t *testing.T) {
        s1 := ""
        s2 := "test"
        if IsAnagram(s1, s2) {
            t.Error("error1")
        }
        s1 = "stte"
        if !IsAnagram(s1, s2) {
            t.Error("error2")
        }
        s1 = "我爱巴黎"
        s2 = "巴黎爱我"
        if !IsAnagram(s1, s2) {
            t.Error("error3")
        }
        s1 = "我爱上海"
        if IsAnagram(s1, s2) {
            t.Error("error4")
        }
}