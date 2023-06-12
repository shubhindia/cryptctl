package edit

import (
	"k8s.io/apimachinery/pkg/runtime"

	secretsv1alpha1 "github.com/shubhindia/encrypted-secrets/api/v1alpha1"
	hacksecretsv1alpha1 "github.com/shubhindia/hcictl/commands/utils/apis/secrets/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Manifest []*Object

type Object struct {
	// The original text as parsed by NewYAMLOrJSONDecoder.
	Raw []byte
	// The original object as decoded by UniversalDeserializer.
	Object runtime.Object
	Meta   metav1.Object

	// Tracking for the various stages of encryption and decryption.
	OrigEnc  *secretsv1alpha1.EncryptedSecret
	OrigDec  *hacksecretsv1alpha1.DecryptedSecret
	AfterDec *hacksecretsv1alpha1.DecryptedSecret
	AfterEnc *secretsv1alpha1.EncryptedSecret
	Kind     string
	Data     map[string]string

	// The KMS KeyId used for this object, if known. If nil, it might be a new
	// object.
	KeyId string
	// The Plaintext Data key and Cipher Key generated using KMS Key ID
	PlainDataKey  *[32]byte
	CipherDataKey []byte

	// Byte coordinates for areas of the raw text we need to edit when re-serializing.
	KindLoc TextLocation
	DataLoc TextLocation
	KeyLocs []KeysLocation
}

type TextLocation struct {
	Start int
	End   int
}

type KeysLocation struct {
	TextLocation
	Key string
}
