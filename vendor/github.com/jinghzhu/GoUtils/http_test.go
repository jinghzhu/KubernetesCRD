package utils

import (
        "testing"
)

func TestIsURL(t *testing.T) {
        url1 := "http://golang.org"
        if !IsURL(url1) {
            t.Error("error1")
        }

        url2 := ""
        if IsURL(url2) {
            t.Error("error2")
        }

        url3 := "https://www.google.com"
        if !IsURL(url3) {
            t.Error("error3")
        }

}