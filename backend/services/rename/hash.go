package rename

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func md5Hash(data []byte) string {
	h := md5.Sum(data)
	return hex.EncodeToString(h[:])
}

func sha1Hash(data []byte) string {
	h := sha1.Sum(data)
	return hex.EncodeToString(h[:])
}

func sha256Hash(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
