package main

import (
	"os"
)

// In
// ("lindsay", ["matt" "mark" "john" "lindsay"]) --> true
func In(needle string, haystack []string) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}

func IsReadOnly(filepath string) (bool, error) {
	fi, err := os.Stat(filepath)

	if err != nil {
		return false, err
	}

	return fi.Mode().Perm()&0200 == 0, nil
}

func IsDir(filepath string) (bool, error) {
	fi, err := os.Stat(filepath)

	if err != nil {
		return false, err
	}

	return fi.Mode().IsDir(), nil
}
