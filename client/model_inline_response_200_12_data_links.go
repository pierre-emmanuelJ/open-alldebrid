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

// InlineResponse20012DataLinks struct for InlineResponse20012DataLinks
type InlineResponse20012DataLinks struct {
	Link *string `json:"link,omitempty"`
	Filename *string `json:"filename,omitempty"`
	Size *int64 `json:"size,omitempty"`
}

// NewInlineResponse20012DataLinks instantiates a new InlineResponse20012DataLinks object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInlineResponse20012DataLinks() *InlineResponse20012DataLinks {
	this := InlineResponse20012DataLinks{}
	return &this
}

// NewInlineResponse20012DataLinksWithDefaults instantiates a new InlineResponse20012DataLinks object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInlineResponse20012DataLinksWithDefaults() *InlineResponse20012DataLinks {
	this := InlineResponse20012DataLinks{}
	return &this
}

// GetLink returns the Link field value if set, zero value otherwise.
func (o *InlineResponse20012DataLinks) GetLink() string {
	if o == nil || o.Link == nil {
		var ret string
		return ret
	}
	return *o.Link
}

// GetLinkOk returns a tuple with the Link field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse20012DataLinks) GetLinkOk() (*string, bool) {
	if o == nil || o.Link == nil {
		return nil, false
	}
	return o.Link, true
}

// HasLink returns a boolean if a field has been set.
func (o *InlineResponse20012DataLinks) HasLink() bool {
	if o != nil && o.Link != nil {
		return true
	}

	return false
}

// SetLink gets a reference to the given string and assigns it to the Link field.
func (o *InlineResponse20012DataLinks) SetLink(v string) {
	o.Link = &v
}

// GetFilename returns the Filename field value if set, zero value otherwise.
func (o *InlineResponse20012DataLinks) GetFilename() string {
	if o == nil || o.Filename == nil {
		var ret string
		return ret
	}
	return *o.Filename
}

// GetFilenameOk returns a tuple with the Filename field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse20012DataLinks) GetFilenameOk() (*string, bool) {
	if o == nil || o.Filename == nil {
		return nil, false
	}
	return o.Filename, true
}

// HasFilename returns a boolean if a field has been set.
func (o *InlineResponse20012DataLinks) HasFilename() bool {
	if o != nil && o.Filename != nil {
		return true
	}

	return false
}

// SetFilename gets a reference to the given string and assigns it to the Filename field.
func (o *InlineResponse20012DataLinks) SetFilename(v string) {
	o.Filename = &v
}

// GetSize returns the Size field value if set, zero value otherwise.
func (o *InlineResponse20012DataLinks) GetSize() int64 {
	if o == nil || o.Size == nil {
		var ret int64
		return ret
	}
	return *o.Size
}

// GetSizeOk returns a tuple with the Size field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse20012DataLinks) GetSizeOk() (*int64, bool) {
	if o == nil || o.Size == nil {
		return nil, false
	}
	return o.Size, true
}

// HasSize returns a boolean if a field has been set.
func (o *InlineResponse20012DataLinks) HasSize() bool {
	if o != nil && o.Size != nil {
		return true
	}

	return false
}

// SetSize gets a reference to the given int64 and assigns it to the Size field.
func (o *InlineResponse20012DataLinks) SetSize(v int64) {
	o.Size = &v
}

func (o InlineResponse20012DataLinks) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Link != nil {
		toSerialize["link"] = o.Link
	}
	if o.Filename != nil {
		toSerialize["filename"] = o.Filename
	}
	if o.Size != nil {
		toSerialize["size"] = o.Size
	}
	return json.Marshal(toSerialize)
}

type NullableInlineResponse20012DataLinks struct {
	value *InlineResponse20012DataLinks
	isSet bool
}

func (v NullableInlineResponse20012DataLinks) Get() *InlineResponse20012DataLinks {
	return v.value
}

func (v *NullableInlineResponse20012DataLinks) Set(val *InlineResponse20012DataLinks) {
	v.value = val
	v.isSet = true
}

func (v NullableInlineResponse20012DataLinks) IsSet() bool {
	return v.isSet
}

func (v *NullableInlineResponse20012DataLinks) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInlineResponse20012DataLinks(val *InlineResponse20012DataLinks) *NullableInlineResponse20012DataLinks {
	return &NullableInlineResponse20012DataLinks{value: val, isSet: true}
}

func (v NullableInlineResponse20012DataLinks) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInlineResponse20012DataLinks) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


