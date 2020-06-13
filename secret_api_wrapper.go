package gcp_secret_manager

import (
	"context"
	"fmt"

	sm "cloud.google.com/go/secretmanager/apiv1"
	smpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// ClientFactory is used to create SecretClient, which is the GRPC Secret Client
// in normal use, but can be mocked for tests.
type ClientFactory interface {
	NewSecretClient(ctx context.Context) (SecretClient, error)
}

// SecretClient is a wrapper around the secretmanager APIs that are used by smcache.
// It is entirely for the purpose of being able to mock these for testing.
type SecretClient interface {
	AccessSecretVersion(req *smpb.AccessSecretVersionRequest) (*smpb.AccessSecretVersionResponse, error)
	ListSecretVersions(req *smpb.ListSecretVersionsRequest) SecretListIterator
	DestroySecretVersion(req *smpb.DestroySecretVersionRequest) (*smpb.SecretVersion, error)
	CreateSecret(req *smpb.CreateSecretRequest) (*smpb.Secret, error)
	AddSecretVersion(req *smpb.AddSecretVersionRequest) (*smpb.SecretVersion, error)
	DeleteSecret(req *smpb.DeleteSecretRequest) error
	GetSecret(req *smpb.GetSecretRequest) (*smpb.Secret, error)
	ListSecrets(req *smpb.ListSecretsRequest) *sm.SecretIterator
	GetSecretVersion(req *smpb.GetSecretVersionRequest) (*smpb.SecretVersion, error)
	DisableSecretVersion(req *smpb.DisableSecretVersionRequest) (*smpb.SecretVersion, error)
	EnableSecretVersion(req *smpb.EnableSecretVersionRequest) (*smpb.SecretVersion, error)
	Close() error
}

// SecretListIterator is an interface for the GRPC secret manager response from ListSecretVersions.
type SecretListIterator interface {
	Next() (*smpb.SecretVersion, error)
}

type secretClientImpl struct {
	client *sm.Client
	ctx    context.Context
}

// SecretClientFactoryImpl implements ClientFactory for the real GRPC client.
type SecretClientFactoryImpl struct{}

// NewSecretClient creates a GRPC NewClient for secretmanager.
func (*SecretClientFactoryImpl) NewSecretClient(ctx context.Context) (SecretClient, error) {
	c, err := sm.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup client: %w", err)
	}

	return &secretClientImpl{client: c, ctx: ctx}, nil
}

func (sc *secretClientImpl) AccessSecretVersion(req *smpb.AccessSecretVersionRequest) (*smpb.AccessSecretVersionResponse, error) {
	return sc.client.AccessSecretVersion(sc.ctx, req)
}
func (sc *secretClientImpl) ListSecretVersions(req *smpb.ListSecretVersionsRequest) SecretListIterator {
	return sc.client.ListSecretVersions(sc.ctx, req)
}
func (sc *secretClientImpl) DestroySecretVersion(req *smpb.DestroySecretVersionRequest) (*smpb.SecretVersion, error) {
	return sc.client.DestroySecretVersion(sc.ctx, req)
}
func (sc *secretClientImpl) CreateSecret(req *smpb.CreateSecretRequest) (*smpb.Secret, error) {
	return sc.client.CreateSecret(sc.ctx, req)
}
func (sc *secretClientImpl) AddSecretVersion(req *smpb.AddSecretVersionRequest) (*smpb.SecretVersion, error) {
	return sc.client.AddSecretVersion(sc.ctx, req)
}
func (sc *secretClientImpl) DeleteSecret(req *smpb.DeleteSecretRequest) error {
	return sc.client.DeleteSecret(sc.ctx, req)
}

func (sc *secretClientImpl) GetSecret(req *smpb.GetSecretRequest) (*smpb.Secret, error) {
	return sc.client.GetSecret(sc.ctx, req)
}

func (sc *secretClientImpl) GetSecretVersion(req *smpb.GetSecretVersionRequest) (*smpb.SecretVersion, error) {
	return sc.client.GetSecretVersion(sc.ctx, req)
}

func (sc *secretClientImpl) ListSecrets(req *smpb.ListSecretsRequest) *sm.SecretIterator {
	return sc.client.ListSecrets(sc.ctx, req)
}

func (sc *secretClientImpl) DisableSecretVersion(req *smpb.DisableSecretVersionRequest) (*smpb.SecretVersion, error) {
	return sc.client.DisableSecretVersion(sc.ctx, req)
}

func (sc *secretClientImpl) EnableSecretVersion(req *smpb.EnableSecretVersionRequest) (*smpb.SecretVersion, error) {
	return sc.client.EnableSecretVersion(sc.ctx, req)
}

func (sc *secretClientImpl) Close() error {
	return sc.client.Close()
}