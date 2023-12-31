// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// UserItem user item
//
// swagger:model UserItem
type UserItem struct {

	// avatar
	Avatar string `json:"avatar,omitempty"`

	// created
	Created string `json:"created,omitempty"`

	// enable
	Enable bool `json:"enable"`

	// id
	ID string `json:"id,omitempty"`

	// image count
	ImageCount int32 `json:"image_count"`

	// is private
	IsPrivate bool `json:"is_private"`

	// roles
	Roles string `json:"roles,omitempty"`

	// stable config
	StableConfig interface{} `json:"stable_config,omitempty"`

	// username
	Username string `json:"username,omitempty"`
}

// Validate validates this user item
func (m *UserItem) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this user item based on context it is used
func (m *UserItem) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *UserItem) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UserItem) UnmarshalBinary(b []byte) error {
	var res UserItem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
