package i18n

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	logging "github.com/op/go-logging"
	"golang.org/x/text/language"
)

//Store i18n store
type Store interface {
	Set(lang string, code, message string)
	Get(lang string, code string) string
	Del(lang string, code string)
	Keys(lang string) ([]string, error)
}

//I18n i18n helper
type I18n struct {
	Store  Store           `inject:""`
	Logger *logging.Logger `inject:""`

	Locales map[string]map[string]string
}

//Exist is lang exist?
func (p *I18n) Exist(lang string) bool {
	_, ok := p.Locales[lang]
	return ok
}

//Load load locales from filesystem
func (p *I18n) Load(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		const ext = ".ini"
		name := info.Name()
		if info.Mode().IsRegular() && filepath.Ext(name) == ext {
			p.Logger.Debugf("Find locale file %s", path)
			if err != nil {
				return err
			}
			lang := name[0 : len(name)-len(ext)]
			if _, err := language.Parse(lang); err != nil {
				return err
			}
			cfg, err := ini.Load(path)
			if err != nil {
				return err
			}

			p.Locales[lang] = cfg.Section(ini.DEFAULT_SECTION).KeysHash()
			p.Logger.Infof("find %d items", len(p.Locales[lang]))
		}
		return nil
	})
}

func (p *I18n) set(lng *language.Tag, code, message string) {
	lang := lng.String()
	if _, ok := p.Locales[lang]; !ok {
		p.Locales[lang] = make(map[string]string)
	}
	p.Locales[lang][code] = message
}

//T translate by lang tag
func (p *I18n) T(lang string, code string, args ...interface{}) string {
	msg := p.Store.Get(lang, code)
	if len(msg) == 0 {
		if items, ok := p.Locales[lang]; ok {
			msg = items[code]
		}
	}
	if len(msg) == 0 {
		msg = code
	}
	return fmt.Sprintf(msg, args...)
}
