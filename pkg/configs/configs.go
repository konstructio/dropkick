package configs

const DefaultVersion = "development"

// Version is used on version command. The value is dynamically updated on build time via ldflag. Built binary
// versions will follow semver value like v1.0.0, when not using the built version, "development" is used.
var Version = DefaultVersion
