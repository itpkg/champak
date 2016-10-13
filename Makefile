dst=dist

build:
	go build -ldflags "-s -X main.version=`git rev-parse --short HEAD`" -o $(dst)/champak demo/main.go
	-cp -rv demo/locales demo/templates demo/db $(dst)/
	cd front-react && npm run build
	-cp -rv front-react/dist $(dst)/public

clean:
	-rm -rv $(dst)
