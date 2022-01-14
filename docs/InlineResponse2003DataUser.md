# InlineResponse2003DataUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Username** | Pointer to **string** |  | [optional] 
**Email** | Pointer to **string** |  | [optional] 
**IsPremium** | Pointer to **bool** |  | [optional] 
**IsSubscribed** | Pointer to **bool** |  | [optional] 
**IsTrial** | Pointer to **bool** |  | [optional] 
**PremiumUntil** | Pointer to **int64** |  | [optional] 
**Lang** | Pointer to **string** |  | [optional] 
**PreferedDomain** | Pointer to **string** |  | [optional] 
**FidelityPoints** | Pointer to **int64** |  | [optional] 
**LimitedHostersQuotas** | Pointer to **map[string]int64** |  | [optional] 
**Notifications** | Pointer to **[]string** |  | [optional] 

## Methods

### NewInlineResponse2003DataUser

`func NewInlineResponse2003DataUser() *InlineResponse2003DataUser`

NewInlineResponse2003DataUser instantiates a new InlineResponse2003DataUser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInlineResponse2003DataUserWithDefaults

`func NewInlineResponse2003DataUserWithDefaults() *InlineResponse2003DataUser`

NewInlineResponse2003DataUserWithDefaults instantiates a new InlineResponse2003DataUser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUsername

`func (o *InlineResponse2003DataUser) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *InlineResponse2003DataUser) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *InlineResponse2003DataUser) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *InlineResponse2003DataUser) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

### GetEmail

`func (o *InlineResponse2003DataUser) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *InlineResponse2003DataUser) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *InlineResponse2003DataUser) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *InlineResponse2003DataUser) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetIsPremium

`func (o *InlineResponse2003DataUser) GetIsPremium() bool`

GetIsPremium returns the IsPremium field if non-nil, zero value otherwise.

### GetIsPremiumOk

`func (o *InlineResponse2003DataUser) GetIsPremiumOk() (*bool, bool)`

GetIsPremiumOk returns a tuple with the IsPremium field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPremium

`func (o *InlineResponse2003DataUser) SetIsPremium(v bool)`

SetIsPremium sets IsPremium field to given value.

### HasIsPremium

`func (o *InlineResponse2003DataUser) HasIsPremium() bool`

HasIsPremium returns a boolean if a field has been set.

### GetIsSubscribed

`func (o *InlineResponse2003DataUser) GetIsSubscribed() bool`

GetIsSubscribed returns the IsSubscribed field if non-nil, zero value otherwise.

### GetIsSubscribedOk

`func (o *InlineResponse2003DataUser) GetIsSubscribedOk() (*bool, bool)`

GetIsSubscribedOk returns a tuple with the IsSubscribed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsSubscribed

`func (o *InlineResponse2003DataUser) SetIsSubscribed(v bool)`

SetIsSubscribed sets IsSubscribed field to given value.

### HasIsSubscribed

`func (o *InlineResponse2003DataUser) HasIsSubscribed() bool`

HasIsSubscribed returns a boolean if a field has been set.

### GetIsTrial

`func (o *InlineResponse2003DataUser) GetIsTrial() bool`

GetIsTrial returns the IsTrial field if non-nil, zero value otherwise.

### GetIsTrialOk

`func (o *InlineResponse2003DataUser) GetIsTrialOk() (*bool, bool)`

GetIsTrialOk returns a tuple with the IsTrial field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsTrial

`func (o *InlineResponse2003DataUser) SetIsTrial(v bool)`

SetIsTrial sets IsTrial field to given value.

### HasIsTrial

`func (o *InlineResponse2003DataUser) HasIsTrial() bool`

HasIsTrial returns a boolean if a field has been set.

### GetPremiumUntil

`func (o *InlineResponse2003DataUser) GetPremiumUntil() int64`

GetPremiumUntil returns the PremiumUntil field if non-nil, zero value otherwise.

### GetPremiumUntilOk

