dst=dist

build:
	go build -ldflags "-s -X main.version=`git rev-parse --short HEAD`" -o $(dst)/champak demo/main.go
	-cp -rv demo/locales demo/templates demo/db $(dst)/
	cd front && ember build --environment production
	-cp -rv front/dist $(dst)/public

clean:
	-rm -rv $(dst)
