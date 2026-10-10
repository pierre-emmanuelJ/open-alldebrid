# OpenAllDebrid (C#)

A .NET client for the [AllDebrid](https://alldebrid.com) API, built on the
[open-alldebrid OpenAPI spec](https://github.com/pierre-emmanuelJ/open-alldebrid/blob/master/openapi.yaml).
.NET 8 and 10, async with cancellation, no dependencies beyond the base
library (HttpClient and System.Text.Json).

```sh
dotnet add package OpenAllDebrid
```

```csharp
using OpenAllDebrid;

var client = new AllDebrid(apiKey, userAgent: "my-app");

var user = await client.UserAsync();

var uploaded = await client.UploadMagnetsAsync(new[] { "magnet:?xt=urn:btih:..." });
var magnet = await client.MagnetAsync(uploaded[0].Id!.Value);
var files = await client.MagnetFilesAsync(new[] { magnet.Id });

try
{
    var unlocked = await client.UnlockAsync("https://1fichier.com/?abcdefghij");
    var link = unlocked.Delayed is { } id ? await client.WaitDelayedLinkAsync(id.Value) : unlocked.Link;
}
catch (AllDebridException e) when (e.Code == ErrorCodes.LinkDown)
{
    Console.WriteLine("the file is gone");
}
```

## What it handles for you

- **Errors.** The API answers errors with HTTP 200 and `"status": "error"`.
  Every call throws them as `AllDebridException` (with `Code`, and `Token` for
  `AUTH_BLOCKED`); other HTTP statuses throw `AllDebridHttpException`
  (`IsRateLimited` for 429/503).
- **List parameters** are sent the way the API expects:
  `magnets[]=a&magnets[]=b`.
- **Inconsistent shapes.** IDs sent as strings or numbers (`FlexString`), a
  magnet status that is an object for one magnet and a list otherwise, regexps
  that are a string or a list, and empty quota maps sent as `[]` all decode.
  Responses missing a required field throw instead of giving half-empty objects.
- **PIN login.** `GetPinAsync()` then `WaitForPinAsync(pin)` returns the new key.
- **Delayed links.** `WaitDelayedLinkAsync(id)` polls until the link is ready.
- **Live mode.** `LiveMagnets().PollAsync()` only transfers what changed and
  keeps the merged list, including deletions.

## Try it without an account

`AllDebrid.DemoPremiumKey`, `DemoTrialKey` and `DemoFreeKey` make the API answer
static sample data.

## Development

`src/OpenAllDebrid/Models.g.cs` is generated from the spec by
`scripts/gen_models.py` at the root of the repository (needs python3 with
PyYAML):

```sh
python3 ../../scripts/gen_models.py csharp src/OpenAllDebrid/Models.g.cs
dotnet test --solution OpenAllDebrid.slnx   # OFFLINE=1 skips the demo API calls
```

The tests check that every operation of the spec has a method calling the
right path, and that `ErrorCodes` matches the spec.
