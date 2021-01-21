package mstore

import (
	"fmt"
	"os"

	"github.com/graymeta/stow"
	"github.com/graymeta/stow/azure"
	"github.com/graymeta/stow/s3"
)

// Vault -
type Vault struct {
	*stow.ConfigMap
	location      stow.Location
	kind          string
	storageGroup  string
	containerName string
	endpoint      string
	envVars       envVars
}

// ConnectToVault connects thservice to the Vault
func ConnectToVault() (s *Vault, err error) {
	s = &Vault{
		ConfigMap: &stow.ConfigMap{},
	}

	if err = s.GetEnvVars(); err != nil {
		return
	}

	if err = s.dial(); err != nil {
		return
	}

	return
}

// GetEnvVars will parse all the local environment variables that are required by mstore
func (s *Vault) GetEnvVars() (err error) {
	s.kind, err = MustEnv("STORAGE_KIND")
	if err != nil {
		return
	}

	s.containerName, err = MustEnv("CONTAINER_NAME")
	if err != nil {
		return
	}

	s.endpoint = os.Getenv("STORAGE_ENDPOINT") // only for use with mock storage (minio etc)

	switch s.kind {
	case "s3":
		if err = s.populateS3ConfigMap(); err != nil {
			return
		}
	case "azure":
		if err = s.populateAzureConfigMap(); err != nil {
			return
		}
	default:
		return fmt.Errorf("could not populate the config map due to unknown storage kind: %s", s.kind)
	}

	return nil
}

func (s *Vault) populateS3ConfigMap() (err error) {
	keyID, err := MustEnv("AWS_ACCESS_KEY_ID")
	if err != nil {
		return
	}

	secretKey, err := MustEnv("AWS_SECRET_ACCESS_KEY")
	if err != nil {
		return
	}

	region, err := MustEnv("AWS_REGION")
	if err != nil {
		return
	}

	s.ConfigMap = &stow.ConfigMap{
		s3.ConfigAccessKeyID: keyID,
		s3.ConfigSecretKey:   secretKey,
		s3.ConfigRegion:      region,
	}

	return nil
}

func (s *Vault) populateAzureConfigMap() (err error) {
	s.storageGroup, err = MustEnv("AZURE_STORAGE_GROUP")
	if err != nil {
		return
	}

	account, err := MustEnv("AZURE_STORAGE_ACCOUNT")
	if err != nil {
		return
	}

	accessKey, err := MustEnv("AZURE_STORAGE_ACCESS_KEY")
	if err != nil {
		return
	}

	s.ConfigMap = &stow.ConfigMap{
		azure.ConfigAccount: account,
		azure.ConfigKey:     accessKey,
	}

	return
}

// Dial to connect to Vault
func (s *Vault) dial() (err error) {
	s.location, err = stow.Dial(s.kind, s.ConfigMap)

	return
}

// Close the location connection
func (s *Vault) Close() error {
	return s.location.Close()
}

// GetContainer will fetch the requested container for the service
func (s *Vault) GetContainer() (cont *Container, err error) {
	walkContainersFunc := func(c stow.Container, err error) error {
		if err != nil {
			return err
		}

		if c.Name() == s.containerName {
			cont = &Container{c}
		}

		return nil
	}

	if err := stow.WalkContainers(s.location, s.containerName, 100, walkContainersFunc); err != nil {
		return nil, err
	}

	return
}
