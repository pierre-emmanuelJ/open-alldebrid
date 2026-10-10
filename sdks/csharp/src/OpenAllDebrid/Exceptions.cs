namespace OpenAllDebrid;

/// <summary>An error answered by the API ("status": "error"). Compare <see cref="Code"/> with <see cref="ErrorCodes"/>.</summary>
public sealed class AllDebridException : Exception
{
    public AllDebridException(string code, string apiMessage, string? token = null)
        : base($"{code}: {apiMessage}")
    {
        Code = code;
        ApiMessage = apiMessage;
        Token = token;
    }

    /// <summary>The error code, e.g. MAGNET_INVALID_ID.</summary>
    public string Code { get; }

    /// <summary>The message of the API.</summary>
    public string ApiMessage { get; }

    /// <summary>With AUTH_BLOCKED, the token to poll with <see cref="AllDebrid.VerificationStatusAsync"/>.</summary>
    public string? Token { get; }
}

/// <summary>An HTTP status other than 200, such as 429 when rate limited.</summary>
public sealed class AllDebridHttpException : Exception
{
    public AllDebridHttpException(int statusCode, AllDebridException? apiError = null)
        : base(apiError is null ? $"HTTP {statusCode}" : $"HTTP {statusCode}: {apiError.Message}", apiError)
    {
        StatusCode = statusCode;
        ApiError = apiError;
    }

    public int StatusCode { get; }

    /// <summary>The API error in the body, when there is one.</summary>
    public AllDebridException? ApiError { get; }

    /// <summary>Whether the API asked to slow down (HTTP 429 or 503).</summary>
    public bool IsRateLimited => StatusCode is 429 or 503;
}
