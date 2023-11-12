package post

import "time"

const ()

type Post struct {
	ID         uint
	BlogID     uint
	Sysname    string
	KeywordIDs *[]uint
	TagIDs     *[]uint
	IsDeleted  bool
	Title      string
	Preview    string
	Content    string
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
}

func (e *Post) Validate() error {
	return nil
}

type PostPreview struct {
	ID         uint
	BlogID     uint
	Sysname    string
	KeywordIDs *[]uint
	TagIDs     *[]uint
	IsDeleted  bool
	Title      string
	Preview    string
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
}
