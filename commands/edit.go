package commands

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/shubhindia/hcictl/common"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"

	editutils "github.com/shubhindia/hcictl/commands/utils/edit"
)

func init() {

	cliCmd := cli.Command{
		Name:  "edit",
		Usage: "edit encryptedSecrets manifest",
		Before: func(ctx *cli.Context) error {
			if ctx.Args().First() == "" {
				return fmt.Errorf("hcictl edit expectes a file to edit")
			}

			if ctx.Args().Len() > 1 {
				return fmt.Errorf("too many arguments")
			}

			return nil
		},
		Action: func(ctx *cli.Context) error {
			fileName := ctx.Args().First()
			encryptedFile, err := os.ReadFile(fileName)
			if err != nil {
				return fmt.Errorf("error reading file %s", err.Error())
			}

			var encryptedSecret editutils.EncryptedSecret

			// unmarshal into EncryptedSecret
			err = yaml.Unmarshal(encryptedFile, &encryptedSecret)
			if err != nil {
				return fmt.Errorf("error unmarshaling file %s", err.Error())
			}

			// prepare decryptedSecret to be edited
			decryptedSecret := editutils.DecryptedSecret{
				ApiVersion: encryptedSecret.ApiVersion,
				Kind:       "DecryptedSecret",
				Metadata:   encryptedSecret.Metadata,
			}

			keyPhrase := os.Getenv("KEYPHRASE")
			if keyPhrase == "" {
				return fmt.Errorf("keyphrase not found")
			}

			decryptedData := make(map[string]string)

			// decrypt the data in encryptedSecrets
			for key, value := range encryptedSecret.Data {
				decryptedString := decodeAndDecrypt(value, keyPhrase)

				decryptedData[key] = string(decryptedString)
			}

			decryptedSecret.Data = decryptedData

			fmt.Printf("%+v", decryptedSecret)

			return nil

		},
	}

	common.RegisterCommand(cliCmd)
}

func mdHashing(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:]) // by referring to it as a string
}

func decodeAndDecrypt(encoded string, keyPhrase string) []byte {
	ciphered, _ := base64.StdEncoding.DecodeString(encoded)
	hashedPhrase := mdHashing(keyPhrase)
	aesBlock, err := aes.NewCipher([]byte(hashedPhrase))
	if err != nil {
		log.Fatalln(err)
	}
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Fatalln(err)
	}
	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]

	originalText, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return originalText

}
