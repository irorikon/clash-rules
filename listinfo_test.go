package main

import (
	"bytes"
	"testing"

	mihomoprovider "github.com/metacubex/mihomo/constant/provider"
	mrsprovider "github.com/metacubex/mihomo/rules/provider"
	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

func TestListInfoMapToMRSPreservesCategories(t *testing.T) {
	listInfoMap := ListInfoMap{
		"APPLE": {GeoSite: &router.GeoSite{
			Domain: []*router.Domain{{Type: router.Domain_RootDomain, Value: "apple.example"}},
		}},
		"GOOGLE": {GeoSite: &router.GeoSite{
			Domain: []*router.Domain{{Type: router.Domain_RootDomain, Value: "google.example"}},
		}},
	}

	got, err := listInfoMap.ToMRS([]string{"apple", "google"})
	if err != nil {
		t.Fatalf("ToMRS() error = %v", err)
	}
	for _, category := range []string{"apple", "google"} {
		if _, ok := got[category+".mrs"]; !ok {
			t.Errorf("ToMRS() missing category file %q", category+".mrs")
		}
	}
	if len(got) != 2 {
		t.Errorf("ToMRS() generated %d files, want 2", len(got))
	}
}

func TestListInfoToMRS(t *testing.T) {
	list := &ListInfo{
		GeoSite: &router.GeoSite{
			Domain: []*router.Domain{
				{Type: router.Domain_Full, Value: "exact.example"},
				{Type: router.Domain_RootDomain, Value: "suffix.example"},
				{Type: router.Domain_Plain, Value: "keyword"},
				{Type: router.Domain_Regex, Value: "regexp"},
			},
		},
	}

	mrs, skipped, err := list.ToMRS()
	if err != nil {
		t.Fatalf("ToMRS() error = %v", err)
	}
	if skipped != 2 {
		t.Fatalf("ToMRS() skipped = %d, want 2", skipped)
	}

	var rules bytes.Buffer
	if err := mrsprovider.ConvertToMrs(mrs, mihomoprovider.Domain, mihomoprovider.MrsRule, &rules); err != nil {
		t.Fatalf("decoding generated MRS: %v", err)
	}
	if got, want := rules.String(), "+.suffix.example\nexact.example\n"; got != want {
		t.Fatalf("decoded MRS rules = %q, want %q", got, want)
	}
}

func TestListInfoToMRSRequiresSupportedRules(t *testing.T) {
	list := &ListInfo{GeoSite: &router.GeoSite{
		Domain: []*router.Domain{
			{Type: router.Domain_Plain, Value: "keyword"},
		},
	}}

	if _, _, err := list.ToMRS(); err == nil {
		t.Fatal("ToMRS() error = nil, want unsupported-only list error")
	}
}
