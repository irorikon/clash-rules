package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	mihomoprovider "github.com/metacubex/mihomo/constant/provider"
	mrsprovider "github.com/metacubex/mihomo/rules/provider"
)

var mrsRuleCategories = []struct {
	name     string
	behavior mihomoprovider.RuleBehavior
}{
	{name: "apple", behavior: mihomoprovider.Domain},
	{name: "cncidr", behavior: mihomoprovider.IPCIDR},
	{name: "direct", behavior: mihomoprovider.Domain},
	{name: "gfw", behavior: mihomoprovider.Domain},
	{name: "google", behavior: mihomoprovider.Domain},
	{name: "greatfire", behavior: mihomoprovider.Domain},
	{name: "icloud", behavior: mihomoprovider.Domain},
	{name: "lancidr", behavior: mihomoprovider.IPCIDR},
	{name: "private", behavior: mihomoprovider.Domain},
	{name: "proxy", behavior: mihomoprovider.Domain},
	{name: "reject", behavior: mihomoprovider.Domain},
	{name: "telegramcidr", behavior: mihomoprovider.IPCIDR},
	{name: "tld-not-cn", behavior: mihomoprovider.Domain},
}

// GenerateCategorizedMRS converts normalized category rule files to MRS.
// Each input file contains one domain or CIDR rule per line and is named
// <category>.rules.
func GenerateCategorizedMRS(inputPath, outputPath string) error {
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	for _, category := range mrsRuleCategories {
		inputFile := filepath.Join(inputPath, category.name+".rules")
		rules, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("%s: %w", category.name, err)
		}
		if len(bytes.TrimSpace(rules)) == 0 {
			return fmt.Errorf("%s: no rules to convert", category.name)
		}

		var mrs bytes.Buffer
		if err := mrsprovider.ConvertToMrs(rules, category.behavior, mihomoprovider.TextRule, &mrs); err != nil {
			return fmt.Errorf("%s: %w", category.name, err)
		}

		var decodedRules bytes.Buffer
		if err := mrsprovider.ConvertToMrs(mrs.Bytes(), category.behavior, mihomoprovider.MrsRule, &decodedRules); err != nil {
			return fmt.Errorf("%s: generated MRS failed validation: %w", category.name, err)
		}
		if len(bytes.TrimSpace(decodedRules.Bytes())) == 0 {
			return fmt.Errorf("%s: generated MRS contains no readable rules", category.name)
		}

		if err := os.WriteFile(filepath.Join(outputPath, category.name+".mrs"), mrs.Bytes(), 0644); err != nil {
			return fmt.Errorf("%s: %w", category.name, err)
		}
	}

	return nil
}
