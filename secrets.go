package gcp_secret_manager

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"errors"
	"fmt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
)

var (
	Client    SecretClient
	ProjectID string
)

func init() {
	var ctx = context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	Client = &secretClientImpl{client: c, ctx: ctx}
}
func CreateEmptySecret(secretName string) (*secretmanagerpb.Secret, error) {
	if SecretExists(secretName) == true {
		log.Println("failed to create secret as secret already exists")
		return nil, errors.New("failed to create secret as secret already exists")
	}
	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", ProjectID),
		SecretId: secretName,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}
	secret, err := Client.CreateSecret(createSecretReq)
	if err != nil {
		log.Printf("failed to create secret: %v", err)
		return nil, err
	}
	return secret, nil
}
func CreateSecretWithData(secretName string, payload []byte) (*secretmanagerpb.SecretVersion, error) {
	if SecretExists(secretName) == true {
		log.Println("failed to create secret as secret already exists")
		return nil, errors.New("failed to create secret as secret already exists")
	}
	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", ProjectID),
		SecretId: secretName,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}
	secret, err := Client.CreateSecret(createSecretReq)
	if err != nil {
		log.Printf("failed to create secret: %v\n", err)
		return nil, err
	}
	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secret.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}
	version, err := Client.AddSecretVersion(addSecretVersionReq)
	if err != nil {
		log.Printf("failed to add secret version: %v\n", err)
		return nil, err
	}
	return version, err
}

func SecretExists(secretName string) bool {
	accessRequest := &secretmanagerpb.GetSecretRequest{
		Name: fmt.Sprintf("projects/%v/secrets/%v", ProjectID, secretName)}
	_, err := Client.GetSecret(accessRequest)
	if err != nil {
		return false
	}
	return true
}

func AddNewSecretVersion(secretName string, payload []byte) (*secretmanagerpb.SecretVersion, error) {
	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: fmt.Sprintf("projects/%v/secrets/%v", ProjectID, secretName),
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}
	version, err := Client.AddSecretVersion(addSecretVersionReq)
	if err != nil {
		log.Printf("failed to add secret version: %v", err)
		return nil, err
	}
	return version, nil
}
func GetSecret(secretName string, version string) (*secretmanagerpb.SecretPayload, error) {
	if version == "" {
		version = "latest"
	}
	getSecret := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%v/secrets/%v/versions/%v", ProjectID, secretName, version),
	}
	result, err := Client.AccessSecretVersion(getSecret)
	if err != nil {
		log.Printf("failed to get secret: %v", err)
		return nil, err
	}
	return result.Payload, nil
}

func DeleteSecretAndVersions(secretName string) error {
	deleteSecretReq := &secretmanagerpb.DeleteSecretRequest{
		Name: fmt.Sprintf("projects/%v/secrets/%v", ProjectID, secretName),
	}
	err := Client.DeleteSecret(deleteSecretReq)
	if err == nil {
		log.Printf("Secret Deleted Successfully")
	}
	return err
}

func DeleteSecretVersion(secretName string, version string) (*secretmanagerpb.SecretVersion, error) {
	destroySecretReq := &secretmanagerpb.DestroySecretVersionRequest{
		Name: fmt.Sprintf("projects/%v/secrets/%v/versions/%v", ProjectID, secretName, version),
	}
	result, err := Client.DestroySecretVersion(destroySecretReq)
	if err != nil {
		log.Printf("failed to get secret: %v", err)
		return nil, err
	}
	return result, nil
}

func GetSecretMetadata(secretName string, version string) (*secretmanagerpb.SecretVersion, error) {
	getSecretReq := &secretmanagerpb.GetSecretVersionRequest{
		Name: fmt.Sprintf("projects/%v/secrets/%v/versions/%v", ProjectID, secretName, version),
	}
	result, err := Client.GetSecretVersion(getSecretReq)
	if err != nil {
		log.Printf("failed to get secret: %v", err)
		return nil, err
	}
	return result, nil
}

func DisableSecret(secretName string, version string) (*secretmanagerpb.SecretVersion, error) {
	disableSecretReq := &secretmanagerpb.DisableSecretVersionRequest{
		Name: fmt.Sprintf("projects/%v/secrets/%v/versions/%v", ProjectID, secretName, version),
	}
	result, err := Client.DisableSecretVersion(disableSecretReq)
	if err != nil {
		log.Printf("failed to get secret: %v", err)
		return nil, err
	}
	return result, nil
}

func EnableSecret(secretName string, version string) (*secretmanagerpb.SecretVersion, error) {
	enableSecretReq := &secretmanagerpb.EnableSecretVersionRequest{
		Name: fmt.Sprintf("projects/%v/secrets/%v/versions/%v", ProjectID, secretName, version),
	}
	result, err := Client.EnableSecretVersion(enableSecretReq)
	if err != nil {
		log.Printf("failed to get secret: %v\n", err)
		return nil, err
	}
	return result, nil
}
