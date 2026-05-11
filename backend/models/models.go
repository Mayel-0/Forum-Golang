package models

import (
	"time"

	"github.com/google/uuid"
)

type PostSubCategory string

const (
	SubCategoryConcerts   PostSubCategory = "concerts"
	SubCategoryArtistes   PostSubCategory = "artistes"
	SubCategoryNouveautes PostSubCategory = "nouveautes"
)

type Data struct {
	User     User
	Post     Post
	Category Category

	Categories []Category
	TopUsers   []User
	Comments   []Comment

	Posts           []Post
	PostsArtists    []Post
	PostsConcerts   []Post
	PostsNouveautes []Post

	TotalArtistes    int64
	TotalConcerts    int64
	TotalNouveautes  int64
	SubCategoryLabel string

	RecentPosts  []Post
	PopularPosts []Post

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
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PostID        uuid.UUID  `gorm:"type:uuid;not null"`
	AuthorID      uuid.UUID  `gorm:"type:uuid;not null"`
	ParentID      *uuid.UUID `gorm:"type:uuid"`
	Body          string     `gorm:"type:text;not null"`
	LikesCount    int        `gorm:"type:integer;not null;default:0"`
	CommentsCount int        `gorm:"type:integer;not null;default:0"`
	CreatedAt     time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt     time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt     *time.Time `gorm:"type:timestamptz"`
	Author        User       `gorm:"foreignKey:AuthorID"`
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

	// Champs calculés — pas en DB
	LikesCount    int `gorm:"->"`
	PostsCount    int `gorm:"->"`
	CommentsCount int `gorm:"->"`
}

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Slug        string    `gorm:"type:varchar(120);uniqueIndex;not null"`
	Description *string   `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

type Post struct {
	ID            uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AuthorID      uuid.UUID       `gorm:"type:uuid;not null"`
	CategoryID    *uuid.UUID      `gorm:"type:uuid"`
	SubCategory   PostSubCategory `gorm:"type:post_sub_category;not null"`
	Title         string          `gorm:"type:varchar(300);uniqueIndex;not null"`
	Body          string          `gorm:"type:text;not null"`
	Slug          string          `gorm:"type:varchar(350);uniqueIndex;not null"`
	IsPinned      bool            `gorm:"type:boolean;not null;default:false"`
	IsLocked      bool            `gorm:"type:boolean;not null;default:false"`
	LikesCount    int             `gorm:"type:integer;not null;default:0"`
	CommentsCount int             `gorm:"type:integer;not null;default:0"`
	CreatedAt     time.Time       `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt     time.Time       `gorm:"type:timestamptz;not null;default:now()"`
	DeletedAt     *time.Time      `gorm:"type:timestamptz"`

	Author User `gorm:"foreignKey:AuthorID"`
}
