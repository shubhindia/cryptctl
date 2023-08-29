package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go/aws"
)

func InitAwsKms(string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("error loading aws config")
	}
	client := kms.NewFromConfig(cfg)

	result, err := client.CreateKey(context.TODO(), &kms.CreateKeyInput{

		Description: aws.String("cryptctl-key"),
	})

	if err != nil {
		return fmt.Errorf("error creating aws-kms key: %s", err.Error())
	}

	// create alias for the same key so that we can fetch the key using alias

	_, err = client.CreateAlias(context.TODO(), &kms.CreateAliasInput{
		AliasName:   aws.String("alias/cryptctl-key"),
		TargetKeyId: result.KeyMetadata.KeyId,
	})
	if err != nil {
		return fmt.Errorf("error creating alias for kms key %s : %s", *result.KeyMetadata.KeyId, err.Error())
	}

	fmt.Printf("created aws-kms key: %s", *result.KeyMetadata.KeyId)
	return nil
}
