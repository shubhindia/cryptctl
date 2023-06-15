package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/shubhindia/cryptctl/common"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"

	editutils "github.com/shubhindia/cryptctl/commands/utils/edit"

	providers "github.com/shubhindia/crypt-core/providers"
)

var whitespaceRegexp *regexp.Regexp

func init() {

	editCmd := cli.Command{
		Name:  "edit",
		Usage: "edit encryptedSecrets manifest",
		Before: func(ctx *cli.Context) error {
			if ctx.Args().First() == "" {
				return fmt.Errorf("cryptctl edit expectes a file to edit")
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

			// get the provider
			// not handling the error here because this should never fail

			// ToDo := use k8s labels here so that we can actually use label methods instead of converting interface to map again
			provider := encryptedSecret.Metadata["labels"].(map[interface{}]interface{})["app.kubernetes.io/provider"]

			// prepare decryptedSecret to be edited
			decryptedSecret := editutils.DecryptedSecret{
				ApiVersion: encryptedSecret.ApiVersion,
				Kind:       "DecryptedSecret",
				Metadata:   encryptedSecret.Metadata,
			}

			decryptedData := make(map[string]string)

			// decrypt the data in encryptedSecrets
			for key, value := range encryptedSecret.Data {
				decryptedString, err := providers.DecodeAndDecrypt(value, provider.(string))
				if err != nil {
					return fmt.Errorf("failed to decrypt value for %s %s", key, err.Error())
				}

				decryptedData[key] = decryptedString
			}

			decryptedSecret.Data = decryptedData

			// marshal into yaml
			decrypted, err := yaml.Marshal(&decryptedSecret)
			if err != nil {
				return fmt.Errorf("error marshaling decryptedSecret %s", err.Error())
			}

			editedManitest, err := editObjects(decrypted)
			if err != nil {
				return fmt.Errorf("error editing objects %s", err.Error())
			}

			// unmarshal the edited yaml into decryptedSecrets again to encrypt the new secrets
			var newDecryptedSecret editutils.DecryptedSecret

			err = yaml.Unmarshal(editedManitest, &newDecryptedSecret)
			if err != nil {
				return fmt.Errorf("error unmarshaling file %s", err.Error())
			}

			// prepare new encryptedSecret to be written
			newEncryptedSecret := encryptedSecret

			newEncryptedData := make(map[string]string)

			for key, value := range newDecryptedSecret.Data {
				encryptedString, err := providers.EncryptAndEncode(value, provider.(string))
				if err != nil {
					return fmt.Errorf("error encrypting new secrets %s", err)
				}
				newEncryptedData[key] = encryptedString
			}

			// write newly encrypted data
			newEncryptedSecret.Data = newEncryptedData

			// write the contents to yaml
			newEncrypted, err := yaml.Marshal(&newEncryptedSecret)
			if err != nil {
				return fmt.Errorf("error marshaling encryptedSecret %s", err.Error())
			}

			err = os.WriteFile(fileName, newEncrypted, 0600)
			if err != nil {
				return fmt.Errorf("error writing EncryptedSecret %s", err)
			}

			return nil

		},
	}

	common.RegisterCommand(editCmd)
}

func editObjects(data []byte) ([]byte, error) {
	manifestBuf := bytes.Buffer{}
	_, _ = manifestBuf.Write(data)

	for {

		// make the YAML to show in the editor.
		editorBuf := bytes.Buffer{}
		_, _ = manifestBuf.WriteTo(&editorBuf)
		editorReader := bytes.NewReader(editorBuf.Bytes())

		// open a temporary file.
		tmpfile, err := os.CreateTemp("", ".*.yaml")
		if err != nil {
			return nil, fmt.Errorf("error making tempfile %s", err.Error())
		}
		defer tmpfile.Close()
		defer os.Remove(tmpfile.Name())
		_, _ = editorReader.WriteTo(tmpfile)
		_ = tmpfile.Sync()

		// show the editor.
		err = runEditor(tmpfile.Name())
		if err != nil {
			return nil, fmt.Errorf("error running editor %s", err.Error())
		}

		// re-read the edited file.
		afterTmpfile, err := os.Open(tmpfile.Name())
		if err != nil {
			return nil, fmt.Errorf("error re-opening tempfile %s %s", tmpfile.Name(), err)
		}
		defer afterTmpfile.Close()
		afterBuf := bytes.Buffer{}
		_, err = afterBuf.ReadFrom(afterTmpfile)
		if err != nil {
			return nil, fmt.Errorf("error reading tempfile %s %s", tmpfile.Name(), err)
		}

		// check if the file was edited at all.
		if bytes.Equal(editorBuf.Bytes(), afterBuf.Bytes()) {
			fmt.Println("Edit cancelled. No changes made")
			os.Exit(0)
		} else {
			// return the edited bytes
			return afterBuf.Bytes(), nil
		}

	}

}

func runEditor(filename string) error {
	whitespaceRegexp = regexp.MustCompile(`\s+`)
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("nso $EDITOR set")
	}

	// deal with an editor that has options.
	editorParts := whitespaceRegexp.Split(editor, -1)
	executable := editorParts[0]
	executable, _ = exec.LookPath(executable)

	editorParts = append(editorParts, filename)
	cmd := exec.Command(executable, editorParts[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running editor %s", err)
	}
	return nil
}
