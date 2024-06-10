package translation

import (
	"golang.org/x/text/language"
	"platform/internal/translations/entity/key"
)

type Value struct {
	KeyID       key.ID
	Translation string
	Language    language.Tag
}
