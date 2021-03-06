package config

import (
	"errors"
	"fmt"
	"github.com/blang/semver"
	"github.com/ghodss/yaml"
	"strconv"
)

const (
	maxCompatManifestVersion = "1.0.0"

	// StringType represents the representation of a string parameter type.
	StringType = "string"

	// IntegerType represents the representation of a integer parameter type.
	IntegerType = "integer"

	// BooleanType represents the representation of a boolean parameter type.
	BooleanType = "boolean"
)

// ManifestFile represents a template's metadata.
type ManifestFile struct {
	Version    string       `json:"version"`
	Parameters []*Parameter `json:"parameters"`
}

// Parameter represents a parameter defined as part of a template's metadata.
type Parameter struct {
	Name         string   `json:"name"`
	Prompt       string   `json:"prompt"`
	Type         string   `json:"type"`
	Enum         []string `json:"enum"`
	Description  string   `json:"description"`
	DefaultValue string   `json:"defaultValue"`
}

// LoadManifestData unmarshals YAML content into a ManifestFile.
func LoadManifestData(b []byte) (*ManifestFile, error) {
	m := &ManifestFile{}
	err := yaml.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// ValidateManifest validates the expected YAML structure of a manifest.
func ValidateManifest(m *ManifestFile) error {
	err := validateManifestVersion(m.Version)
	if err != nil {
		return err
	}
	err = validateManifestParams(m.Parameters)
	if err != nil {
		return err
	}
	return nil
}

func validateManifestVersion(version string) error {
	if version == "" {
		return errors.New("manifest file needs to provide a version")
	}

	v, err := semver.Make(version)
	if err != nil {
		return err
	}

	maxCompatVersion, err := semver.Make(maxCompatManifestVersion)
	if err != nil {
		return err
	}

	if v.GT(maxCompatVersion) {
		return fmt.Errorf("manifest version needs to be less than %s", maxCompatManifestVersion)
	}

	return nil
}

func validateManifestParams(params []*Parameter) error {
	for _, p := range params {
		if p.Type == "" {
			return errors.New("every parameter defined in manifest needs to provide a type")
		}
		if p.Type == IntegerType {
			if p.DefaultValue != "" {
				_, err := strconv.Atoi(p.DefaultValue)
				if err != nil {
					return err
				}
			}
			if p.Enum != nil {
				for _, e := range p.Enum {
					_, err := strconv.Atoi(e)
					if err != nil {
						return err
					}
				}
			}
		}
		if p.Type == BooleanType {
			if p.DefaultValue != "" {
				_, err := strconv.ParseBool(p.DefaultValue)
				if err != nil {
					return err
				}
			}
			if p.Enum != nil {
				return errors.New("boolean type does not allow enums")
			}
		}
	}
	return nil
}
