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

// StorageCreateAzureblobStorageRequest storage create azureblob storage request
//
// swagger:model storage.createAzureblobStorageRequest
type StorageCreateAzureblobStorageRequest struct {

	// config for underlying HTTP client
	ClientConfig struct {
		ModelClientConfig
	} `json:"clientConfig,omitempty"`

	// config for the storage
	Config struct {
		StorageAzureblobConfig
	} `json:"config,omitempty"`

	// Name of the storage, must be unique
	// Example: my-storage
	Name string `json:"name,omitempty"`

	// Path of the storage
	Path string `json:"path,omitempty"`
}

// Validate validates this storage create azureblob storage request
func (m *StorageCreateAzureblobStorageRequest) Validate(formats strfmt.Registry) error {
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

func (m *StorageCreateAzureblobStorageRequest) validateClientConfig(formats strfmt.Registry) error {
	if swag.IsZero(m.ClientConfig) { // not required
		return nil
	}

	return nil
}

func (m *StorageCreateAzureblobStorageRequest) validateConfig(formats strfmt.Registry) error {
	if swag.IsZero(m.Config) { // not required
		return nil
	}

	return nil
}

// ContextValidate validate this storage create azureblob storage request based on the context it is used
func (m *StorageCreateAzureblobStorageRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
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

func (m *StorageCreateAzureblobStorageRequest) contextValidateClientConfig(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

func (m *StorageCreateAzureblobStorageRequest) contextValidateConfig(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

// MarshalBinary interface implementation
func (m *StorageCreateAzureblobStorageRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StorageCreateAzureblobStorageRequest) UnmarshalBinary(b []byte) error {
	var res StorageCreateAzureblobStorageRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
