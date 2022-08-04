package main

/*
	可以使用 go tool dist list 查看Go支持的操作系统和平台列表:
	aix/ppc64
	android/386
	android/amd64
	android/arm
	android/arm64
	darwin/amd64
	darwin/arm64
	dragonfly/amd64
	freebsd/386
	freebsd/amd64
	freebsd/arm
	..........

	Go为Gopher提供了主流编程语言中最好的跨平台交叉编译能力, 在编译时仅需
	指定目标平台的操作系统类型(GOOS)和处理器架构类型(GOARCH)即可; 但这种
	跨平台编译能力仅限于纯go代码;
	如果跨平台编译使用了cgo技术的go源文件, 会输出如下结果:
	GOOS=linux GOARCH=arm go build cgo_sleep.go
	go: no Go source files

	当Go编译器执行跨平台编译时会将CGO_ENABLED置为0, 即关闭cgo; 即找不到
	cgo_sleep.go的原因; 显式开启cgo并再来跨平台编译一下上面的cgo_sleep.go文件:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build cgo_sleep.go
	# runtime/cgo
	gcc: error: unrecognized command line option '-marm'; did you mean '-mabm'?
	即便显式开启cgo, cgo调用的linux上的外部链接器gcc也会因无法识别目标
	平台的命令而报错;


*/
