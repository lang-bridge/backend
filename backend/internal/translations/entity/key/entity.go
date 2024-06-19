package key

import (
	"context"
	"errors"
	"platform/internal/translations/entity/project"
)

type ID int64
type TagID int64

type Platform string

const (
	PlatformWeb     Platform = "WEB"
	PlatformIOS     Platform = "IOS"
	PlatformAndroid Platform = "ANDROID"
	PlatformOther   Platform = "OTHER"
)

type Tag struct {
	ID        TagID
	ProjectID project.ID
	Value     string
}

type Key struct {
	ID        ID
	ProjectID project.ID
	Name      string
	Platforms []Platform
	Tags      []TagID
}

type TagsRepository interface {
	EnsureTags(ctx context.Context, projectID project.ID, tags []string) ([]Tag, error)
}

type CreateKeyParam struct {
	ProjectID project.ID
	Name      string
	Platforms []Platform
	Tags      []TagID
}

var (
	ErrKeyAlreadyExists = errors.New("key already exists")
)

type KeysRepository interface {
	CreateKey(ctx context.Context, key CreateKeyParam) (Key, error)
}
