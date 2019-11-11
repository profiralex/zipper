package zipper

import (
	"os"
	"testing"
)

const (
	testEnvStringKey1   = "testEnvKey1"
	testEnvStringValue1 = "testEnvValue1"
)

func TestEnvProviderKeyDoesNotExist(t *testing.T) {
	provider := EnvProvider{}
	value, err := provider.Lookup(testEnvStringKey1)

	if value != "" {
		t.Errorf("Unexpected value %s", value)
	}

	if err == nil {
		t.Errorf("Expected error not received")
	}
}

func TestEnvProviderKeyExists(t *testing.T) {
	os.Setenv(testEnvStringKey1, testEnvStringValue1)
	defer os.Unsetenv(testEnvStringKey1)

	provider := EnvProvider{}
	value, err := provider.Lookup(testEnvStringKey1)

	if value != testEnvStringValue1 {
		t.Errorf("Expected value %s got %s", testEnvStringValue1, value)
	}

	if err != nil {
		t.Errorf("Unexpected error recived %w", err)
	}
}
