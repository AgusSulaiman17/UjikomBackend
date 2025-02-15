package utils

import "strconv"

// ParseUint converts a string to uint
func ParseUint(str string) (uint, error) {
    val, err := strconv.ParseUint(str, 10, 32)
    if err != nil {
        return 0, err
    }
    return uint(val), nil
}
