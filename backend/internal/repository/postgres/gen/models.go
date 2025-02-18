// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package gen

import (
	"database/sql/driver"
	"fmt"

	"golang.org/x/text/language"
)

type Platform string

const (
	PlatformWEB     Platform = "WEB"
	PlatformIOS     Platform = "IOS"
	PlatformANDROID Platform = "ANDROID"
	PlatformOTHER   Platform = "OTHER"
)

func (e *Platform) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Platform(s)
	case string:
		*e = Platform(s)
	default:
		return fmt.Errorf("unsupported scan type for Platform: %T", src)
	}
	return nil
}

type NullPlatform struct {
	Platform Platform
	Valid    bool // Valid is true if Platform is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPlatform) Scan(value interface{}) error {
	if value == nil {
		ns.Platform, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Platform.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPlatform) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Platform), nil
}

type Key struct {
	ID        int64      `db:"id"`
	ProjectID int64      `db:"project_id"`
	Name      string     `db:"name"`
	Platforms []Platform `db:"platforms"`
	Tags      []int64    `db:"tags"`
}

type KeyTag struct {
	ID        int64  `db:"id"`
	ProjectID int64  `db:"project_id"`
	Value     string `db:"value"`
}

type Project struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type Translation struct {
	ID       int64        `db:"id"`
	KeyID    int64        `db:"key_id"`
	Language language.Tag `db:"language"`
	Value    string       `db:"value"`
}
