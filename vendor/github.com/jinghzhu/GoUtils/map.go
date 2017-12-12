package utils

import (
    "fmt"
    "sort"
)

// Check whether two maps are equcal, whose key-value types are String-Int
// The performance is better than reflect.DeepEqual
func EqualStringInt(x, y map[string]int) bool {
    if x == nil && y == nil {
        return true
    } else if x == nil || y == nil {
        return false
    }
    if len(x) != len(y) {
        return false
    }
    for k, xv := range x {
        if yv, ok := y[k]; !ok || yv != xv {
            return false
        }
    }
    return true
}


// Print map in the order of keys
func PrintByValueOrder(m map[string]int) {
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, v := range keys {
        fmt.Printf("key = %s, value = %d\n", v, m[v])
    }
}


// The following 3 functions are developed for supporting string slice type as key in map
func KeyStrings(list []string) string {
    return fmt.Sprintf("%q", list)
}

func SetValueIntForKeyStrings(list []string, v int, m map[string]int) {
    m[KeyStrings(list)] = v
}

func GetValueIntForKeyStrings(list []string, m map[string]int) int {
    return m[KeyStrings(list)]
}
