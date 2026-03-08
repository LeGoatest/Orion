package prompt

import (
	"strings"
)

// SectionDefinition provides templates for specific sections
func GetSectionHeader(s Section) string {
	return "### " + string(s) + " ###"
}

func Deduplicate(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
