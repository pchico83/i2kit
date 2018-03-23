package service

import (
	"testing"

	"github.com/pchico83/i2kit/cli/schemas/environment"
)

func TestGetSize(t *testing.T) {
	e := &environment.Environment{}
	s := &Service{}
	if s.GetInstanceType(e) != "t2.small" {
		t.Fatal("Expected 't2.nano' default value")
	}
	e.Provider = &environment.Provider{InstanceType: "t2.micro"}
	if s.GetInstanceType(e) != "t2.micro" {
		t.Fatal("Expected 't2.micro' value from environment.yml")
	}
	s.InstanceType = "t2.small"
	if s.GetInstanceType(e) != "t2.small" {
		t.Fatal("Expected 't2.small' value from service.yml")
	}
}
