package edit

import (
	"regexp"

	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"

	secretsv1alpha1 "github.com/shubhindia/encrypted-secrets/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	nonceLength = 24
)

var dataRegexp *regexp.Regexp
var keyRegexp *regexp.Regexp
var nonStringRegexp *regexp.Regexp

type Payload struct {
	Key     []byte
	Nonce   *[nonceLength]byte
	Message []byte
}

func init() {
	dataRegexp = regexp.MustCompile(`(?ms)kind: (EncryptedSecret|DecryptedSecret).*?(^data:.*?)\z`)
	keyRegexp = regexp.MustCompile("" +
		// Turn on multiline mode for the whole pattern, ^ and $ will match on lines rather than start and end of whole string.
		`(?m)` +
		// Look for the key, some whitespace, then some non-space-or-:, then :
		`^[ \t]+([^:\n\r]+):` +
		// Whitespace between the key's : and the value
		`[ \t]+` +
		// Start an alternation for block scalars and normal values.
		`(?:` +
		// Match block scalars first because they would otherwise match the normal value pattern.
		// Looks for the | or >, optional flags, then lines with 4 spaces of indentation. A better version of this
		// would look more like ([|>]\n([ \t]+).+?\n(?:\3.+?\n)*) and would use a backreference instead of hardwiring
		// things but Go, or rather RE2, refuses to support backrefs because they can be slow. Blaaaaaaah.
		`([|>].*?(?:\n    .+?$)+)` +
		// Alternation between block scalar and normal values.
		`|` +
		// Look for a normal value, something on a single line with optional trailing whitespace.
		`(.+?)[ \t]*$` +
		// Close the block vs. normal alternation.
		`)`,
	)
	nonStringRegexp = regexp.MustCompile(`^(\d+(\.\d+)?|true|false|null|\[.*\]|)$`)
}

func NewObject(raw []byte) (*Object, error) {

	// Here, we need to be able to edit the objects which are not registered in ridectl
	// So we are deserializing the all the object, if object is not registered, UniversalDeserializer()
	// will return error 'no kind "xyz" is registered for version "abc"', with return the object with Raw
	// field set in it.

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
