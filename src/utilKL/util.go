package utilKL

import (
	"crypto/md5"
	"fmt"
)

func hashString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
