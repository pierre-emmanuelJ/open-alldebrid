import Foundation

/// An error of the client.
public enum AllDebridError: Error, Sendable, Equatable, CustomStringConvertible {
    /// An error answered by the API (`"status": "error"`). Compare `code` with ``ErrorCode``.
    /// With `AUTH_BLOCKED`, `token` is the token to poll with ``AllDebrid/verificationStatus(token:)``.
    case api(code: String, message: String, token: String?)
    /// An HTTP status other than 200, such as 429 when rate limited.
    case http(status: Int, code: String?)
    /// The response did not have the expected shape.
    case decoding(String)

    /// The API error code, when this error carries one.
    public var code: String? {
        switch self {
        case .api(let code, _, _): code
        case .http(_, let code): code
        case .decoding: nil
        }
    }

    /// Whether the API asked to slow down (HTTP 429 or 503).
    public var isRateLimited: Bool {
        if case .http(let status, _) = self { return status == 429 || status == 503 }
        return false
    }

    public var description: String {
        switch self {
        case .api(let code, let message, _): "alldebrid: \(code): \(message)"
        case .http(let status, let code): "alldebrid: HTTP \(status)\(code.map { ": \($0)" } ?? "")"
        case .decoding(let detail): "alldebrid: unexpected response: \(detail)"
        }
    }
}
