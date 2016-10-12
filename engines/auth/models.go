package auth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/itpkg/champak/web"
)

//LeaveWord leave word
type LeaveWord struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

//Notice notice
type Notice struct {
	web.Model
	Lang    string `json:"lang"`
	Content string `json:"content"`
}

//Setting setting
type Setting struct {
	web.Model

	Key    string
	Val    []byte
	Encode bool
}

//User user model
type User struct {
	web.Model
	Email    string `json:"email"`
	UID      string `json:"uid"`
	Home     string `json:"home"`
	Logo     string `json:"logo"`
	Name     string `json:"name"`
	Password []byte `json:"-"`

	ProviderType string `json:"provider_type"`
	ProviderID   string `json:"provider_id"`

	LastSignIn  *time.Time `json:"last_sign_in"`
	SignInCount uint       `json:"sign_in_count"`
	ConfirmedAt *time.Time `json:"confirmed_at"`
	LockedAt    *time.Time `json:"locked_at"`

	Permissions []Permission `json:"permissions"`
	Logs        []Log        `json:"logs"`
}

//IsConfirmed confirmed?
func (p *User) IsConfirmed() bool {
	return p.ConfirmedAt != nil
}

//IsLocked locked?
func (p *User) IsLocked() bool {
	return p.LockedAt != nil
}

//IsAvailable is valid?
func (p *User) IsAvailable() bool {
	return p.IsConfirmed() && !p.IsLocked()
}

//SetGravatarLogo set logo by gravatar
func (p *User) SetGravatarLogo() {
	buf := md5.Sum([]byte(strings.ToLower(p.Email)))
	p.Logo = fmt.Sprintf("https://gravatar.com/avatar/%s.png", hex.EncodeToString(buf[:]))
}

func (p User) String() string {
	return fmt.Sprintf("%s<%s>", p.Name, p.Email)
}

//Log model
type Log struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"-"`
	User      User      `json:"-"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

//Role role model
type Role struct {
	web.Model

	Name         string
	ResourceType string
	ResourceID   uint
}

func (p Role) String() string {
	return fmt.Sprintf("%s@%s://%d", p.Name, p.ResourceType, p.ResourceID)
}

//Permission permission model
type Permission struct {
	web.Model
	User   User
	UserID uint
	Role   Role
	RoleID uint
	Begin  time.Time
	End    time.Time
}

//EndS end to string
func (p *Permission) EndS() string {
	return p.End.Format("2006-01-02")
}

//BeginS begin to string
func (p *Permission) BeginS() string {
	return p.Begin.Format("2006-01-02")
}

//Enable is enable?
func (p *Permission) Enable() bool {
	now := time.Now()
	return now.After(p.Begin) && now.Before(p.End)
}

//Attachment attachment
type Attachment struct {
	web.Model

	Title     string
	Name      string
	MediaType string
	Summary   string

	UserID uint
	User   User
}
