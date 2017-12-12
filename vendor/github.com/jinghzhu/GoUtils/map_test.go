package utils

import (
        "testing"
)


func TestEqualStringInt(t *testing.T) {
    var m1 map[string]int = nil
    var m2 map[string]int
    m3 := make(map[string]int)
    // var m4 *map[string]int
    // m4 = new(map[string]int)
    m5 := map[string]int{}
    m6 := make(map[string]int, 1)
    if !EqualStringInt(m1, m2) {
        t.Error("error1")
    }
    if EqualStringInt(m1, m3) {
        t.Error("error2")
    }
    // if !EqualStringInt(m1, m4) {
    //     t.Error("error3")
    // }
    if EqualStringInt(m1, m5) {
        t.Error("error4")
    }
    if EqualStringInt(m1, m6) {
        t.Error("error5")
    }
    if EqualStringInt(map[string]int{"a": 0}, map[string]int{"b": 1}) {
        t.Error("error6")
    }
    if !EqualStringInt(map[string]int{"a": 0}, map[string]int{"a": 0}) {
        t.Error("error7")
    }
}