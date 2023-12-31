// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NodeItem node item
//
// swagger:model NodeItem
type NodeItem struct {

	// host
	Host string `json:"host,omitempty"`

	// max concurrent
	MaxConcurrent int32 `json:"max_concurrent"`

	// name
	Name string `json:"name,omitempty"`

	// pending
	Pending int32 `json:"pending"`

	// running
	Running int32 `json:"running"`
}

// Validate validates this node item
func (m *NodeItem) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this node item based on context it is used
func (m *NodeItem) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *NodeItem) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NodeItem) UnmarshalBinary(b []byte) error {
	var res NodeItem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
