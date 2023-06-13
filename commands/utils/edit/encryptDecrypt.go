package edit

import (
	"crypto/md5"
	"encoding/hex"

	staticProvider "github.com/shubhindia/crypt-core/providers/static"
)

func mdHashing(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:]) // by referring to it as a string
}

func DecodeAndDecrypt(encoded string, keyPhrase string) []byte {
	return staticProvider.DecodeAndDecrypt(encoded, keyPhrase)

}

func EncryptAndEncode(value string, keyPhrase string) (string, error) {

	return staticProvider.EncryptAndEncode(value, keyPhrase)

}
