package anonbll

import (
	"github.com/dargzero/anonymization/anonmodel"
	"strings"
	"unicode"
)

type prefixDimension struct {
	fieldName        string
	originalRange    int
	currentRange     int
	prefixGroups     map[string]int
	fullPrefixGroups map[string]int
}

func (p *prefixDimension) initialize(anonCollectionName string, fieldName string) {
	p.fieldName = fieldName
}

func (p *prefixDimension) getInitialBoundaries() anonmodel.Boundary {
	return &anonmodel.PrefixBoundary{Prefix: "", Filters: make(map[string]struct{})}
}

func (p *prefixDimension) getDimensionForStatistics(stat interface{}, firstRun bool) mondrianDimension {
	prefixGroups := stat.(map[string]int)
	currentRange := len(prefixGroups)

	if firstRun {
		p.originalRange = currentRange
	}
	return &prefixDimension{fieldName: p.fieldName, originalRange: p.originalRange, currentRange: currentRange, prefixGroups: prefixGroups}
}

func (p *prefixDimension) prepare(partition anonmodel.Partition, count int) {
	p.initFullPrefixGroups()

	boundary := partition[p.fieldName].(*anonmodel.PrefixBoundary)
	currentPrefix := boundary.Prefix

	commonAncestor := p.getCommonAncestorPrefix(count)
	if commonAncestor == currentPrefix {
		return
	}

	boundary.SetPrefix(commonAncestor)
}

func (p *prefixDimension) initFullPrefixGroups() {
	p.fullPrefixGroups = make(map[string]int)
	for prefix, prefixCount := range p.prefixGroups {
		for ; prefix != ""; prefix = getPrefixString(prefix) {
			if _, found := p.fullPrefixGroups[prefix]; found {
				p.fullPrefixGroups[prefix] += prefixCount
			} else {
				p.fullPrefixGroups[prefix] = prefixCount
			}
		}
	}
}

func getPrefixString(str string) string {
	index := strings.LastIndexFunc(str, unicode.IsSpace)
	if index == -1 {
		return ""
	}

	return strings.TrimRight(str[:index], " ")
}

func (p *prefixDimension) getCommonAncestorPrefix(count int) string {
	ancestor := ""
	for prefix, prefixCount := range p.fullPrefixGroups {
		if prefixCount == count && len(prefix) > len(ancestor) {
			ancestor = prefix
		}
	}
	return ancestor
}

func (p *prefixDimension) getNormalizedRange() float64 {
	// incentivize division on text fields (faster, and probably more info)
	return float64(p.currentRange)/float64(p.originalRange) + 2.0
}

func (p *prefixDimension) tryGetAllowableCut(k int, partition anonmodel.Partition, count int) (bool, anonmodel.Partition, anonmodel.Partition, error) {
	selectedCount, closestPrefix := p.getPossibleCut(count)
	if selectedCount < k || count-selectedCount < k {
		return false, nil, nil, nil
	}

	left := partition.Clone()
	leftBoundary := left[p.fieldName].(*anonmodel.PrefixBoundary)
	leftBoundary.AddFilter(closestPrefix)

	right := partition.Clone()
	rightBoundary := right[p.fieldName].(*anonmodel.PrefixBoundary)
	rightBoundary.SetPrefix(closestPrefix)

	return true, left, right, nil
}

func (p *prefixDimension) getPossibleCut(count int) (int, string) {
	half := (count + 1) / 2
	closestDiff := abs(count - half)
	closestPrefix := ""
	for prefix, prefixCount := range p.fullPrefixGroups {
		diff := abs(prefixCount - half)
		if diff < closestDiff {
			closestDiff = diff
			closestPrefix = prefix
		} else if diff == closestDiff && len(prefix) > len(closestPrefix) {
			closestPrefix = prefix
		}
	}

	selectedCount := p.fullPrefixGroups[closestPrefix]
	return selectedCount, closestPrefix
}

func abs(x int) int {
	if x >= 0 {
		return x
	}

	return -x
}
