package i18n_test

import (
	"testing"

	"github.com/itpkg/champak/web/i18n"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	logging "github.com/op/go-logging"
	"golang.org/x/text/language"
)

var logger = logging.MustGetLogger("test")
var lang = language.SimplifiedChinese.String()

func TestDatabase(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	db.LogMode(true)
	i18n.Migrate(db)
	testStore(t, &i18n.GormStore{Db: db, Logger: logger})
}

func testStore(t *testing.T, p i18n.Store) {
	key := "hello"
	val := "你好"
	p.Set(lang, key, val)
	p.Set(lang, key+".1", val)
	if val1 := p.Get(lang, key); val != val1 {
		t.Errorf("want %s, get %s", val, val1)
	}
	ks, err := p.Keys(lang)
	if err != nil {
		t.Fatal(err)
	}
	if len(ks) == 0 {
		t.Errorf("empty keys")
	} else {
		t.Log(ks)
	}
	p.Del(lang, key)
}
