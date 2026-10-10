import org.jetbrains.kotlin.gradle.dsl.JvmTarget

plugins {
    kotlin("jvm") version "2.4.20"
    kotlin("plugin.serialization") version "2.4.20"
    `java-library`
    `maven-publish`
}

group = "io.github.pierre-emmanuelj"
version = "1.0.0"

repositories {
    mavenCentral()
}

dependencies {
    api("com.squareup.okhttp3:okhttp:5.5.0")
    api("org.jetbrains.kotlinx:kotlinx-serialization-json:1.11.0")
    api("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.11.0")

    testImplementation(kotlin("test"))
    testImplementation("org.junit.jupiter:junit-jupiter:5.14.4")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
    testImplementation("com.squareup.okhttp3:mockwebserver3:5.5.0")
    testImplementation("org.jetbrains.kotlinx:kotlinx-coroutines-test:1.11.0")
    testImplementation("org.yaml:snakeyaml:2.7")
}

kotlin {
    explicitApi()
    jvmToolchain(21)
    compilerOptions {
        // Android and older JVMs.
        jvmTarget.set(JvmTarget.JVM_11)
        allWarningsAsErrors.set(true)
    }
}

java {
    targetCompatibility = JavaVersion.VERSION_11
    withSourcesJar()
}

tasks.test {
    useJUnitPlatform()
    environment("OPENAPI_SPEC", rootProject.file("../../openapi.yaml").absolutePath)
    testLogging {
        events("failed")
        exceptionFormat = org.gradle.api.tasks.testing.logging.TestExceptionFormat.FULL
    }
}

// Regenerates Models.kt from the spec (needs python3 with PyYAML).
tasks.register<Exec>("generateModels") {
    workingDir = rootProject.file("../..")
    commandLine(
        "python3", "scripts/gen_models.py", "kotlin",
        "sdks/kotlin/src/main/kotlin/io/github/pierreemmanuelj/alldebrid/Models.kt",
    )
}

publishing {
    publications {
        create<MavenPublication>("maven") {
            from(components["java"])
            pom {
                name.set("open-alldebrid")
                description.set("Kotlin client for the AllDebrid API, built from the open-alldebrid OpenAPI spec")
                url.set("https://github.com/pierre-emmanuelJ/open-alldebrid")
                licenses {
                    license {
                        name.set("MIT")
                        url.set("https://opensource.org/license/mit")
                    }
                }
            }
        }
    }
}
