package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func NewUUIDV7() string {
	var uuid [16]byte
	if _, err := rand.Read(uuid[:]); err != nil {
		panic(err)
	}

	now := time.Now().UnixMilli()
	uuid[0] = byte(now >> 40)
	uuid[1] = byte(now >> 32)
	uuid[2] = byte(now >> 24)
	uuid[3] = byte(now >> 16)
	uuid[4] = byte(now >> 8)
	uuid[5] = byte(now)
	uuid[6] = (uuid[6] & 0x0f) | 0x70
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	encoded := hex.EncodeToString(uuid[:])
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		encoded[0:8],
		encoded[8:12],
		encoded[12:16],
		encoded[16:20],
		encoded[20:32],
	)
}
