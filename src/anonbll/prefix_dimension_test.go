package anonbll

import (
	"testing"
)

func TestGetPrefixString(t *testing.T) {
	str := "Microsoft Windows 7"
	expectedPrefix := "Microsoft Windows"

	prefix := getPrefixString(str)

	if prefix != expectedPrefix {
		t.Fatalf("the value of 'prefix' should be '%v', got '%v'", expectedPrefix, prefix)
	}
}

func TestGetCommonAncestorPrefix(t *testing.T) {
	prefixGroups := map[string]int{
		"Microsoft Windows 7 Professionnel N":           6,
		"Microsoft Windows 10 Famille":                  9,
		"Microsoft Windows 7 Édition Familiale Premium": 24,
		"Microsoft Windows Server 2008 R2 Standard":     1,
		"Microsoft Windows 8 Pro":                       1,
		"Microsoft Windows 8":                           6,
		"Microsoft Windows 10 家庭中文版":                    7,
		"Microsoft Windows 10 专业版":                      3,
		"Microsoft Windows Server 2012 R2 Standard":     13,
		"Microsoft Windows Server 2008 R2 Datacenter":   9,
		"Microsoft Windows 10 企业版":                      2,
	}
	expectedAncestor := "Microsoft Windows"

	dim := prefixDimension{prefixGroups: prefixGroups}
	dim.initFullPrefixGroups()
	ancestor := dim.getCommonAncestorPrefix(81)

	if ancestor != expectedAncestor {
		t.Fatalf("Expected: '%v'; Got: '%v'", expectedAncestor, ancestor)
	}
}

func TestGetPossibleCut(t *testing.T) {
	prefixGroups := map[string]int{
		"Microsoft Windows 7 Professionnel N":           6,
		"Microsoft Windows 10 Famille":                  9,
		"Microsoft Windows 7 Édition Familiale Premium": 24,
		"Microsoft Windows Server 2008 R2 Standard":     1,
		"Microsoft Windows 8 Pro":                       1,
		"Microsoft Windows 8":                           6,
		"Microsoft Windows 10 家庭中文版":                    7,
		"Microsoft Windows 10 专业版":                      3,
		"Microsoft Windows Server 2012 R2 Standard":     13,
		"Microsoft Windows Server 2008 R2 Datacenter":   9,
		"Microsoft Windows 10 企业版":                      2,
	}
	expectedCount := 30
	expectedClosest := "Microsoft Windows 7"

	dim := prefixDimension{prefixGroups: prefixGroups}
	dim.initFullPrefixGroups()
	count, closest := dim.getPossibleCut(81)

	if count != expectedCount || closest != expectedClosest {
		t.Fatalf("Expected: '%v' - %v; Got: '%v' - %v", expectedClosest, expectedCount, closest, count)
	}
}
