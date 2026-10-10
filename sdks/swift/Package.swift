// swift-tools-version: 6.0
import PackageDescription

let package = Package(
    name: "OpenAllDebrid",
    platforms: [.macOS(.v12), .iOS(.v15), .tvOS(.v15), .watchOS(.v8)],
    products: [
        .library(name: "OpenAllDebrid", targets: ["OpenAllDebrid"])
    ],
    targets: [
        .target(name: "OpenAllDebrid"),
        .testTarget(name: "OpenAllDebridTests", dependencies: ["OpenAllDebrid"]),
    ]
)
