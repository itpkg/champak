package shop_test

import (
	"fmt"
	"testing"

	"github.com/itpkg/champak/engines/auth"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func testAttachments(t *testing.T, db *gorm.DB) {

	var count int
	if err := db.Model(&auth.Attachment{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		for i := 0; i < 5; i++ {
			if err := db.Create(&auth.Attachment{
				Name:      fmt.Sprintf("/assets/attachments/%d.png", i),
				Title:     fmt.Sprintf("图片 %d", i),
				MediaType: "images/png",
				Summary:   fmt.Sprintf("简介 %d", i),
				UserID:    1,
			}).Error; err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestSeeds(t *testing.T) {
	db, err := auth.OpenDatabase()
	if err != nil {
		t.Fatal(err)
	}
	testAttachments(t, db)
}
