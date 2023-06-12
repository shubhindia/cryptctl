package edit

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/pkg/errors"
)

var emptyRegexp *regexp.Regexp

func init() {
	emptyRegexp = regexp.MustCompile(`(?m)\A(^(\s*#.*|\s*)$\s*)*\z`)

}

func NewManifest(in io.Reader) (Manifest, error) {

	// Read in the whole file.
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(in)
	if err != nil {
		return nil, errors.Wrap(err, "error reading manifest")
	}

	objects := []*Object{}
	if emptyRegexp.MatchString(buf.String()) {
		return nil, fmt.Errorf("empty object")
	}

	obj, err := NewObject([]byte(buf.Bytes()))
	if err != nil {
		return nil, errors.Wrap(err, "error decoding object")
	}

	objects = append(objects, obj)
	return objects, nil
}

func (m Manifest) Decrypt() error {
	for _, obj := range m {
		err := obj.Decrypt()
		if err != nil {
			return errors.Wrap(err, "error decrypting object")
		}
	}
	return nil
}

func (m Manifest) Serialize(out io.Writer) error {
	first := true
	for _, obj := range m {
		if !first {
			_, _ = out.Write([]byte("---\n"))
		}
		first = false
		err := obj.Serialize(out)
		if err != nil {
			return errors.Wrapf(err, "error serializing %s/%s", obj.Meta.GetNamespace(), obj.Meta.GetName())
		}
	}
	return nil
}
