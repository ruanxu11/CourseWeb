package main

import (
	"crypto/md5"
	"fmt"
)

func hash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
