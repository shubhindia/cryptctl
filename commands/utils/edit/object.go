package edit

import (
	"regexp"

	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"

	secretsv1alpha1 "github.com/shubhindia/encrypted-secrets/api/v1alpha1"
	hacksecretsv1alpha1 "github.com/shubhindia/hcictl/commands/utils/apis/secrets/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	nonceLength = 24
)

var dataRegexp *regexp.Regexp

type Payload struct {
	Key     []byte
	Nonce   *[nonceLength]byte
	Message []byte
}

func init() {
	dataRegexp = regexp.MustCompile(`(?ms)kind: (EncryptedSecret|DecryptedSecret).*?(^data:.*?)\z`)

}

func NewObject(raw []byte) (*Object, error) {

	o := &Object{Raw: raw}

	// Create new codec with strict mode on; this will strictly check objects spec
	codecs := serializer.NewCodecFactory(scheme.Scheme, serializer.EnableStrict)
	obj, _, err := codecs.UniversalDeserializer().Decode(raw, nil, nil)

	if err != nil {
		if ok, _ := regexp.MatchString("no kind(.*)is registered for version", err.Error()); ok {
			return o, nil
		}
		return nil, err
	}

	o.Object = obj
	o.Meta = obj.(metav1.Object)

	// Check if this an EncryptedSecret.
	enc, ok := obj.(*secretsv1alpha1.EncryptedSecret)
	if ok {

		o.OrigEnc = enc
		o.Kind = "EncryptedSecret"
		o.Data = enc.Data
	}

	// // or a DecryptedSecret.
	// dec, ok := obj.(*hacksecretsv1beta2.DecryptedSecret)
	// if ok {
	// 	o.AfterDec = dec
	// 	o.Kind = "DecryptedSecret"
	// 	o.Data = dec.Data
	// }

	if o.Kind != "" {
		// Run the regex parse. If you are reading this code, I am sorry and yes I
		// feel bad about it. This is used when re-encoding to allow output that
		// preserves comments, whitespace, key ordering, etc.
		match := dataRegexp.FindSubmatchIndex(raw)
		if match == nil {
			// This shouldn't happen.
			panic("EncryptedSecret or DecryptedSecret didn't match dataRegexp")
		}
		// match[0] and [1] are for the whole regexp, we don't need that.
		o.KindLoc.Start = match[2]
		o.KindLoc.End = match[3]
		o.DataLoc.Start = match[4]
		o.DataLoc.End = match[5]

	}
	return o, nil
}

func (o *Object) Decrypt() error {

	if o.Kind == "" {
		return nil
	}

	dec := &hacksecretsv1alpha1.DecryptedSecret{ObjectMeta: o.OrigEnc.ObjectMeta, Data: map[string]string{}}
	for key, value := range o.OrigEnc.Data {
		dec.Data[key] = value
	}

	o.OrigDec = dec
	o.Kind = "DecryptedSecret"
	o.Data = dec.Data
	return nil

}
