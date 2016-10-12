package main

import (
	"log"

	_ "github.com/itpkg/champak/engines/auth"
	_ "github.com/itpkg/champak/engines/forum"
	_ "github.com/itpkg/champak/engines/ops"
	_ "github.com/itpkg/champak/engines/reading"
	_ "github.com/itpkg/champak/engines/shop"
	"github.com/itpkg/champak/web"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

var version string

func main() {
	if err := web.Run(version); err != nil {
		log.Fatal(err)
	}
}
