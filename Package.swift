// swift-tools-version: 6.0
// Lets apps add the Swift SDK with the repository URL. The SDK and its tests
// live in sdks/swift.
import PackageDescription

let package = Package(
    name: "OpenAllDebrid",
    platforms: [.macOS(.v12), .iOS(.v15), .tvOS(.v15), .watchOS(.v8)],
    products: [
        .library(name: "OpenAllDebrid", targets: ["OpenAllDebrid"])
    ],
    targets: [
        .target(name: "OpenAllDebrid", path: "sdks/swift/Sources/OpenAllDebrid")
    ]
)
