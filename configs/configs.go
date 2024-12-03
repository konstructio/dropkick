package configs

const (
	DefaultVersion = "development"
)

//nolint:gochecknoglobals // used to store the version of the built binary
var Version = DefaultVersion

type Config struct {
	Version string
}
