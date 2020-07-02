package entity

import (
"html"
"strings"
"time"
)

// Food . . .
type Video struct {
	ID          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	UserID      uint64     `gorm:"size:100;not null;" json:"user_id"`
	User      User       `json:"user"`
	VideoName string     `gorm:"size:100;not null;unique" json:"video_name"`
	VideoURL  string     `gorm:"text;null;" json:"video_url"`
	Title       string     `gorm:"size:100;not null;unique" json:"title"`
	Description string     `gorm:"text;not null;" json:"description"`
	VideoImage   string     `gorm:"size:255;null;" json:"video_image"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}



// BeforeSave . .
func (f *Video) BeforeSave() {
	f.Title = html.EscapeString(strings.TrimSpace(f.Title))
}

// Prepare . .
func (f *Video) Prepare() {
	f.Title = html.EscapeString(strings.TrimSpace(f.Title))
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
}

// Validate . .
func (f *Video) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if f.Title == "" || f.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
		if f.Description == "" || f.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
	default:
		if f.Title == "" || f.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
		if f.Description == "" || f.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
	}
	return errorMessages
}
