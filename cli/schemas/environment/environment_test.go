package environment

import (
	"testing"
)

func TestDockerConfig(t *testing.T) {
	e := &Environment{}
	b64DockerConfig := e.B64DockerConfig()
	if b64DockerConfig != "" {
		t.Fatal("Expected empty encoded value")
	}

	e.Docker = &Docker{
		Username: "username",
		Password: "password",
	}
	b64DockerConfig = e.B64DockerConfig()
	expectedValue := "CnsKCSJhdXRocyI6IHsKCQkiaHR0cHM6Ly9pbmRleC5kb2NrZXIuaW8vdjEvIjogewoJCQkiYXV0aCI6ICJkWE5sY201aGJXVTZjR0Z6YzNkdmNtUT0iCgkJfQoJfQp9Cg=="
	if b64DockerConfig != expectedValue {
		t.Fatalf("Wrong encoded value: %s", b64DockerConfig)
	}
}
