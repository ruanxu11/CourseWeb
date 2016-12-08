package koala

import (
	"crypto/md5"
	"fmt"
)

func HashString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
