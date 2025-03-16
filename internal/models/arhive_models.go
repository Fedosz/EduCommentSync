package models

type StudentArchive struct {
	ID       int    `gorm:"primaryKey"`
	Name     string `gorm:"size:255;not null"`
	SurName  string `gorm:"size:255;not null"`
	MailHash string `gorm:"size:255;not null"`
}

type ColabLinkArchive struct {
	ID        int    `gorm:"primaryKey"`
	ColabLink string `gorm:"size:255;not null"`
	StudentID int    `gorm:"not null"`
	WorkName  string `gorm:"not null"`
}

type RawCommentArchive struct {
	ID        int    `gorm:"primaryKey"`
	Text      string `gorm:"type:text;not null"`
	Author    string
	StudentID int    `gorm:"not null"`
	WorkName  string `gorm:"not null"`
}

type CommentArchive struct {
	ID         int    `gorm:"primaryKey"`
	StudentID  int    `gorm:"not null"`
	Text       string `gorm:"type:text;not null"`
	TaskNumber int    `gorm:"not null"`
	IsDone     bool   `gorm:"not null"`
	WorkName   string `gorm:"not null"`
}
