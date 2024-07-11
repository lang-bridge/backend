package ptesting

import (
	"golang.org/x/text/language"
	"platform/internal/translations/entity/key"
	"platform/internal/translations/entity/project"
	"platform/internal/translations/entity/translation"
	"strings"
)

const alpha = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func (g *Gen) NextString(min, max int) string {
	n := g.NextInt(min, max)
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = alpha[g.r.IntN(len(alpha))]
	}
	return string(buf)
}

func (g *Gen) NextInt(min, max int) int {
	return min + g.r.IntN(max-min+1)
}

func (g *Gen) NextKey(id project.ID) key.Key {
	return key.Key{
		ID:        0,
		ProjectID: id,
		Name:      g.NextString(5, 20),
		Platforms: []key.Platform{key.PlatformOther},
		Tags:      []key.TagID{},
	}
}

func (g *Gen) NextLanguage() language.Tag {
	return OneOf(g,
		language.English,
		language.Russian,
		language.German,
		language.Dutch,
		language.French,
		language.Spanish,
		language.Italian,
		language.Portuguese,
		language.Chinese,
		language.Japanese,
		language.Korean,
		language.Arabic,
		language.Turkish,
		language.Hindi,
		language.Bengali,
	)
}

func (g *Gen) NextTranslation(keyID key.ID) translation.Value {
	terms := Array(g, g.NextInt(1, 10), func(g *Gen) string {
		return g.NextString(2, 30)
	})
	return translation.Value{
		KeyID:       keyID,
		Translation: strings.Join(terms, " "),
		Language:    g.NextLanguage(),
	}
}

func Array[T any](g *Gen, size int, f func(*Gen) T) []T {
	var res = make([]T, size)
	for i := 0; i < size; i++ {
		res[i] = f(g)
	}
	return res
}

func Elems[T any](g *Gen, min, max int, elems ...T) []T {
	var arr = make([]T, len(elems))
	copy(arr, elems)
	g.r.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr[:g.NextInt(min, max)]
}

func OneOf[T any](g *Gen, values ...T) T {
	return values[g.r.IntN(len(values))]
}
