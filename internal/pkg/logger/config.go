package logger

// Config defines zap logger configuration.
type Config struct {
	Level            string   `yaml:"level"`            // debug, info, warn, error
	Encoding         string   `yaml:"encoding"`         // json or console
	OutputPaths      []string `yaml:"outputPaths"`      // e.g. ["stdout"]
	ErrorOutputPaths []string `yaml:"errorOutputPaths"` // e.g. ["stderr"]
}
