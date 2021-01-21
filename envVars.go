package mstore

import (
	"fmt"
	"os"
)

type envVars struct {
	s3            *s3EnvVars
	az            *azureEnvVars
	storageKind   string
	containerName string
	endpoint      string
}

type s3EnvVars struct {
	keyID     string
	secretKey string
	region    string
}

type azureEnvVars struct {
	account      string
	accessKey    string
	storageGroup string
}

func parseEnvVars() (ev *envVars, err error) {
	ev = &envVars{}

	ev.storageKind, err = MustEnv("STORAGE_KIND")
	if err != nil {
		return
	}

	ev.containerName, err = MustEnv("CONTAINER_NAME")
	if err != nil {
		return
	}

	ev.endpoint = os.Getenv("S3_ENDPOINT") // only for use with mock storage (minio etc)

	switch ev.storageKind {
	case "s3":
		ev.s3, err = parseS3EnvVars()
		if err != nil {
			return
		}
	case "azure":
		ev.az, err = parseAzureEnvVars()
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("unknown kind - %s. Could not get env vars", ev.storageKind)
	}

	return
}

func parseS3EnvVars() (ev *s3EnvVars, err error) {
	ev = &s3EnvVars{}

	ev.keyID, err = MustEnv("AWS_ACCESS_KEY_ID")
	if err != nil {
		return
	}

	ev.secretKey, err = MustEnv("AWS_SECRET_ACCESS_KEY")
	if err != nil {
		return
	}

	ev.region, err = MustEnv("AWS_REGION")
	if err != nil {
		return
	}

	return
}

func parseAzureEnvVars() (ev *azureEnvVars, err error) {
	ev = &azureEnvVars{}

	ev.account, err = MustEnv("AZURE_STORAGE_ACCOUNT")
	if err != nil {
		return
	}

	ev.accessKey, err = MustEnv("AZURE_STORAGE_ACCESS_KEY")
	if err != nil {
		return
	}

	ev.storageGroup, err = MustEnv("AZURE_STORAGE_GROUP")
	if err != nil {
		return
	}

	return
}

func (ev *envVars) GetContainerURL() string {
	if ev.endpoint != "" {
		return fmt.Sprintf("%s/%s", ev.endpoint, ev.containerName)
	}

	switch ev.storageKind {
	case "s3":
		return fmt.Sprintf("https://s3.%s.amazonaws.com/%s", ev.s3.region, ev.containerName)
	default:
		return fmt.Sprintf("https://%s.blob.windows.net/%s", ev.az.storageGroup, ev.containerName)
	}
}

func (ev *envVars) ContainerName() string {
	return ev.containerName
}
