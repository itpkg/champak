# champak
A web application build by go-lang and reactjs.

## Install go

### For ubuntu

```
add-apt-repository ppa:ubuntu-lxc/lxd-stable
sudo apt-get update
sudo apt-get install golang
```

### For archlinux

```
sudo pacman -S go go-tools
```

### Add to your .bashrc or .zshrc

```
GOPATH=$HOME/go
PATH=$GOPATH/bin:$PATH
export GOPATH PATH
```

### For development

```
go get -u github.com/nsf/gocode
go get -u github.com/derekparker/delve/cmd/dlv
go get -u github.com/alecthomas/gometalinter
go get -u github.com/golang/lint/golint

go get -u github.com/kardianos/govendor

go get -u bitbucket.org/liamstask/goose/cmd/goose
go get -u github.com/itpkg/champak
```

### Start
```
cd $GOPATH/src/github.com/itpkg/champak
govendor sync
cd demo && go run main.go   # backend server
cd front-react && npm run dev # frontend server
```

## Database creation

### postgresql

```
psql -U postgres
CREATE DATABASE db-name WITH ENCODING = 'UTF8';
CREATE USER user-name WITH PASSWORD 'change-me';
GRANT ALL PRIVILEGES ON DATABASE db-name TO user-name;
```

* ExecStartPre=/usr/bin/postgresql-check-db-dir ${PGROOT}/data (code=exited, status=1/FAILURE)

```
initdb  -D '/var/lib/postgres/data'
```

## Build

```
cd $GOPATH/src/github.com/itpkg/champak
make
ls -lh dist
```

## Documents
- [react](https://facebook.github.io/react/docs/getting-started.html)
- [react-bootstrap](http://react-bootstrap.github.io/)
- [gin](https://github.com/gin-gonic/gin)
- [goose](https://bitbucket.org/liamstask/goose/)
- [go-plus](https://atom.io/packages/go-plus)
- [gorm](http://jinzhu.me/gorm/)
- [locale](https://blog.golang.org/matchlang)
- [govendor](https://github.com/kardianos/govendor)
