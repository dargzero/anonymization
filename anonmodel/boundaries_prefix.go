package anonmodel

import (
	"regexp"
	"strings"
)

// PrefixBoundary represents a boundary for a prefix type field
type PrefixBoundary struct {
	Prefix  string
	Filters map[string]struct{}
}

// Clone copies a boundary
func (p *PrefixBoundary) Clone() Boundary {
	filtersClone := make(map[string]struct{})
	for key, value := range p.Filters {
		filtersClone[key] = value
	}

	return &PrefixBoundary{Prefix: p.Prefix, Filters: filtersClone}
}

// GetGeneralizedValue gets the string representation of the given bound
func (p *PrefixBoundary) GetGeneralizedValue() string {
	if p.Prefix == "" {
		return "-"
	}

	return p.Prefix
}

// SetPrefix sets the prefix of the prefix boundary
func (p *PrefixBoundary) SetPrefix(prefix string) {
	p.Prefix = prefix

	// if the prefix is not prefix for a filter, then that filter is obsolete
	for filter := range p.Filters {
		if !strings.HasPrefix(filter, p.Prefix) {
			delete(p.Filters, filter)
		}
	}
}

// AddFilter adds a new filter to the prefix boundary
func (p *PrefixBoundary) AddFilter(newFilter string) {
	// if the new filter is a prefix for a current one, that makes the current one obsolete
	regexString := "^" + regexp.QuoteMeta(newFilter) + "\\s"
	regex := regexp.MustCompile(regexString)

	for currentFilter := range p.Filters {
		if currentFilter == newFilter {
			return // the filter is already in the filters, no need to do anything
		} else if regex.MatchString(currentFilter) {
			delete(p.Filters, currentFilter)
		}
	}

	p.Filters[newFilter] = struct{}{}
}
