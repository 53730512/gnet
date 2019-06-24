module gitee.com/liyp/gnet

go 1.12

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190621222207-cc06ce4a13d4
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190621203818-d432491b9138
)

require (
	github.com/fatih/color v1.7.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gorilla/websocket v1.4.0
	github.com/mattn/go-colorable v0.1.2 // indirect
)
