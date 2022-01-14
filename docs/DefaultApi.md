# \DefaultApi

All URIs are relative to *https://api.alldebrid.com/v4*

Method | HTTP request | Description
------------- | ------------- | -------------
[**HostsDomainsGet**](DefaultApi.md#HostsDomainsGet) | **Get** /hosts/domains | Use this endpoint to only retrieve the list of supported hosts domains and redirectors as an array.
[**HostsGet**](DefaultApi.md#HostsGet) | **Get** /hosts | Use this endpoint to retrieve informations about what hosts we support.
[**HostsPriorityGet**](DefaultApi.md#HostsPriorityGet) | **Get** /hosts/priority | Not all hosts are created equal, so some hosts are more limited than other.
[**LinkDelayedGet**](DefaultApi.md#LinkDelayedGet) | **Get** /link/delayed | This endpoint give the status of a delayed link.
[**LinkInfosGet**](DefaultApi.md#LinkInfosGet) | **Get** /link/infos | Use this endpoint to retrieve informations about a link.
[**LinkRedirectorGet**](DefaultApi.md#LinkRedirectorGet) | **Get** /link/redirector | Use this endpoint to retrieve links protected by a redirector or link protector.
[**LinkStreamingGet**](DefaultApi.md#LinkStreamingGet) | **Get** /link/streaming | The unlocking flow for streaming link is a bit more complex.
[**LinkUnlockGet**](DefaultApi.md#LinkUnlockGet) | **Get** /link/unlock | This endpoint unlocks a given link.
[**MagnetDeleteGet**](DefaultApi.md#MagnetDeleteGet) | **Get** /magnet/delete | Delete a magnet.
[**MagnetInstantGet**](DefaultApi.md#MagnetInstantGet) | **Get** /magnet/instant | Check if a magnet is available instantly.
[**MagnetRestartGet**](DefaultApi.md#MagnetRestartGet) | **Get** /magnet/restart | Restart a failed magnet, or multiple failed magnets at once.
[**MagnetStatusGet**](DefaultApi.md#MagnetStatusGet) | **Get** /magnet/status | Get the status of current magnets, or only one if you specify a magnet ID.
[**MagnetUploadFilePost**](DefaultApi.md#MagnetUploadFilePost) | **Post** /magnet/upload/file | Upload torrent files.
[**MagnetUploadGet**](DefaultApi.md#MagnetUploadGet) | **Get** /magnet/upload | Upload a magnet with its URI or hash.
[**UserGet**](DefaultApi.md#UserGet) | **Get** /user | Use this endpoint to get user informations.
[**UserHistoryDeleteGet**](DefaultApi.md#UserHistoryDeleteGet) | **Get** /user/history/delete | Use this endpoint to delete all links currently in your recent links history.
[**UserHistoryGet**](DefaultApi.md#UserHistoryGet) | **Get** /user/history | Use this endpoint to get recent links.
[**UserHostsGet**](DefaultApi.md#UserHostsGet) | **Get** /user/hosts | This endpoint retrieves a complete list of all available hosts for this user.
[**UserLinksDeleteGet**](DefaultApi.md#UserLinksDeleteGet) | **Get** /user/links/delete | Delete a saved link.
[**UserLinksGet**](DefaultApi.md#UserLinksGet) | **Get** /user/links | Use this endpoint to get links the user saved for later use.
[**UserLinksSaveGet**](DefaultApi.md#UserLinksSaveGet) | **Get** /user/links/save | Save a link.
[**UserNotificationClearGet**](DefaultApi.md#UserNotificationClearGet) | **Get** /user/notification/clear | This endpoint clears a user notification with its code.



## HostsDomainsGet

> InlineResponse2001 HostsDomainsGet(ctx).Agent(agent).Execute()

Use this endpoint to only retrieve the list of supported hosts domains and redirectors as an array.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.HostsDomainsGet(context.Background()).Agent(agent).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.HostsDomainsGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `HostsDomainsGet`: InlineResponse2001
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.HostsDomainsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiHostsDomainsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]

### Return type

[**InlineResponse2001**](InlineResponse2001.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## HostsGet

> InlineResponse200 HostsGet(ctx).Agent(agent).HostOnly(hostOnly).Execute()

Use this endpoint to retrieve informations about what hosts we support.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    hostOnly := "hostOnly_example" // string | Endpoint will only return \"hosts\" data (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.HostsGet(context.Background()).Agent(agent).HostOnly(hostOnly).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.HostsGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `HostsGet`: InlineResponse200
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.HostsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiHostsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **hostOnly** | **string** | Endpoint will only return \&quot;hosts\&quot; data | 

### Return type

[**InlineResponse200**](InlineResponse200.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## HostsPriorityGet

> InlineResponse2002 HostsPriorityGet(ctx).Agent(agent).Execute()

Not all hosts are created equal, so some hosts are more limited than other.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.HostsPriorityGet(context.Background()).Agent(agent).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.HostsPriorityGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `HostsPriorityGet`: InlineResponse2002
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.HostsPriorityGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiHostsPriorityGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]

### Return type

[**InlineResponse2002**](InlineResponse2002.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## LinkDelayedGet

> InlineResponse2009 LinkDelayedGet(ctx).Agent(agent).Id(id).Apikey(apikey).Execute()

This endpoint give the status of a delayed link.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    id := "id_example" // string | Delayed ID received in /link/unlock.
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.LinkDelayedGet(context.Background()).Agent(agent).Id(id).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.LinkDelayedGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `LinkDelayedGet`: InlineResponse2009
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.LinkDelayedGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLinkDelayedGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **id** | **string** | Delayed ID received in /link/unlock. | 
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse2009**](InlineResponse2009.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## LinkInfosGet

> InlineResponse2005 LinkInfosGet(ctx).Agent(agent).Link(link).Apikey(apikey).Password(password).Execute()

Use this endpoint to retrieve informations about a link.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    link := []string{"Inner_example"} // []string | The link or array of links you request informations about.
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)
    password := "password_example" // string | Link password (supported on uptobox / 1fichier). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.LinkInfosGet(context.Background()).Agent(agent).Link(link).Apikey(apikey).Password(password).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.LinkInfosGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `LinkInfosGet`: InlineResponse2005
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.LinkInfosGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLinkInfosGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **link** | **[]string** | The link or array of links you request informations about. | 
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 
 **password** | **string** | Link password (supported on uptobox / 1fichier). | 

### Return type

[**InlineResponse2005**](InlineResponse2005.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## LinkRedirectorGet

> InlineResponse2006 LinkRedirectorGet(ctx).Agent(agent).Link(link).Apikey(apikey).Execute()

Use this endpoint to retrieve links protected by a redirector or link protector.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    link := "link_example" // string | The redirector or protector link to extract links.
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.LinkRedirectorGet(context.Background()).Agent(agent).Link(link).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.LinkRedirectorGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `LinkRedirectorGet`: InlineResponse2006
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.LinkRedirectorGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLinkRedirectorGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **link** | **string** | The redirector or protector link to extract links. | 
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse2006**](InlineResponse2006.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## LinkStreamingGet

> InlineResponse2008 LinkStreamingGet(ctx).Agent(agent).Id(id).Stream(stream).Apikey(apikey).Execute()

The unlocking flow for streaming link is a bit more complex.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    id := "id_example" // string | The link ID you received from the /link/unlock call.
    stream := "stream_example" // string | The stream ID you choosed from the stream qualities list returned by /link/unlock.
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.LinkStreamingGet(context.Background()).Agent(agent).Id(id).Stream(stream).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.LinkStreamingGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `LinkStreamingGet`: InlineResponse2008
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.LinkStreamingGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLinkStreamingGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **id** | **string** | The link ID you received from the /link/unlock call. | 
 **stream** | **string** | The stream ID you choosed from the stream qualities list returned by /link/unlock. | 
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse2008**](InlineResponse2008.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## LinkUnlockGet

> InlineResponse2007 LinkUnlockGet(ctx).Agent(agent).Link(link).Apikey(apikey).Password(password).Execute()

This endpoint unlocks a given link.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    link := "link_example" // string | The redirector or protector link to extract links.
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)
    password := "password_example" // string | Link password (supported on uptobox / 1fichier). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.LinkUnlockGet(context.Background()).Agent(agent).Link(link).Apikey(apikey).Password(password).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.LinkUnlockGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `LinkUnlockGet`: InlineResponse2007
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.LinkUnlockGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLinkUnlockGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **link** | **string** | The redirector or protector link to extract links. | 
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 
 **password** | **string** | Link password (supported on uptobox / 1fichier). | 

### Return type

[**InlineResponse2007**](InlineResponse2007.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## MagnetDeleteGet

> InlineResponse20013 MagnetDeleteGet(ctx).Agent(agent).Id(id).Apikey(apikey).Execute()

Delete a magnet.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    id := "id_example" // string | Magnet ID.
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.MagnetDeleteGet(context.Background()).Agent(agent).Id(id).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.MagnetDeleteGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `MagnetDeleteGet`: InlineResponse20013
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.MagnetDeleteGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMagnetDeleteGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **id** | **string** | Magnet ID. | 
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse20013**](InlineResponse20013.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## MagnetInstantGet

> InlineResponse20015 MagnetInstantGet(ctx).Agent(agent).Magnets(magnets).Apikey(apikey).Execute()

Check if a magnet is available instantly.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    magnets := []string{"Inner_example"} // []string | Magnets URI or hash you wish to check instant availability, can be one or many links
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.MagnetInstantGet(context.Background()).Agent(agent).Magnets(magnets).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.MagnetInstantGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `MagnetInstantGet`: InlineResponse20015
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.MagnetInstantGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMagnetInstantGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **magnets** | **[]string** | Magnets URI or hash you wish to check instant availability, can be one or many links | 
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse20015**](InlineResponse20015.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## MagnetRestartGet

> InlineResponse20014 MagnetRestartGet(ctx).Agent(agent).Ids(ids).Apikey(apikey).Execute()

Restart a failed magnet, or multiple failed magnets at once.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    ids := []string{"Inner_example"} // []string | Array of Magnet ID.
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.MagnetRestartGet(context.Background()).Agent(agent).Ids(ids).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.MagnetRestartGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `MagnetRestartGet`: InlineResponse20014
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.MagnetRestartGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMagnetRestartGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **ids** | **[]string** | Array of Magnet ID. | 
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse20014**](InlineResponse20014.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## MagnetStatusGet

> InlineResponse20012 MagnetStatusGet(ctx).Agent(agent).Apikey(apikey).Id(id).Status(status).Session(session).Counter(counter).Execute()

Get the status of current magnets, or only one if you specify a magnet ID.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)
    id := "id_example" // string | Magnet ID. (optional)
    status := "status_example" // string | Magnets status filter. Either active, ready, expired or error (optional)
    session := "session_example" // string | Session ID for Live mode (see Live Mode). (optional)
    counter := "counter_example" // string | Counter for Live mode (see Live Mode). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.MagnetStatusGet(context.Background()).Agent(agent).Apikey(apikey).Id(id).Status(status).Session(session).Counter(counter).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.MagnetStatusGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `MagnetStatusGet`: InlineResponse20012
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.MagnetStatusGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMagnetStatusGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 
 **id** | **string** | Magnet ID. | 
 **status** | **string** | Magnets status filter. Either active, ready, expired or error | 
 **session** | **string** | Session ID for Live mode (see Live Mode). | 
 **counter** | **string** | Counter for Live mode (see Live Mode). | 

### Return type

[**InlineResponse20012**](InlineResponse20012.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## MagnetUploadFilePost

> InlineResponse20011 MagnetUploadFilePost(ctx).Agent(agent).Apikey(apikey).File(file).Execute()

Upload torrent files.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)
    file := os.NewFile(1234, "some_file") // *os.File |  (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.MagnetUploadFilePost(context.Background()).Agent(agent).Apikey(apikey).File(file).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.MagnetUploadFilePost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `MagnetUploadFilePost`: InlineResponse20011
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.MagnetUploadFilePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMagnetUploadFilePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 
 **file** | ***os.File** |  | 

### Return type

[**InlineResponse20011**](InlineResponse20011.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## MagnetUploadGet

> InlineResponse20010 MagnetUploadGet(ctx).Agent(agent).Magnets(magnets).Apikey(apikey).Execute()

Upload a magnet with its URI or hash.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    magnets := []string{"Inner_example"} // []string | Magnet(s) URI or hash. Must send magnet either in GET param or in POST data.
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.MagnetUploadGet(context.Background()).Agent(agent).Magnets(magnets).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.MagnetUploadGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `MagnetUploadGet`: InlineResponse20010
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.MagnetUploadGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMagnetUploadGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **magnets** | **[]string** | Magnet(s) URI or hash. Must send magnet either in GET param or in POST data. | 
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse20010**](InlineResponse20010.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserGet

> InlineResponse2003 UserGet(ctx).Agent(agent).Apikey(apikey).Execute()

Use this endpoint to get user informations.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.UserGet(context.Background()).Agent(agent).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.UserGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserGet`: InlineResponse2003
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.UserGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse2003**](InlineResponse2003.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserHistoryDeleteGet

> InlineResponse20019 UserHistoryDeleteGet(ctx).Agent(agent).Apikey(apikey).Execute()

Use this endpoint to delete all links currently in your recent links history.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.UserHistoryDeleteGet(context.Background()).Agent(agent).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.UserHistoryDeleteGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserHistoryDeleteGet`: InlineResponse20019
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.UserHistoryDeleteGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserHistoryDeleteGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse20019**](InlineResponse20019.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserHistoryGet

> InlineResponse20016 UserHistoryGet(ctx).Agent(agent).Apikey(apikey).Execute()

Use this endpoint to get recent links.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.UserHistoryGet(context.Background()).Agent(agent).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.UserHistoryGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserHistoryGet`: InlineResponse20016
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.UserHistoryGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserHistoryGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse20016**](InlineResponse20016.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserHostsGet

> InlineResponse200 UserHostsGet(ctx).Agent(agent).Apikey(apikey).HostOnly(hostOnly).Execute()

This endpoint retrieves a complete list of all available hosts for this user.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)
    hostOnly := "hostOnly_example" // string | Endpoint will only return \"hosts\" data (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.UserHostsGet(context.Background()).Agent(agent).Apikey(apikey).HostOnly(hostOnly).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.UserHostsGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserHostsGet`: InlineResponse200
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.UserHostsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserHostsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 
 **hostOnly** | **string** | Endpoint will only return \&quot;hosts\&quot; data | 

### Return type

[**InlineResponse200**](InlineResponse200.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserLinksDeleteGet

> InlineResponse20018 UserLinksDeleteGet(ctx).Agent(agent).Apikey(apikey).Link(link).Link2(link2).Execute()

Delete a saved link.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)
    link := "link_example" // string | Link to delete. (optional)
    link2 := []string{"Inner_example"} // []string | Links to delete. (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.UserLinksDeleteGet(context.Background()).Agent(agent).Apikey(apikey).Link(link).Link2(link2).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.UserLinksDeleteGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserLinksDeleteGet`: InlineResponse20018
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.UserLinksDeleteGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserLinksDeleteGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 
 **link** | **string** | Link to delete. | 
 **link2** | **[]string** | Links to delete. | 

### Return type

[**InlineResponse20018**](InlineResponse20018.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserLinksGet

> InlineResponse20016 UserLinksGet(ctx).Agent(agent).Apikey(apikey).Execute()

Use this endpoint to get links the user saved for later use.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.UserLinksGet(context.Background()).Agent(agent).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.UserLinksGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserLinksGet`: InlineResponse20016
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.UserLinksGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserLinksGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse20016**](InlineResponse20016.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserLinksSaveGet

> InlineResponse20017 UserLinksSaveGet(ctx).Agent(agent).Link(link).Apikey(apikey).Execute()

Save a link.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    link := []string{"Inner_example"} // []string | Links to save.
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.UserLinksSaveGet(context.Background()).Agent(agent).Link(link).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.UserLinksSaveGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserLinksSaveGet`: InlineResponse20017
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.UserLinksSaveGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserLinksSaveGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **link** | **[]string** | Links to save. | 
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse20017**](InlineResponse20017.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserNotificationClearGet

> InlineResponse2004 UserNotificationClearGet(ctx).Agent(agent).Code(code).Apikey(apikey).Execute()

This endpoint clears a user notification with its code.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    agent := "agent_example" // string | Your software user-agent. (default to "open-alldebrid")
    code := "code_example" // string | Notification code to clear
    apikey := "apikey_example" // string | Deprecated User apikey (Use Bearer Auth in header). (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.UserNotificationClearGet(context.Background()).Agent(agent).Code(code).Apikey(apikey).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.UserNotificationClearGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserNotificationClearGet`: InlineResponse2004
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.UserNotificationClearGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserNotificationClearGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **agent** | **string** | Your software user-agent. | [default to &quot;open-alldebrid&quot;]
 **code** | **string** | Notification code to clear | 
 **apikey** | **string** | Deprecated User apikey (Use Bearer Auth in header). | 

### Return type

[**InlineResponse2004**](InlineResponse2004.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

