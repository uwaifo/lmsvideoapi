package entity

import "time"

type Snippet struct {
	ID          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	UserID      uint64     `gorm:"size:100;not null;" json:"user_id"`
	User        User       `json:"user"`
	SnippetFile string     `gorm:"size:100;not null;unique" json:"snippet_file"`
	SnippetURL  string     `gorm:"text;null;" json:"snippet_url"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
