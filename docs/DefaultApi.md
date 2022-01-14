# {{classname}}

All URIs are relative to *https://api.alldebrid.com/v4*

Method | HTTP request | Description
------------- | ------------- | -------------
[**HostsGet**](DefaultApi.md#HostsGet) | **Get** /hosts | Use these endpoints to get a list of all the hosts we support, and the services of redirection (adf.ly, bit.ly) or protection we can extract links from.

# **HostsGet**
> InlineResponse200 HostsGet(ctx, agent, optional)
Use these endpoints to get a list of all the hosts we support, and the services of redirection (adf.ly, bit.ly) or protection we can extract links from.

Use this endpoint to retrieve informations about what hosts we support and all related informations about it.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **agent** | **string**| Your software user-agent. | 
 **optional** | ***DefaultApiHostsGetOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DefaultApiHostsGetOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **hostOnly** | **optional.String**| Endpoint will only return \&quot;hosts\&quot; data | 

### Return type

[**InlineResponse200**](inline_response_200.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

