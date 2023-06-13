package edit

import (
	staticProvider "github.com/shubhindia/crypt-core/providers/static"
)

func DecodeAndDecrypt(encoded string, keyPhrase string) []byte {
	return staticProvider.DecodeAndDecrypt(encoded, keyPhrase)

}

func EncryptAndEncode(value string, keyPhrase string) (string, error) {

	return staticProvider.EncryptAndEncode(value, keyPhrase)

}
