package models

import (
	"time"

	"github.com/google/uuid"
)

type Data struct {
	User            User
	messagesError   string
	messagesSuccess string
}

type Session struct {
	Userid string
	Expiry time.Time
}

type Likes struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"`
	PostID    *uuid.UUID `gorm:"type:uuid"`
	CommentID *uuid.UUID `gorm:"type:uuid"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}

type Comment struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PostID    uuid.UUID  `gorm:"type:uuid;not null"`
	AuthorID  uuid.UUID  `gorm:"type:uuid;not null"`
	ParentID  *uuid.UUID `gorm:"type:uuid"`
	Body      string     `gorm:"type:text;not null"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt *time.Time `gorm:"type:timestamptz"`
}

type Follow struct {
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null;primaryKey;check:follower_id <> following_id"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

type Bookmark struct {
	UserID    uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

type UserRole string

const (
	UserRoleMember UserRole = "member"
	UserRoleAdmin  UserRole = "admin"
)

type User struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username        string     `gorm:"type:varchar(50);uniqueIndex;not null"`
	Email           string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash    string     `gorm:"type:text;not null"`
	DisplayName     *string    `gorm:"type:varchar(100)"`
	AvatarURL       *string    `gorm:"type:text"`
	Bio             *string    `gorm:"type:text"`
	Role            UserRole   `gorm:"type:user_role;not null;default:'member'"`
	IsBanned        bool       `gorm:"type:boolean;not null;default:false"`
	EmailVerifiedAt *time.Time `gorm:"type:timestamptz"`
	CreatedAt       time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt       time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Slug        string    `gorm:"type:varchar(120);uniqueIndex;not null"`
	Description *string   `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

type Post struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AuthorID   uuid.UUID  `gorm:"type:uuid;not null"`
	CategoryID *uuid.UUID `gorm:"type:uuid"`
	Title      string     `gorm:"type:varchar(300);not null"`
	Body       string     `gorm:"type:text;not null"`
	Slug       string     `gorm:"type:varchar(350);uniqueIndex;not null"`
	IsPinned   bool       `gorm:"type:boolean;not null;default:false"`
	IsLocked   bool       `gorm:"type:boolean;not null;default:false"`
	ViewsCount int        `gorm:"type:integer;not null;default:0"`
	CreatedAt  time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt  time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt  *time.Time `gorm:"type:timestamptz"`
}
