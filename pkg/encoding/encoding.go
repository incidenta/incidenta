package encoding

import (
	"crypto/sha1"
	"encoding/hex"
)

func EncodeSha1(str string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
