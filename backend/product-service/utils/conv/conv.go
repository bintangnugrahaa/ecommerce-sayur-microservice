package conv

import "strings"

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}
