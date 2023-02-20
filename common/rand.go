package common

import (
	"math/rand"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

// get random bytes of given length made up of
// english lowercase letters
func NewRandomString(length int) string {
	bytes := make([]byte, length)
        for i := 0; i < length; i++ {
                bytes[i] = byte(97 + rand.Intn(25))
        }
        return string(bytes)
}
