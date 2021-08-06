package config

type Config struct {
	App  `yaml:"app"`
	Log  `yaml:"logger"`
}

type App struct {
	Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
	Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
}

// @TODO: Find a logger.
type Log struct {
}
