package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	mihomoprovider "github.com/metacubex/mihomo/constant/provider"
	mrsprovider "github.com/metacubex/mihomo/rules/provider"
)

func TestGenerateCategorizedMRS(t *testing.T) {
	inputPath := t.TempDir()
	outputPath := t.TempDir()

	for _, category := range mrsRuleCategories {
		rule := "+.example.test\n"
		if category.behavior == mihomoprovider.IPCIDR {
			rule = "192.0.2.0/24\n"
		}
		if err := os.WriteFile(filepath.Join(inputPath, category.name+".rules"), []byte(rule), 0600); err != nil {
			t.Fatalf("writing input for %s: %v", category.name, err)
		}
	}

	if err := GenerateCategorizedMRS(inputPath, outputPath); err != nil {
		t.Fatalf("GenerateCategorizedMRS() error = %v", err)
	}

	for _, category := range mrsRuleCategories {
		path := filepath.Join(outputPath, category.name+".mrs")
		encoded, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("reading generated %s: %v", category.name, err)
			continue
		}

		var decoded bytes.Buffer
		if err := mrsprovider.ConvertToMrs(encoded, category.behavior, mihomoprovider.MrsRule, &decoded); err != nil {
			t.Errorf("decoding generated %s MRS: %v", category.name, err)
			continue
		}

		want := "+.example.test\n"
		if category.behavior == mihomoprovider.IPCIDR {
			want = "192.0.2.0/24\n"
		}
		if decoded.String() != want {
			t.Errorf("%s decoded rules = %q, want %q", category.name, decoded.String(), want)
		}
	}
}

func TestGenerateCategorizedMRSRequiresEveryCategory(t *testing.T) {
	if err := GenerateCategorizedMRS(t.TempDir(), t.TempDir()); err == nil {
		t.Fatal("GenerateCategorizedMRS() error = nil, want missing input error")
	}
}

func TestGenerateCategorizedMRSRejectsInvalidRules(t *testing.T) {
	inputPath := t.TempDir()
	outputPath := t.TempDir()

	for _, category := range mrsRuleCategories {
		rule := "+.example.test\n"
		if category.behavior == mihomoprovider.IPCIDR {
			rule = "192.0.2.0/24\n"
		}
		if err := os.WriteFile(filepath.Join(inputPath, category.name+".rules"), []byte(rule), 0600); err != nil {
			t.Fatalf("writing input for %s: %v", category.name, err)
		}
	}

	if err := os.WriteFile(filepath.Join(inputPath, "apple.rules"), []byte("invalid/domain\n"), 0600); err != nil {
		t.Fatalf("writing invalid input: %v", err)
	}

	if err := GenerateCategorizedMRS(inputPath, outputPath); err == nil {
		t.Fatal("GenerateCategorizedMRS() error = nil, want invalid rule error")
	}
}
