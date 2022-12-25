package config

// ConfigSection defines interface used by configuration sections
type ConfigSection interface {
	IsValid() error
	LoadEnvOverrides()
}

// LoadOverride load config parameters override from .env file
// and validates the configuration options
func LoadOverride(configs []ConfigSection) error {

	for _, section := range configs {
		section.LoadEnvOverrides()
		if err := section.IsValid(); err != nil {
			return err
		}
	}
	return nil
}
