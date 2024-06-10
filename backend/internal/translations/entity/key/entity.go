package key

import "platform/internal/translations/entity/project"

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
	Tags      []Tag
}
