// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// StorageCreateOosEnvAuthStorageRequest storage create oos env auth storage request
//
// swagger:model storage.createOosEnv_authStorageRequest
type StorageCreateOosEnvAuthStorageRequest struct {

	// config for underlying HTTP client
	ClientConfig struct {
		ModelClientConfig
	} `json:"clientConfig,omitempty"`

	// config for the storage
	Config struct {
		StorageOosEnvAuthConfig
	} `json:"config,omitempty"`

	// Name of the storage, must be unique
	// Example: my-storage
	Name string `json:"name,omitempty"`

	// Path of the storage
	Path string `json:"path,omitempty"`
}

// Validate validates this storage create oos env auth storage request
func (m *StorageCreateOosEnvAuthStorageRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateClientConfig(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateConfig(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StorageCreateOosEnvAuthStorageRequest) validateClientConfig(formats strfmt.Registry) error {
	if swag.IsZero(m.ClientConfig) { // not required
		return nil
	}

	return nil
}

func (m *StorageCreateOosEnvAuthStorageRequest) validateConfig(formats strfmt.Registry) error {
	if swag.IsZero(m.Config) { // not required
		return nil
	}

	return nil
}

// ContextValidate validate this storage create oos env auth storage request based on the context it is used
func (m *StorageCreateOosEnvAuthStorageRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateClientConfig(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateConfig(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StorageCreateOosEnvAuthStorageRequest) contextValidateClientConfig(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

func (m *StorageCreateOosEnvAuthStorageRequest) contextValidateConfig(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

// MarshalBinary interface implementation
func (m *StorageCreateOosEnvAuthStorageRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StorageCreateOosEnvAuthStorageRequest) UnmarshalBinary(b []byte) error {
	var res StorageCreateOosEnvAuthStorageRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
