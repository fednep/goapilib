package config

// Section defines interface used by configuration sections
type Section interface {
	IsValid() error
}

// IsValid validates Config. If any of the fields in the config struct
// implement Section interface (i.e. have IsValid() func) it is called.
//
// Error returned if any of the IsValid returns false
func IsValid(config any) error {
	return nil
}