`func (o *InlineResponse2003DataUser) GetPremiumUntilOk() (*int64, bool)`

GetPremiumUntilOk returns a tuple with the PremiumUntil field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPremiumUntil

`func (o *InlineResponse2003DataUser) SetPremiumUntil(v int64)`

SetPremiumUntil sets PremiumUntil field to given value.

### HasPremiumUntil

`func (o *InlineResponse2003DataUser) HasPremiumUntil() bool`

HasPremiumUntil returns a boolean if a field has been set.

### GetLang

`func (o *InlineResponse2003DataUser) GetLang() string`

GetLang returns the Lang field if non-nil, zero value otherwise.

### GetLangOk

`func (o *InlineResponse2003DataUser) GetLangOk() (*string, bool)`

GetLangOk returns a tuple with the Lang field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLang

`func (o *InlineResponse2003DataUser) SetLang(v string)`

SetLang sets Lang field to given value.

### HasLang

`func (o *InlineResponse2003DataUser) HasLang() bool`

HasLang returns a boolean if a field has been set.

### GetPreferedDomain

`func (o *InlineResponse2003DataUser) GetPreferedDomain() string`

GetPreferedDomain returns the PreferedDomain field if non-nil, zero value otherwise.

### GetPreferedDomainOk

`func (o *InlineResponse2003DataUser) GetPreferedDomainOk() (*string, bool)`

GetPreferedDomainOk returns a tuple with the PreferedDomain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreferedDomain

`func (o *InlineResponse2003DataUser) SetPreferedDomain(v string)`

SetPreferedDomain sets PreferedDomain field to given value.

### HasPreferedDomain

`func (o *InlineResponse2003DataUser) HasPreferedDomain() bool`

HasPreferedDomain returns a boolean if a field has been set.

### GetFidelityPoints

`func (o *InlineResponse2003DataUser) GetFidelityPoints() int64`

GetFidelityPoints returns the FidelityPoints field if non-nil, zero value otherwise.

### GetFidelityPointsOk

`func (o *InlineResponse2003DataUser) GetFidelityPointsOk() (*int64, bool)`

GetFidelityPointsOk returns a tuple with the FidelityPoints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFidelityPoints

`func (o *InlineResponse2003DataUser) SetFidelityPoints(v int64)`

SetFidelityPoints sets FidelityPoints field to given value.

### HasFidelityPoints

`func (o *InlineResponse2003DataUser) HasFidelityPoints() bool`

HasFidelityPoints returns a boolean if a field has been set.

### GetLimitedHostersQuotas

`func (o *InlineResponse2003DataUser) GetLimitedHostersQuotas() map[string]int64`

GetLimitedHostersQuotas returns the LimitedHostersQuotas field if non-nil, zero value otherwise.

### GetLimitedHostersQuotasOk

`func (o *InlineResponse2003DataUser) GetLimitedHostersQuotasOk() (*map[string]int64, bool)`

GetLimitedHostersQuotasOk returns a tuple with the LimitedHostersQuotas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLimitedHostersQuotas

`func (o *InlineResponse2003DataUser) SetLimitedHostersQuotas(v map[string]int64)`

SetLimitedHostersQuotas sets LimitedHostersQuotas field to given value.

### HasLimitedHostersQuotas

`func (o *InlineResponse2003DataUser) HasLimitedHostersQuotas() bool`

HasLimitedHostersQuotas returns a boolean if a field has been set.

### GetNotifications

`func (o *InlineResponse2003DataUser) GetNotifications() []string`

GetNotifications returns the Notifications field if non-nil, zero value otherwise.

### GetNotificationsOk

`func (o *InlineResponse2003DataUser) GetNotificationsOk() (*[]string, bool)`

GetNotificationsOk returns a tuple with the Notifications field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotifications

`func (o *InlineResponse2003DataUser) SetNotifications(v []string)`

SetNotifications sets Notifications field to given value.

### HasNotifications

`func (o *InlineResponse2003DataUser) HasNotifications() bool`

HasNotifications returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


