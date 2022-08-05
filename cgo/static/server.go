package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.FileServer(http.Dir(cwd)),
	}
	log.Fatal(srv.ListenAndServe())
}

/*
	使用cgo代码的静态构建
	TODO: 熟悉编译, 链接, go编译器等
	<go语言精进之路2>, 60.6 节, 此节内容基于go1.14, 后续版本相关内容有所变化

	go build server.go
	ldd server

	linux-vdso.so.1 (0x00007fff4d72e000)
	libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007fd4e1fc1000)
	libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd4e1dcf000)
	/lib64/ld-linux-x86-64.so.2 (0x00007fd4e2002000)
	默认构建出的Go应用有多个对外部动态共享库的依赖

	默认情况下, Go的运行时环境变量CGO_ENABLED=1(通过go env命令可以查看),
	即默认开启cgo, 允许在Go代码中调用C代码; Go的预编译标准库的.a文件也是
	在这种情况下编译出来的; 在$GOROOT/pkg/linux_amd64中, 遍历所有预编译好
	的标准库.a文件, 并用nm输出每个.a文件中的未定义符号(状态为U) TODO: nm 命令
	cd $GOROOT/pkg/linux_amd64
	nm -uA net.a
	nm: __.PKGDEF: file format not recognized
	nm: _go_.o: file format not recognized
	net.a:_x003.o:                 U _cgo_topofstack
	net.a:_x003.o:                 U __errno_location
	net.a:_x003.o:                 U getnameinfo
	net.a:_x003.o:                 U _GLOBAL_OFFSET_TABLE_
	net.a:_x005.o:                 U _cgo_topofstack
	net.a:_x005.o:                 U __errno_location
	net.a:_x005.o:                 U freeaddrinfo
	net.a:_x005.o:                 U gai_strerror
	net.a:_x005.o:                 U getaddrinfo
	net.a:_x005.o:                 U _GLOBAL_OFFSET_TABLE_

	nm -uA os/user.a
	nm: __.PKGDEF: file format not recognized
	nm: _go_.o: file format not recognized
	user.a:_x001.o:                 U _GLOBAL_OFFSET_TABLE_
	user.a:_x001.o:                 U malloc
	user.a:_x003.o:                 U _cgo_topofstack
	user.a:_x003.o:                 U free
	user.a:_x003.o:                 U getgrgid_r
	user.a:_x003.o:                 U getgrnam_r
	user.a:_x003.o:                 U getpwnam_r
	user.a:_x003.o:                 U getpwuid_r
	user.a:_x003.o:                 U _GLOBAL_OFFSET_TABLE_
	user.a:_x003.o:                 U realloc
	user.a:_x003.o:                 U sysconf
	user.a:_x004.o:                 U _cgo_topofstack
	user.a:_x004.o:                 U getgrouplist
	user.a:_x004.o:                 U _GLOBAL_OFFSET_TABLE_
	可以看到上面的包依赖的外部链接

	以os/user为例, 在CGO_ENABLED=1 即cgo开启的情况下, os/user包中的
	lookup-UserXXX系列函数采用了cgo版本的实现;
	在$GOROOT/src/os/user/cgo_lookup_unix.go源文件中的build tag中包含了
	+build cgo的构建指示器; CGO_ENABLED=1的情况下该文件才会被编译,
	该文件中的cgo版本实现的lookupUser将被使用:

	//go:build (aix || darwin || dragonfly || freebsd || (!android && linux) || netbsd || openbsd || solaris) && cgo && !osusergo

	package user
	......


	func lookupUser(username string) (*User, error) {
		var pwd C.struct_passwd
		var result *C.struct_passwd
		nameC := make([]byte, len(username)+1)
		copy(nameC, username)

		buf := alloc(userBuffer)
		defer buf.free()

		err := retryWithBuffer(buf, func() syscall.Errno {
			// mygetpwnam_r is a wrapper around getpwnam_r to avoid
			// passing a size_t to getpwnam_r, because for unknown
			// reasons passing a size_t to getpwnam_r doesn't work on
			// Solaris.
			return syscall.Errno(C.mygetpwnam_r((*C.char)(unsafe.Pointer(&nameC[0])),
				&pwd,
				(*C.char)(buf.ptr),
				C.size_t(buf.size),
				&result))
		})
		if err != nil {
			return nil, fmt.Errorf("user: lookup username %s: %v", username, err)
		}
		if result == nil {
			return nil, UnknownUserError(username)
		}
		return buildUser(&pwd), err
	}

	凡是依赖上述包的Go代码最终编译的可执行文件都要有外部依赖, 即默认情况下
	编译出的server有外部依赖的原因(server至少依赖net.a)

	以上cgo版本实现都有对应的go版本实现,  os/user 包的  lookupUser 函数的
	go版本如下: $GOROOT/src/os/user/lookup.go

	// Lookup looks up a user by username. If the user cannot be found, the
	// returned error is of type UnknownUserError.
	func Lookup(username string) (*User, error) {
		if u, err := Current(); err == nil && u.Username == username {
			return u, err
		}
		return lookupUser(username)
	}

	可以通过设置CGO_ENABLED=0关闭cgo是促使编译器选用Go版本实现的前提条件:
	CGO_ENABLED=0 go build -o server server.go
	ldd server
		not a dynamic executable
	nm -u A server
	关闭cgo后编译得到的server是一个静态编译的程序, 没有对外部的任何依赖;
	如果使用go build的-x -v选项可以看到Go编译器会重新编译依赖的包的静态版本
	(包括net等), 并将编译后的.a(以包为单位)放入编译器构建缓存目录下
	(比如~/.cache/go-build/xxx, 后续复用), 然后再静态链接这些版本;


	即使在CGO_ENABLED=1默认值的情况下, 也可以实现纯静态链接
	即告诉链接器在最后的链接时采用静态链接方式, 哪怕依赖的Go标准库中某些包
	使用的是C版本的实现;
	根据Go官方文档($GOROOT/cmd/cgo/doc.go), Go链接器支持两种工作模式:
	内部链接(internal linking)和外部链接(external linking)

	TODO: 理解
	如果用户代码中仅仅使用了net、os/user等几个标准库中的依赖cgo的包,
	Go链接器默认使用内部链接, 而无须启动外部链接器(如gcc、clang等);
	不过Go链接器功能有限, 仅仅将.o和预编译好的标准库的.a写到最终二进
	制文件中; 因此如果标准库中是在CGO_ENABLED=1的情况下编译的, 那么编译
	出来的最终二进制文件依旧是动态链接的, 即便在go build时
	传入-ldflags 'extldflags "-static"'也是如此, 因为根本没有用到外部链接器;

	go build -o server -ldflags '-extldflags "-static"' server.go
	ldd server
	    linux-vdso.so.1 (0x00007ffd14ec7000)
        libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007fb27b454000)
        libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fb27b262000)
        /lib64/ld-linux-x86-64.so.2 (0x00007fb27b495000)
	可以看到依旧依赖了动态链接

	TODO: 因为go版本变化, 以下的描述会有所偏差
	而外部链接机制则是Go链接器将所有生成的.o都写到一个.o文件中, 再将其交
	给外部链接器(比如gcc或clang)去做最终的链接处理; 如果此时在go build
	的命令行参数中传入-ldflags ‘extldflags “-static”’, 那么gcc/clang将会
	做静态链接, 将.o中未定义(undefined)的符号都替换为真正的代码指令;
	可以通过-linkmode=external来强制Go链接器采用外部链接; TODO: 实践
	go build -o server -ldflags '-linkmode "external" -extldflags  "-static"' server.go
	按照书籍说明执行上述命令, 编译失败

*/
