/*
Alldebrid API

Welcome to the OpenAPI Alldebrid API v4 !<br /> You can use this API to access various Alldebrid services from custom applications or scripts.<br /> <br /> API is organized around REST,<br /> returns JSON-encoded responses and use standard HTTP response codes.<br /> <br /> All calls are to be made on the HTTPS endpoints.<br /> Some are public, others require to be authentificated with an apikey (see Authentication).<br /> <br /> You must also identify your apps or script with a meaningfull agent parameter.<br /> <br /> This API version is namespaced as v4, as such all endpoint start with /v4/,<br /> such like http://api.alldebrid.com/v4/ping?agent=apiShowcase.<br /> <br /> This API v4 should be the final version regarding general response format and errors (hopefully).<br />

API version: 4.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
)

// InlineResponse20010Data struct for InlineResponse20010Data
type InlineResponse20010Data struct {
	Magnets *[]InlineResponse20010DataMagnets `json:"magnets,omitempty"`
}

// NewInlineResponse20010Data instantiates a new InlineResponse20010Data object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInlineResponse20010Data() *InlineResponse20010Data {
	this := InlineResponse20010Data{}
	return &this
}

// NewInlineResponse20010DataWithDefaults instantiates a new InlineResponse20010Data object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInlineResponse20010DataWithDefaults() *InlineResponse20010Data {
	this := InlineResponse20010Data{}
	return &this
}

// GetMagnets returns the Magnets field value if set, zero value otherwise.
func (o *InlineResponse20010Data) GetMagnets() []InlineResponse20010DataMagnets {
	if o == nil || o.Magnets == nil {
		var ret []InlineResponse20010DataMagnets
		return ret
	}
	return *o.Magnets
}

// GetMagnetsOk returns a tuple with the Magnets field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse20010Data) GetMagnetsOk() (*[]InlineResponse20010DataMagnets, bool) {
	if o == nil || o.Magnets == nil {
		return nil, false
	}
	return o.Magnets, true
}

// HasMagnets returns a boolean if a field has been set.
func (o *InlineResponse20010Data) HasMagnets() bool {
	if o != nil && o.Magnets != nil {
		return true
	}

	return false
}

// SetMagnets gets a reference to the given []InlineResponse20010DataMagnets and assigns it to the Magnets field.
func (o *InlineResponse20010Data) SetMagnets(v []InlineResponse20010DataMagnets) {
	o.Magnets = &v
}

func (o InlineResponse20010Data) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Magnets != nil {
		toSerialize["magnets"] = o.Magnets
	}
	return json.Marshal(toSerialize)
}

type NullableInlineResponse20010Data struct {
	value *InlineResponse20010Data
	isSet bool
}

func (v NullableInlineResponse20010Data) Get() *InlineResponse20010Data {
	return v.value
}

func (v *NullableInlineResponse20010Data) Set(val *InlineResponse20010Data) {
	v.value = val
	v.isSet = true
}

func (v NullableInlineResponse20010Data) IsSet() bool {
	return v.isSet
}

func (v *NullableInlineResponse20010Data) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInlineResponse20010Data(val *InlineResponse20010Data) *NullableInlineResponse20010Data {
	return &NullableInlineResponse20010Data{value: val, isSet: true}
}

func (v NullableInlineResponse20010Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInlineResponse20010Data) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

