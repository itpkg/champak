package auth

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/itpkg/champak/web"
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

//Dao db helper
type Dao struct {
	Db        *gorm.DB        `inject:""`
	Encryptor *web.Encryptor  `inject:""`
	Logger    *logging.Logger `inject:""`
}

//Set save setting
func (p *Dao) Set(k string, v interface{}, f bool) error {
	var m Setting
	null := p.Db.Where("key = ?", k).First(&m).RecordNotFound()
	if null {
		m = Setting{Key: k}
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return err
	}
	if f {
		m.Val, err = p.Encryptor.Encrypt(buf.Bytes())
		if err != nil {
			return err
		}
	} else {
		m.Val = buf.Bytes()
	}
	m.Encode = f

	if null {
		err = p.Db.Create(&m).Error
	} else {
		err = p.Db.Model(&m).Updates(map[string]interface{}{
			"encode": f,
			"val":    buf,
		}).Error
	}
	return err
}

//Get get setting value by key
func (p *Dao) Get(k string, v interface{}) error {
	var m Setting
	err := p.Db.Where("key = ?", k).First(&m).Error
	if err != nil {
		return err
	}
	if m.Encode {
		if m.Val, err = p.Encryptor.Decrypt(m.Val); err != nil {
			return err
		}
	}

	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	buf.Write(m.Val)
	return dec.Decode(v)
}

//-----------------------------------------------------------------------------

//Log add log
func (p *Dao) Log(user uint, msg string) {
	l := Log{UserID: user, Message: msg}
	if err := p.Db.Create(&l).Error; err != nil {
		p.Logger.Error(err)
	}
}

//-----------------------------------------------------------------------------

//UserClaims generate user claims
func (p *Dao) UserClaims(u *User) jws.Claims {
	cm := jws.Claims{}
	cm.SetSubject(u.Name)
	cm.Set("uid", u.UID)
	cm.Set("id", u.ID)

	cm.Set("roles", p.Authority(u.ID, "-", 0))
	return cm
}

//SignIn sign in
func (p *Dao) SignIn(email, password string) (*User, error) {
	var u User
	err := p.Db.Where("provider_id = ? AND provider_type = ?", "email", email).First(&u).Error
	if err == nil {
		if !p.Encryptor.Chk([]byte(password), u.Password) {
			err = errors.New("email and password not match")
		}
	}
	return &u, err
}

//SignUp sign up
func (p *Dao) SignUp(email, name, password string) (*User, error) {
	var u User
	var err error
	now := time.Now()
	if p.Db.Where("provider_id = ? AND provider_type = ?", "email", email).First(&u).RecordNotFound() {
		uid := uuid.NewV4().String()
		u.Email = email
		u.Name = name
		u.Home = fmt.Sprintf("%s/users/%s", viper.GetString("server.front"), uid)
		u.UID = uid
		u.ProviderID = email
		u.ProviderType = "email"
		u.SignInCount = 1
		u.LastSignIn = &now

		u.SetGravatarLogo()
		u.Password = p.Encryptor.Sum([]byte(password))

		err = p.Db.Create(&u).Error
	} else {
		err = fmt.Errorf("email %s already exists", email)
	}
	return &u, err
}

//AddOpenIDUser add openid user
func (p *Dao) AddOpenIDUser(pid, pty, email, name, home, logo string) (*User, error) {
	var u User
	var err error
	now := time.Now()
	if p.Db.Where("provider_id = ? AND provider_type = ?", pid, pty).First(&u).RecordNotFound() {
		u.Email = email
		u.Name = name
		u.Logo = logo
		u.Home = home
		u.UID = uuid.NewV4().String()
		u.ProviderID = pid
		u.ProviderType = pty
		u.ConfirmedAt = &now
		u.SignInCount = 1
		u.LastSignIn = &now
		err = p.Db.Create(&u).Error
	} else {
		err = p.Db.Model(&u).Updates(map[string]interface{}{
			"email":         email,
			"name":          name,
			"logo":          logo,
			"home":          home,
			"sign_in_count": u.SignInCount + 1,
			"last_sign_in":  &now,
		}).Error
	}
	return &u, err
}

//GetUserByUID get user by uid
func (p *Dao) GetUserByUID(uid string) (*User, error) {
	var u User
	err := p.Db.Where("uid = ?", uid).First(&u).Error
	return &u, err
}

//GetUserByEmail get user by email
func (p *Dao) GetUserByEmail(email string) (*User, error) {
	var u User
	err := p.Db.Where("provider_type = ? AND provider_id = ?", "email", email).First(&u).Error
	return &u, err
}

//Authority get user's role names
func (p *Dao) Authority(user uint, rty string, rid uint) []string {
	var items []Role
	if err := p.Db.
		Where("resource_type = ? AND resource_id = ?", rty, rid).
		Find(&items).Error; err != nil {
		p.Logger.Error(err)
	}
	var roles []string
	for _, r := range items {
		var pm Permission
		if err := p.Db.
			Where("role_id = ? AND user_id = ? ", r.ID, user).
			First(&pm).Error; err != nil {
			p.Logger.Error(err)
			continue
		}
		if pm.Enable() {
			roles = append(roles, r.Name)
		}
	}

	return roles
}

//Is is role ?
func (p *Dao) Is(user uint, name string) bool {
	return p.Can(user, name, "-", 0)
}

//Can can?
func (p *Dao) Can(user uint, name string, rty string, rid uint) bool {
	var r Role
	if p.Db.
		Where("name = ? AND resource_type = ? AND resource_id = ?", name, rty, rid).
		First(&r).
		RecordNotFound() {
		return false
	}
	var pm Permission
	if p.Db.
		Where("user_id = ? AND role_id = ?", user, r.ID).
		First(&pm).
		RecordNotFound() {
		return false
	}

	return pm.Enable()
}

//Role check role exist
func (p *Dao) Role(name string, rty string, rid uint) (*Role, error) {
	var e error
	r := Role{}
	db := p.Db
	if db.
		Where("name = ? AND resource_type = ? AND resource_id = ?", name, rty, rid).
		First(&r).
		RecordNotFound() {
		r = Role{
			Name:         name,
			ResourceType: rty,
			ResourceID:   rid,
		}
		e = db.Create(&r).Error

	}
	return &r, e
}

//Deny deny permission
func (p *Dao) Deny(role uint, user uint) error {
	return p.Db.
		Where("role_id = ? AND user_id = ?", role, user).
		Delete(Permission{}).Error
}

//Allow allow permission
func (p *Dao) Allow(role uint, user uint, years, months, days int) error {
	begin := time.Now()
	end := begin.AddDate(years, months, days)
	var count int
	p.Db.
		Model(&Permission{}).
		Where("role_id = ? AND user_id = ?", role, user).
		Count(&count)
	if count == 0 {
		return p.Db.Create(&Permission{
			UserID: user,
			RoleID: role,
			Begin:  begin,
			End:    end,
		}).Error
	}
	return p.Db.
		Model(&Permission{}).
		Where("role_id = ? AND user_id = ?", role, user).
		UpdateColumns(map[string]interface{}{"begin": begin, "end": end}).Error

}
