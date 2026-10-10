import Foundation

/// A value the API sends either as a JSON string or as a JSON number, such as a
/// delayed link ID. It keeps the value as a string.
public struct FlexString: Decodable, Sendable, Equatable, Hashable, CustomStringConvertible {
    /// The value as sent by the API.
    public let value: String

    public init(_ value: String) {
        self.value = value
    }

    public init(from decoder: any Decoder) throws {
        let c = try decoder.singleValueContainer()
        if let s = try? c.decode(String.self) {
            value = s
        } else if let n = try? c.decode(Int64.self) {
            value = String(n)
        } else if let d = try? c.decode(Double.self) {
            value = String(d)
        } else {
            throw DecodingError.typeMismatch(
                FlexString.self,
                .init(codingPath: decoder.codingPath, debugDescription: "expected a string or a number"))
        }
    }

    /// The value as an integer, when it is one.
    public var int64Value: Int64? { Int64(value) }

    public var description: String { value }
}

/// Decodes a value that is either one item or a list of items.
private enum OneOrMany<T: Decodable>: Decodable {
    case one(T), many([T])

    init(from decoder: any Decoder) throws {
        let c = try decoder.singleValueContainer()
        if let many = try? c.decode([T].self) {
            self = .many(many)
        } else {
            self = .one(try c.decode(T.self))
        }
    }

    var list: [T] {
        switch self {
        case .one(let item): [item]
        case .many(let items): items
        }
    }
}

/// An object, or an empty list standing for an empty object.
private struct EmptyListOrMap<V: Decodable>: Decodable {
    let map: [String: V]

    init(from decoder: any Decoder) throws {
        let c = try decoder.singleValueContainer()
        if let empty = try? c.decode([V].self) {
            guard empty.isEmpty else {
                throw DecodingError.dataCorrupted(
                    .init(codingPath: decoder.codingPath, debugDescription: "quotas sent as a non-empty list"))
            }
            map = [:]
        } else {
            map = try c.decode([String: V].self)
        }
    }
}

extension KeyedDecodingContainer {
    func decodeStrings(forKey key: Key) throws -> [String] {
        try decode(OneOrMany<String>.self, forKey: key).list
    }

    func decodeStringsIfPresent(forKey key: Key) throws -> [String]? {
        try decodeIfPresent(OneOrMany<String>.self, forKey: key)?.list
    }

    func decodeMagnets(forKey key: Key) throws -> [Magnet] {
        try decodeIfPresent(OneOrMany<Magnet>.self, forKey: key)?.list ?? []
    }

    func decodeMagnetsIfPresent(forKey key: Key) throws -> [Magnet]? {
        try decodeIfPresent(OneOrMany<Magnet>.self, forKey: key)?.list
    }

    func decodeQuotas(forKey key: Key) throws -> [String: Int64] {
        try decode(EmptyListOrMap<Int64>.self, forKey: key).map
    }

    func decodeQuotasIfPresent(forKey key: Key) throws -> [String: Int64]? {
        try decodeIfPresent(EmptyListOrMap<Int64>.self, forKey: key)?.map
    }
}
