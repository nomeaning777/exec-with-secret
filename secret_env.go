package secret

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	"context"
	"fmt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
	"net/url"
	"os"
	"strings"
)

type secretPath struct {
	Project string
	Secret  string
	Version string
}

func getSecret(client *secretmanager.Client, path secretPath) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", path.Project, path.Secret, path.Version),
	}
	resp, err := client.AccessSecretVersion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}
	return string(resp.Payload.Data), nil
}

func parseSecretUrl(key string) (secretPath, error) {
	secretUrl, err := url.Parse(key)
	if err != nil {
		return secretPath{}, fmt.Errorf("failed to parse secretUrl: %w", err)
	}
	projectName := secretUrl.Host
	secretName := secretUrl.Path[1:]
	version := "latest"
	if strings.Contains(secretName, "/") {
		splitSecretName := strings.SplitN(secretName, "/", 2)
		secretName = splitSecretName[0]
		version = splitSecretName[1]
	}
	return secretPath{
		Project: projectName,
		Secret:  secretName,
		Version: version,
	}, nil
}

func isSecretUrl(value string) bool {
	return strings.HasPrefix(value, "secretmanager://")
}

func InjectSecretToEnvironment() error {
	client, err := secretmanager.NewClient(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create secret manager client: %w", err)
	}

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}

		if !isSecretUrl(pair[1]) {
			continue
		}

		path, err := parseSecretUrl(pair[1])
		if err != nil {
			return fmt.Errorf("parse failed to environment variable %s: %w", pair[0], err)
		}

		secret, err := getSecret(client, path)
		if err != nil {
			return fmt.Errorf("get secret failed to environemnt variable %s: %w", pair[0], err)
		}

		if err := os.Setenv(pair[0], secret); err != nil {
			return fmt.Errorf("failed to setenv to %s: %+v", pair[0], err)
		}
	}
	return nil
}
