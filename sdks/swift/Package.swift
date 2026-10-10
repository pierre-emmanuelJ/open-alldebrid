// swift-tools-version: 6.0
import PackageDescription

let package = Package(
    name: "OpenAllDebrid",
    platforms: [.macOS(.v13), .iOS(.v16), .tvOS(.v16), .watchOS(.v9)],
    products: [
        .library(name: "OpenAllDebrid", targets: ["OpenAllDebrid"])
    ],
    targets: [
        .target(name: "OpenAllDebrid"),
        .testTarget(name: "OpenAllDebridTests", dependencies: ["OpenAllDebrid"]),
    ]
)
