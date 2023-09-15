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

// StorageCreateKoofrOtherStorageRequest storage create koofr other storage request
//
// swagger:model storage.createKoofrOtherStorageRequest
type StorageCreateKoofrOtherStorageRequest struct {

	// config for underlying HTTP client
	ClientConfig struct {
		ModelClientConfig
	} `json:"clientConfig,omitempty"`

	// config for the storage
	Config struct {
		StorageKoofrOtherConfig
	} `json:"config,omitempty"`

	// Name of the storage, must be unique
	// Example: my-storage
	Name string `json:"name,omitempty"`

	// Path of the storage
	Path string `json:"path,omitempty"`
}

// Validate validates this storage create koofr other storage request
func (m *StorageCreateKoofrOtherStorageRequest) Validate(formats strfmt.Registry) error {
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

func (m *StorageCreateKoofrOtherStorageRequest) validateClientConfig(formats strfmt.Registry) error {
	if swag.IsZero(m.ClientConfig) { // not required
		return nil
	}

	return nil
}

func (m *StorageCreateKoofrOtherStorageRequest) validateConfig(formats strfmt.Registry) error {
	if swag.IsZero(m.Config) { // not required
		return nil
	}

	return nil
}

// ContextValidate validate this storage create koofr other storage request based on the context it is used
func (m *StorageCreateKoofrOtherStorageRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
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

func (m *StorageCreateKoofrOtherStorageRequest) contextValidateClientConfig(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

func (m *StorageCreateKoofrOtherStorageRequest) contextValidateConfig(ctx context.Context, formats strfmt.Registry) error {

	return nil
}

// MarshalBinary interface implementation
func (m *StorageCreateKoofrOtherStorageRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StorageCreateKoofrOtherStorageRequest) UnmarshalBinary(b []byte) error {
	var res StorageCreateKoofrOtherStorageRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
