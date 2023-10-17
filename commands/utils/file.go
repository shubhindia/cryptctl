// this file contains the file related functions

package utils

import (
	"fmt"
	"os"
	"regexp"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"

	secretsv1alpha1 "github.com/shubhindia/encrypted-secrets/api/v1alpha1"
)

const (
	SecretApiVersion = "secrets.shubhindia.xyz/v1alpha1"
)

func ParseEncryptedSecret(filename string) (*secretsv1alpha1.EncryptedSecret, error) {
	encryptedFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s", err.Error())
	}

	codecs := serializer.NewCodecFactory(scheme.Scheme, serializer.EnableStrict)
	obj, _, err := codecs.UniversalDeserializer().Decode(encryptedFile, &schema.GroupVersionKind{
		Group:   secretsv1alpha1.GroupVersion.Group,
		Version: SecretApiVersion,
		Kind:    "EncryptedSecret",
	}, nil)
	if err != nil {
		if ok, _ := regexp.MatchString("no kind(.*)is registered for version", err.Error()); ok {
			panic("no kind(.*)is registered for version")
		}
		panic(err)
	}
	// convert the runtimeObj to encryptedSecret object
	encryptedSecret, ok := obj.(*secretsv1alpha1.EncryptedSecret)
	if !ok {
		// should never happen
		panic("failed to convert runtimeObject to encryptedSecret")
	}
	return encryptedSecret, nil
}
