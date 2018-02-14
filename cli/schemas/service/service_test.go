package service

import (
	"testing"

	"github.com/pchico83/i2kit/cli/schemas/environment"
)

func TestGetSize(t *testing.T) {
	e := &environment.Environment{}
	s := &Service{}
	if s.GetSize(e) != "t2.nano" {
		t.Fatal("Expected 't2.nano' default value")
	}
	e.Provider = &environment.Provider{Size: "t2.micro"}
	if s.GetSize(e) != "t2.micro" {
		t.Fatal("Expected 't2.micro' value from environment.yml")
	}
	s.Size = "t2.small"
	if s.GetSize(e) != "t2.small" {
		t.Fatal("Expected 't2.small' value from service.yml")
	}
}
