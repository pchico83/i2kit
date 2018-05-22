package environment

import (
	"bytes"
	"fmt"
	"strings"
)

//MarshalYAML serializes e into a YAML document. The return value is a string; It will fail if s has an empty name or value.
func (s *Secret) MarshalYAML() (interface{}, error) {
	if s.Name == "" {
		return "", fmt.Errorf("missing secret name")
	}
	if s.Value == "" {
		return "", fmt.Errorf("missing secret value")
	}

	var buffer bytes.Buffer
	buffer.WriteString(s.Name)
	buffer.WriteString("=")
	buffer.WriteString(s.Value)

	return buffer.String(), nil
}

//UnmarshalYAML parses the yaml element and sets the values of s; it will return an error if the parsing fails, or
//if the format is incorrect
func (s *Secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var secret string
	if err := unmarshal(&secret); err != nil {
		return err
	}

	parts := strings.Split(secret, "=")
	if len(parts) != 2 {
		return fmt.Errorf("Invalid secret syntax")
	}

	s.Name = parts[0]
	s.Value = parts[1]
	return nil
}
