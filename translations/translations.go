// Package translations provides for translations.
package translations

import (
	"golang.org/x/text/message/catalog"
)

type (
	// Catalog is a wrapper of catalog.Catalog.
	Catalog struct {
		catalog.Catalog
	}
	// CatalogRegistrationFunc is a function that registers translation data to catalog.Builder.
	CatalogRegistrationFunc func(c *catalog.Builder)
)

var (
	// AllPredefinedCatalogRegistrationFunc is all catalog registration functions predefined by this library.
	AllPredefinedCatalogRegistrationFunc = [...]CatalogRegistrationFunc{
		DefaultEnglish,
		DefaultJapanese,
	}
)

// NewCatalog returns a new catalog instance.
func NewCatalog(opts ...catalog.Option) *Catalog {
	return &Catalog{Catalog: catalog.NewBuilder(opts...)}
}

// Set calls CatalogRegistrationFunc, and registers translation data.
func (t *Catalog) Set(f CatalogRegistrationFunc) {
	f(t.Catalog.(*catalog.Builder))
}
