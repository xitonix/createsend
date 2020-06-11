package test

import "errors"

func CheckError(actual, expected error) bool {
	if actual == nil {
		return expected == nil
	}
	if expected == nil {
		return false
	}
	return errors.Is(actual, expected)
}
