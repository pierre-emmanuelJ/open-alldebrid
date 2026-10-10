// swift-tools-version: 6.0
// Lets apps add the Swift SDK with the repository URL. The SDK and its tests
// live in sdks/swift.
import PackageDescription

let package = Package(
    name: "OpenAllDebrid",
    platforms: [.macOS(.v13), .iOS(.v16), .tvOS(.v16), .watchOS(.v9)],
    products: [
        .library(name: "OpenAllDebrid", targets: ["OpenAllDebrid"])
    ],
    targets: [
        .target(name: "OpenAllDebrid", path: "sdks/swift/Sources/OpenAllDebrid")
    ]
)
