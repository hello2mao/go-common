package factory

import (
	"sync"

	"github.com/hello2mao/go-common/incubator/csp"
	"github.com/pkg/errors"
)

var (
	defaultCSP         csp.CSP   // default CSP
	factoriesInitOnce  sync.Once // factories' Sync on Initialization
	factoriesInitError error     // Factories' Initialization Error

	// when InitFactories has not been called yet (should only happen
	// in test cases), use this CSP temporarily
	bootCSP         csp.CSP
	bootCSPInitOnce sync.Once
)

// CSPFactory is used to get instances of the CSP interface.
// A Factory has name used to address it.
type CSPFactory interface {

	// Name returns the name of this factory
	Name() string

	// Get returns an instance of CSP using opts.
	Get(opts *FactoryOpts) (csp.CSP, error)
}

// GetDefault returns a non-ephemeral (long-term) CSP
func GetDefault() csp.CSP {
	if defaultCSP == nil {
		bootCSPInitOnce.Do(func() {
			var err error
			bootCSP, err = (&SWFactory{}).Get(GetDefaultOpts())
			if err != nil {
				panic("CSP Internal error, failed initialization with GetDefaultOpts!")
			}
		})
		return bootCSP
	}
	return defaultCSP
}

func initCSP(f CSPFactory, config *FactoryOpts) (csp.CSP, error) {
	cspInstance, err := f.Get(config)
	if err != nil {
		return nil, errors.Errorf("Could not initialize CSP %s [%s]", f.Name(), err)
	}

	return cspInstance, nil
}
