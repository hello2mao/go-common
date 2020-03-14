package factory

import (
	"github.com/hello2mao/go-common/csp"
	"github.com/pkg/errors"
)

const pkcs11Enabled = false

// FactoryOpts holds configuration information used to initialize factory implementations
type FactoryOpts struct {
	ProviderName string  `mapstructure:"default" json:"default" yaml:"Default"`
	SwOpts       *SwOpts `mapstructure:"SW,omitempty" json:"SW,omitempty" yaml:"SwOpts"`
}

// InitFactories must be called before using factory interfaces
// It is acceptable to call with config = nil, in which case
// some defaults will get used
// Error is returned only if defaultCSP cannot be found
func InitFactories(config *FactoryOpts) error {
	factoriesInitOnce.Do(func() {
		factoriesInitError = initFactories(config)
	})

	return factoriesInitError
}

func initFactories(config *FactoryOpts) error {
	// Take some precautions on default opts
	if config == nil {
		config = GetDefaultOpts()
	}

	if config.ProviderName == "" {
		config.ProviderName = "SW"
	}

	if config.SwOpts == nil {
		config.SwOpts = GetDefaultOpts().SwOpts
	}

	// Software-Based CSP
	if config.ProviderName == "SW" && config.SwOpts != nil {
		f := &SWFactory{}
		var err error
		defaultCSP, err = initCSP(f, config)
		if err != nil {
			return errors.Wrapf(err, "Failed initializing CSP")
		}
	}

	if defaultCSP == nil {
		return errors.Errorf("Could not find default `%s` CSP", config.ProviderName)
	}

	return nil
}

// GetCSPFromOpts returns a CSP created according to the options passed in input.
func GetCSPFromOpts(config *FactoryOpts) (csp.CSP, error) {
	var f CSPFactory
	switch config.ProviderName {
	case "SW":
		f = &SWFactory{}
	default:
		return nil, errors.Errorf("Could not find CSP, no '%s' provider", config.ProviderName)
	}

	csp, err := f.Get(config)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not initialize CSP %s", f.Name())
	}
	return csp, nil
}
