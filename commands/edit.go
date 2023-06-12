package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/shubhindia/hcictl/commands/utils/edit"
	"github.com/shubhindia/hcictl/common"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/kubernetes/scheme"

	secretsv1alpha1 "github.com/shubhindia/encrypted-secrets/api/v1alpha1"
)

var whitespaceRegexp *regexp.Regexp

func init() {
	whitespaceRegexp = regexp.MustCompile(`\s+`)

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
			// register the apis
			err := secretsv1alpha1.AddToScheme(scheme.Scheme)
			if err != nil {
				return fmt.Errorf("error registering apis %s", err.Error())
			}
			return nil
		},
		Action: func(ctx *cli.Context) error {

			fileName := ctx.Args().First()

			// read the file
			var inStream io.Reader
			inFile, err := os.Open(fileName)
			if err != nil {
				if os.IsNotExist(err) {
					return errors.Wrapf(err, "error reading input file %s", fileName)

				} else {
					return errors.Wrapf(err, "error reading input file %s", fileName)
				}
			} else {
				defer inFile.Close()
				inStream = inFile
			}

			// Parse the input file to objects.
			inManifest, err := edit.NewManifest(inStream)
			if err != nil {
				return errors.Wrap(err, "error decoding input YAML")
			}

			// decrypt the data
			err = inManifest.Decrypt()
			if err != nil {
				return errors.Wrap(err, "error decrypting input manifest")
			}

			// Edit!
			_, err = editObjects(inManifest, "")
			if err != nil {
				return errors.Wrap(err, "error editing objects")
			}

			// ToDo:
			// 1. Convert back to encryptedSecrets after editing
			// 2. Actually decrypt and encrypt the yaml

			return nil

		},
	}

	common.RegisterCommand(cliCmd)
}

func editObjects(manifest edit.Manifest, comment string) (edit.Manifest, error) {
	manifestBuf := bytes.Buffer{}
	err := manifest.Serialize(&manifestBuf)
	if err != nil {
		return nil, errors.Wrap(err, "error encoding objects to YAML")
	}
	for {
		// Format the comment.
		commentBuf := bytes.Buffer{}
		if comment != "" {
			for _, line := range strings.Split(comment, "\n") {
				commentBuf.WriteString("# ")
				commentBuf.WriteString(line)
				commentBuf.WriteString("\n")
			}
			commentBuf.WriteString("#\n")
		}
		commentReader := bytes.NewReader(commentBuf.Bytes())

		// Make the YAML to show in the editor.
		editorBuf := bytes.Buffer{}
		_, _ = commentReader.WriteTo(&editorBuf)
		_, _ = manifestBuf.WriteTo(&editorBuf)
		editorReader := bytes.NewReader(editorBuf.Bytes())

		// Open a temporary file.
		tmpfile, err := os.CreateTemp("", ".*.yml")
		if err != nil {
			return nil, errors.Wrap(err, "error making tempfile")
		}
		defer tmpfile.Close()
		defer os.Remove(tmpfile.Name())
		_, _ = editorReader.WriteTo(tmpfile)
		_ = tmpfile.Sync()

		// Show the editor.
		err = runEditor(tmpfile.Name())
		if err != nil {
			return nil, errors.Wrap(err, "error running editor")
		}

		// Re-read the edited file.
		afterTmpfile, err := os.Open(tmpfile.Name())
		if err != nil {
			return nil, errors.Wrapf(err, "error re-opening tempfile %s", tmpfile.Name())
		}
		defer afterTmpfile.Close()
		afterBuf := bytes.Buffer{}
		_, err = afterBuf.ReadFrom(afterTmpfile)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading tempfile %s", tmpfile.Name())
		}

		// If we're reencrypting ignore this equality check.
		// Check if the file was edited at all.
		if bytes.Equal(editorBuf.Bytes(), afterBuf.Bytes()) {
			fmt.Println("Edit cancelled. No changes made")
			os.Exit(0)
		}

		// Try strip off the comment.
		afterReader := bytes.NewReader(afterBuf.Bytes())
		seekPos := int64(0)
		if bytes.Equal(commentBuf.Bytes(), afterBuf.Bytes()[:commentBuf.Len()]) {
			seekPos = int64(commentBuf.Len())
		}
		_, _ = afterReader.Seek(seekPos, 0)

		outManifest, err := edit.NewManifest(afterReader)
		if err == nil {
			// Decode success, we're done!
			return outManifest, nil
		}

		// Some kind decoding error, probably bad syntax, show the editor again.
		comment = fmt.Sprintf("Error parsing file:\n%s", err)
		manifestBuf.Reset()
		_, _ = afterReader.Seek(seekPos, 0)
		_, _ = afterReader.WriteTo(&manifestBuf)
	}
}

func runEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return errors.New("No $EDITOR set")
	}

	// Deal with an editor that has options.
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
		return errors.Wrap(err, "error running editor")
	}
	return nil
}
