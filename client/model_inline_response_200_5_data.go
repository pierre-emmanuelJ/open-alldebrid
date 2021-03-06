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

// InlineResponse2005Data struct for InlineResponse2005Data
type InlineResponse2005Data struct {
	// Array of info objects
	Infos *[]InlineResponse2005DataInfos `json:"infos,omitempty"`
}

// NewInlineResponse2005Data instantiates a new InlineResponse2005Data object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInlineResponse2005Data() *InlineResponse2005Data {
	this := InlineResponse2005Data{}
	return &this
}

// NewInlineResponse2005DataWithDefaults instantiates a new InlineResponse2005Data object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInlineResponse2005DataWithDefaults() *InlineResponse2005Data {
	this := InlineResponse2005Data{}
	return &this
}

// GetInfos returns the Infos field value if set, zero value otherwise.
func (o *InlineResponse2005Data) GetInfos() []InlineResponse2005DataInfos {
	if o == nil || o.Infos == nil {
		var ret []InlineResponse2005DataInfos
		return ret
	}
	return *o.Infos
}

// GetInfosOk returns a tuple with the Infos field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse2005Data) GetInfosOk() (*[]InlineResponse2005DataInfos, bool) {
	if o == nil || o.Infos == nil {
		return nil, false
	}
	return o.Infos, true
}

// HasInfos returns a boolean if a field has been set.
func (o *InlineResponse2005Data) HasInfos() bool {
	if o != nil && o.Infos != nil {
		return true
	}

	return false
}

// SetInfos gets a reference to the given []InlineResponse2005DataInfos and assigns it to the Infos field.
func (o *InlineResponse2005Data) SetInfos(v []InlineResponse2005DataInfos) {
	o.Infos = &v
}

func (o InlineResponse2005Data) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Infos != nil {
		toSerialize["infos"] = o.Infos
	}
	return json.Marshal(toSerialize)
}

type NullableInlineResponse2005Data struct {
	value *InlineResponse2005Data
	isSet bool
}

func (v NullableInlineResponse2005Data) Get() *InlineResponse2005Data {
	return v.value
}

func (v *NullableInlineResponse2005Data) Set(val *InlineResponse2005Data) {
	v.value = val
	v.isSet = true
}

func (v NullableInlineResponse2005Data) IsSet() bool {
	return v.isSet
}

func (v *NullableInlineResponse2005Data) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInlineResponse2005Data(val *InlineResponse2005Data) *NullableInlineResponse2005Data {
	return &NullableInlineResponse2005Data{value: val, isSet: true}
}

func (v NullableInlineResponse2005Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInlineResponse2005Data) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


