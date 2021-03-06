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

// InlineResponse2007Data struct for InlineResponse2007Data
type InlineResponse2007Data struct {
	Filename *string `json:"filename,omitempty"`
	Filesize *int64 `json:"filesize,omitempty"`
	Host *string `json:"host,omitempty"`
	HostDomain *string `json:"hostDomain,omitempty"`
	Id *string `json:"id,omitempty"`
	Link *string `json:"link,omitempty"`
	Paws *bool `json:"paws,omitempty"`
	Streams *[]InlineResponse2007DataStreams `json:"streams,omitempty"`
}

// NewInlineResponse2007Data instantiates a new InlineResponse2007Data object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInlineResponse2007Data() *InlineResponse2007Data {
	this := InlineResponse2007Data{}
	return &this
}

// NewInlineResponse2007DataWithDefaults instantiates a new InlineResponse2007Data object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInlineResponse2007DataWithDefaults() *InlineResponse2007Data {
	this := InlineResponse2007Data{}
	return &this
}

// GetFilename returns the Filename field value if set, zero value otherwise.
func (o *InlineResponse2007Data) GetFilename() string {
	if o == nil || o.Filename == nil {
		var ret string
		return ret
	}
	return *o.Filename
}

// GetFilenameOk returns a tuple with the Filename field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse2007Data) GetFilenameOk() (*string, bool) {
	if o == nil || o.Filename == nil {
		return nil, false
	}
	return o.Filename, true
}

// HasFilename returns a boolean if a field has been set.
func (o *InlineResponse2007Data) HasFilename() bool {
	if o != nil && o.Filename != nil {
		return true
	}

	return false
}

// SetFilename gets a reference to the given string and assigns it to the Filename field.
func (o *InlineResponse2007Data) SetFilename(v string) {
	o.Filename = &v
}

// GetFilesize returns the Filesize field value if set, zero value otherwise.
func (o *InlineResponse2007Data) GetFilesize() int64 {
	if o == nil || o.Filesize == nil {
		var ret int64
		return ret
	}
	return *o.Filesize
}

// GetFilesizeOk returns a tuple with the Filesize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse2007Data) GetFilesizeOk() (*int64, bool) {
	if o == nil || o.Filesize == nil {
		return nil, false
	}
	return o.Filesize, true
}

// HasFilesize returns a boolean if a field has been set.
func (o *InlineResponse2007Data) HasFilesize() bool {
	if o != nil && o.Filesize != nil {
		return true
	}

	return false
}

// SetFilesize gets a reference to the given int64 and assigns it to the Filesize field.
func (o *InlineResponse2007Data) SetFilesize(v int64) {
	o.Filesize = &v
}

// GetHost returns the Host field value if set, zero value otherwise.
func (o *InlineResponse2007Data) GetHost() string {
	if o == nil || o.Host == nil {
		var ret string
		return ret
	}
	return *o.Host
}

// GetHostOk returns a tuple with the Host field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse2007Data) GetHostOk() (*string, bool) {
	if o == nil || o.Host == nil {
		return nil, false
	}
	return o.Host, true
}

// HasHost returns a boolean if a field has been set.
func (o *InlineResponse2007Data) HasHost() bool {
	if o != nil && o.Host != nil {
		return true
	}

	return false
}

// SetHost gets a reference to the given string and assigns it to the Host field.
func (o *InlineResponse2007Data) SetHost(v string) {
	o.Host = &v
}

// GetHostDomain returns the HostDomain field value if set, zero value otherwise.
func (o *InlineResponse2007Data) GetHostDomain() string {
	if o == nil || o.HostDomain == nil {
		var ret string
		return ret
	}
	return *o.HostDomain
}

// GetHostDomainOk returns a tuple with the HostDomain field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse2007Data) GetHostDomainOk() (*string, bool) {
	if o == nil || o.HostDomain == nil {
		return nil, false
	}
	return o.HostDomain, true
}

// HasHostDomain returns a boolean if a field has been set.
func (o *InlineResponse2007Data) HasHostDomain() bool {
	if o != nil && o.HostDomain != nil {
		return true
	}

	return false
}

