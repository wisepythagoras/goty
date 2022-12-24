package main

import (
	"errors"
	"io/ioutil"
	"strings"
)

// WriteToFile writes the contents to the file at the given path.
func WriteToFile(path string, contents []byte) error {
	return ioutil.WriteFile(path, contents, 0600)
}

// ReadFile reads the contents of the file at the given path.
func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// ParseAddress parses a string address.
func ParseAddr(addr string) (string, string, error) {
	parts := strings.Split(addr, "://")

	if len(parts) != 2 {
		return "", "", errors.New("Invalid address")
	}

	if parts[0] != "tcp" && parts[0] != "tcp4" && parts[0] != "tcp6" && parts[0] != "unix" {
		return "", "", errors.New("Invalid protocol")
	}

	return parts[0], parts[1], nil
}
