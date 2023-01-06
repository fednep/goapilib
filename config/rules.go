package config

// Rules contains configuration options through which one can specify how
// configuration should be loaded
type Rules struct {

	// Toml specified whether loading from Toml will be tried
	Toml bool

	// TomlFile specifies default config name (e.g. "/etc/service/config.toml")
	// can be specified as relative (to the CWD) or absolute path
	TomlFile string

	// TomlRequired specifies whether execution should continue if config file
	// not found or error should be returned.
	TomlRequired bool

	// TomlArgsOption specifies which command line option specifies the actual
	// config file name (e.g. "-config").
	TomlArgsOption string

	// DoEnv specifies, whether loading of ".env" file will be attempted.
	// All variables from .env will be loaded as environment variables
	// and will be available for os.Getenv calls, not only to config file
	DotEnv bool

	// DotEnvFile specifies filename of the configuration file (e.g. ".env.dev)
	DotEnvFile string

	// DotEnvRequired specifies whether execution should continue if config file
	// not found or error should be returned.
	DotEnvRequired bool

	// DotEnvArgsOption specifies which command line option specifies the actual
	// config file name. (e.g. "-env").
	DotEnvArgsOption string

	// Env specifies whether configuration options should be loaded from
	// environment variables.
	Env bool

	// Validate specifies whether config validation should be performed after
	// it is loaded
	Validate bool
}

var DefaultRules = Rules{
	Toml:           true,
	TomlFile:       "config.toml",
	TomlRequired:   false,
	TomlArgsOption: "-config",

	DotEnv:           true,
	DotEnvFile:       ".env",
	DotEnvRequired:   false,
	DotEnvArgsOption: "-env",

	Env: true,

	Validate: true,
}