// SetHostDomain gets a reference to the given string and assigns it to the HostDomain field.
func (o *InlineResponse2007Data) SetHostDomain(v string) {
	o.HostDomain = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *InlineResponse2007Data) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse2007Data) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *InlineResponse2007Data) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *InlineResponse2007Data) SetId(v string) {
	o.Id = &v
}

// GetLink returns the Link field value if set, zero value otherwise.
func (o *InlineResponse2007Data) GetLink() string {
	if o == nil || o.Link == nil {
		var ret string
		return ret
	}
	return *o.Link
}

// GetLinkOk returns a tuple with the Link field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse2007Data) GetLinkOk() (*string, bool) {
	if o == nil || o.Link == nil {
		return nil, false
	}
	return o.Link, true
}

// HasLink returns a boolean if a field has been set.
func (o *InlineResponse2007Data) HasLink() bool {
	if o != nil && o.Link != nil {
		return true
	}

	return false
}

// SetLink gets a reference to the given string and assigns it to the Link field.
func (o *InlineResponse2007Data) SetLink(v string) {
	o.Link = &v
}

// GetPaws returns the Paws field value if set, zero value otherwise.
func (o *InlineResponse2007Data) GetPaws() bool {
	if o == nil || o.Paws == nil {
		var ret bool
		return ret
	}
	return *o.Paws
}

// GetPawsOk returns a tuple with the Paws field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse2007Data) GetPawsOk() (*bool, bool) {
	if o == nil || o.Paws == nil {
		return nil, false
	}
	return o.Paws, true
}

// HasPaws returns a boolean if a field has been set.
func (o *InlineResponse2007Data) HasPaws() bool {
	if o != nil && o.Paws != nil {
		return true
	}

	return false
}

// SetPaws gets a reference to the given bool and assigns it to the Paws field.
func (o *InlineResponse2007Data) SetPaws(v bool) {
	o.Paws = &v
}

// GetStreams returns the Streams field value if set, zero value otherwise.
func (o *InlineResponse2007Data) GetStreams() []InlineResponse2007DataStreams {
	if o == nil || o.Streams == nil {
		var ret []InlineResponse2007DataStreams
		return ret
	}
	return *o.Streams
}

// GetStreamsOk returns a tuple with the Streams field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse2007Data) GetStreamsOk() (*[]InlineResponse2007DataStreams, bool) {
	if o == nil || o.Streams == nil {
		return nil, false
	}
	return o.Streams, true
}

// HasStreams returns a boolean if a field has been set.
func (o *InlineResponse2007Data) HasStreams() bool {
	if o != nil && o.Streams != nil {
		return true
	}

	return false
}

// SetStreams gets a reference to the given []InlineResponse2007DataStreams and assigns it to the Streams field.
func (o *InlineResponse2007Data) SetStreams(v []InlineResponse2007DataStreams) {
	o.Streams = &v
}

func (o InlineResponse2007Data) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Filename != nil {
		toSerialize["filename"] = o.Filename
	}
	if o.Filesize != nil {
		toSerialize["filesize"] = o.Filesize
	}
	if o.Host != nil {
		toSerialize["host"] = o.Host
	}
	if o.HostDomain != nil {
		toSerialize["hostDomain"] = o.HostDomain
	}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.Link != nil {
		toSerialize["link"] = o.Link
	}
	if o.Paws != nil {
		toSerialize["paws"] = o.Paws
	}
	if o.Streams != nil {
		toSerialize["streams"] = o.Streams
	}
	return json.Marshal(toSerialize)
}

type NullableInlineResponse2007Data struct {
	value *InlineResponse2007Data
	isSet bool
}

func (v NullableInlineResponse2007Data) Get() *InlineResponse2007Data {
	return v.value
}

func (v *NullableInlineResponse2007Data) Set(val *InlineResponse2007Data) {
	v.value = val
	v.isSet = true
}

func (v NullableInlineResponse2007Data) IsSet() bool {
	return v.isSet
}

func (v *NullableInlineResponse2007Data) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInlineResponse2007Data(val *InlineResponse2007Data) *NullableInlineResponse2007Data {
	return &NullableInlineResponse2007Data{value: val, isSet: true}
}

func (v NullableInlineResponse2007Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInlineResponse2007Data) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


