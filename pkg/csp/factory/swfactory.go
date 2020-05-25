package factory

import (
	"github.com/hello2mao/go-common/pkg/csp"
	"github.com/hello2mao/go-common/pkg/csp/sw"
	"github.com/pkg/errors"
)

const (
	// SoftwareBasedFactoryName is the name of the factory of the software-based CSP implementation
	SoftwareBasedFactoryName = "SW"
)

// SWFactory is the factory of the software-based CSP.
type SWFactory struct{}

// Name returns the name of this factory
func (f *SWFactory) Name() string {
	return SoftwareBasedFactoryName
}

// Get returns an instance of CSP using Opts.
func (f *SWFactory) Get(config *FactoryOpts) (csp.CSP, error) {
	// Validate arguments
	if config == nil || config.SwOpts == nil {
		return nil, errors.New("Invalid config. It must not be nil.")
	}

	swOpts := config.SwOpts

	var ks csp.KeyStore
	switch {
	case swOpts.Ephemeral:
		ks = sw.NewDummyKeyStore()
	case swOpts.FileKeystore != nil:
		fks, err := sw.NewFileBasedKeyStore(nil, swOpts.FileKeystore.KeyStorePath, false)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to initialize software key store")
		}
		ks = fks
	case swOpts.InmemKeystore != nil:
		ks = sw.NewInMemoryKeyStore()
	default:
		// Default to ephemeral key store
		ks = sw.NewDummyKeyStore()
	}

	return sw.NewWithParams(swOpts.SecLevel, swOpts.HashFamily, ks)
}

// SwOpts contains options for the SWFactory
type SwOpts struct {
	// Default algorithms when not specified (Deprecated?)
	SecLevel   int    `mapstructure:"security" json:"security" yaml:"Security"`
	HashFamily string `mapstructure:"hash" json:"hash" yaml:"Hash"`

	// Keystore Options
	Ephemeral     bool               `mapstructure:"tempkeys,omitempty" json:"tempkeys,omitempty"`
	FileKeystore  *FileKeystoreOpts  `mapstructure:"filekeystore,omitempty" json:"filekeystore,omitempty" yaml:"FileKeyStore"`
	DummyKeystore *DummyKeystoreOpts `mapstructure:"dummykeystore,omitempty" json:"dummykeystore,omitempty"`
	InmemKeystore *InmemKeystoreOpts `mapstructure:"inmemkeystore,omitempty" json:"inmemkeystore,omitempty"`
}

// Pluggable Keystores, could add JKS, P12, etc..
type FileKeystoreOpts struct {
	KeyStorePath string `mapstructure:"keystore" yaml:"KeyStore"`
}

type DummyKeystoreOpts struct{}

// InmemKeystoreOpts - empty, as there is no config for the in-memory keystore
type InmemKeystoreOpts struct{}
