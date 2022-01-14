# InlineResponse20014

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | Pointer to **string** |  | [optional] 
**Data** | Pointer to [**InlineResponse20014Data**](InlineResponse20014Data.md) |  | [optional] 
**Error** | Pointer to [**InlineResponse20013Error**](InlineResponse20013Error.md) |  | [optional] 

## Methods

### NewInlineResponse20014

`func NewInlineResponse20014() *InlineResponse20014`

NewInlineResponse20014 instantiates a new InlineResponse20014 object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInlineResponse20014WithDefaults

`func NewInlineResponse20014WithDefaults() *InlineResponse20014`

NewInlineResponse20014WithDefaults instantiates a new InlineResponse20014 object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *InlineResponse20014) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *InlineResponse20014) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *InlineResponse20014) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *InlineResponse20014) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetData

`func (o *InlineResponse20014) GetData() InlineResponse20014Data`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *InlineResponse20014) GetDataOk() (*InlineResponse20014Data, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *InlineResponse20014) SetData(v InlineResponse20014Data)`

SetData sets Data field to given value.

### HasData

`func (o *InlineResponse20014) HasData() bool`

HasData returns a boolean if a field has been set.

### GetError

`func (o *InlineResponse20014) GetError() InlineResponse20013Error`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *InlineResponse20014) GetErrorOk() (*InlineResponse20013Error, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *InlineResponse20014) SetError(v InlineResponse20013Error)`

SetError sets Error field to given value.

### HasError

`func (o *InlineResponse20014) HasError() bool`

HasError returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


