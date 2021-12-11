package main

import "fmt"

/* 机器码生成---链接*/

/*
 * 程序可能使用其他程序或程序库(library), 编写的程序必须与这些程序或程序库组合
 * 在一起才能执行, 链接就是将编写的程序与外部程序组合在一起的过程; 链接分为:
 * - 静态链接: 链接器会将程序中使用的所有库程序复制到最后的可质询文件中; 因此
 *     静态链接更快, 并且可移植, 不需要在运行它的系统上存在该库, 但是会占用
 *     更多的磁盘和内存空间; 静态链接发生在编译时的最后一步.
 * - 动态链接: 只会在最后的可执行文件中存储动态链接库的位置, 并在运行时调用,
 *     动态链接发生在程序加载到内存时.
 *
 */

// go 默认使用静态链接, 但在一些特殊情况下, 如在使用了 CGO(即引用了C代码) 时,
// 则会使用操作系统的动态链接库; 例如, go 语言的 net/http 包在默认情况下会使
// 用 libpthread 与 lib c 的动态链接库; go 语言也支持在 go build 编译时通过
// 传递参数来指定要生成的链接库的方式, 可以使用 go help build 命令查看.

// 使用 hello, world 程序为例说明编译与链接的过程
func main() {
	fmt.Println("hello, world")
}

// go build -x link.go
// 输出结果分段说明(TODO)

// 创建一个临时目录, 用于存放临时文件; 在默认情况下, 命令结束时会自动删除此
// 目录, 如果需要保留则添加 -work 参数.
//WORK=/tmp/go-build574291369
//mkdir -p $WORK/b001/
//cat >$WORK/b001/_gomod_.go << 'EOF' # internal
//package main
//import _ "unsafe"
////go:linkname __debug_modinfo__ runtime.modinfo
//var __debug_modinfo__ = "0w\xaf\f\x92t\b\x02A\xe1\xc1\a\xe6\xd6\x18\xe6path\tcommand-line-arguments\nmod\tcommand-line-arguments\t(devel)\t\n\xf92C1\x86\x18 r\x00\x82B\x10A\x16\xd8\xf2"
//        EOF
//cat >$WORK/b001/importcfg << 'EOF' # internal

//// 生成编译配置文件, 主要为编译过程需要的外部依赖(如引用其他包的函数定义)
//# import config
//packagefile fmt=/usr/local/go/pkg/linux_amd64/fmt.a
//packagefile runtime=/usr/local/go/pkg/linux_amd64/runtime.a
//EOF

// 编译阶段会生成中间结果 $WORK/b001/_pkg_.a
//cd /home/lucas/gogo/gotta/compiler
///usr/local/go/pkg/tool/linux_amd64/compile -o $WORK/b001/_pkg_.a -trimpath "$WORK/b001=>" -p main -complete -buildid m15EitMB9lGrgk-sJIEh/m15EitMB9lGrgk-sJIEh -goversion go1.15.6 -D _/home/lucas/gogo/gotta/compiler -importcfg $WORK/b001/importcfg -pack -c=4 ./link.go $WORK/b001/_gomod_.go
///usr/local/go/pkg/tool/linux_amd64/buildid -w $WORK/b001/_pkg_.a # internal
//cp $WORK/b001/_pkg_.a /home/lucas/.cache/go-build/fb/fb3610e50807ab9192774644f16e2a81963751f9bfaee90b6648980ad112dc25-d # internal

/*
 * .a 类型的文件叫目标文件(object file), 是一个压缩包, 其内部包含 __.PKGDEF 和
 * _go_.o 两个文件, 分别为编译目标文件和链接目标文件.
 * 可进入 /tmp/go-build192124789/b001 中查看:
 * $ file _pkg_.a  # 检查文件格式
 * _pkg_.a: current ar archive # 说明是 ar 格式的打包文件
 * $ ar x _pkg_.a # 解包文件
 * $ ls
 * __.PKGDEF  _go_.o
 *
 * 这两个文件的文件内容由导出的函数, 变量及引用的其他包的信息组成; 了解包含的
 * 信息需要查看 go 编译器实现的代码, 核心逻辑在:
 * $GOROOT/src/cmd/compile/internal/gc/obj.go 中; (./ar.go) TODO: 源码
 */

// 生成配置文件, 主要包含了需要链接的依赖项
//cat >$WORK/b001/importcfg.link << 'EOF' # internal
//packagefile command-line-arguments=$WORK/b001/_pkg_.a
//packagefile fmt=/usr/local/go/pkg/linux_amd64/fmt.a
//packagefile runtime=/usr/local/go/pkg/linux_amd64/runtime.a
//packagefile errors=/usr/local/go/pkg/linux_amd64/errors.a
//packagefile internal/fmtsort=/usr/local/go/pkg/linux_amd64/internal/fmtsort.a
//packagefile io=/usr/local/go/pkg/linux_amd64/io.a
//packagefile math=/usr/local/go/pkg/linux_amd64/math.a
//packagefile os=/usr/local/go/pkg/linux_amd64/os.a
//packagefile reflect=/usr/local/go/pkg/linux_amd64/reflect.a
//packagefile strconv=/usr/local/go/pkg/linux_amd64/strconv.a
//packagefile sync=/usr/local/go/pkg/linux_amd64/sync.a
//packagefile unicode/utf8=/usr/local/go/pkg/linux_amd64/unicode/utf8.a
//packagefile internal/bytealg=/usr/local/go/pkg/linux_amd64/internal/bytealg.a
//packagefile internal/cpu=/usr/local/go/pkg/linux_amd64/internal/cpu.a
//packagefile runtime/internal/atomic=/usr/local/go/pkg/linux_amd64/runtime/internal/atomic.a
//packagefile runtime/internal/math=/usr/local/go/pkg/linux_amd64/runtime/internal/math.a
//packagefile runtime/internal/sys=/usr/local/go/pkg/linux_amd64/runtime/internal/sys.a
//packagefile internal/reflectlite=/usr/local/go/pkg/linux_amd64/internal/reflectlite.a
//packagefile sort=/usr/local/go/pkg/linux_amd64/sort.a
//packagefile math/bits=/usr/local/go/pkg/linux_amd64/math/bits.a
//packagefile internal/oserror=/usr/local/go/pkg/linux_amd64/internal/oserror.a
//packagefile internal/poll=/usr/local/go/pkg/linux_amd64/internal/poll.a
//packagefile internal/syscall/execenv=/usr/local/go/pkg/linux_amd64/internal/syscall/execenv.a
//packagefile internal/syscall/unix=/usr/local/go/pkg/linux_amd64/internal/syscall/unix.a
//packagefile internal/testlog=/usr/local/go/pkg/linux_amd64/internal/testlog.a
//packagefile sync/atomic=/usr/local/go/pkg/linux_amd64/sync/atomic.a
//packagefile syscall=/usr/local/go/pkg/linux_amd64/syscall.a
//packagefile time=/usr/local/go/pkg/linux_amd64/time.a
//packagefile internal/unsafeheader=/usr/local/go/pkg/linux_amd64/internal/unsafeheader.a
//packagefile unicode=/usr/local/go/pkg/linux_amd64/unicode.a
//packagefile internal/race=/usr/local/go/pkg/linux_amd64/internal/race.a
//EOF
//mkdir -p $WORK/b001/exe/
//cd .

// 执行链接器, 生成最终可执行文件 link, 同时将可执行文件复制到当前路径下并
// 删除临时文件
///usr/local/go/pkg/tool/linux_amd64/link -o $WORK/b001/exe/a.out -importcfg $WORK/b001/importcfg.link -buildmode=exe -buildid=hvaQq6mxRKFr4Q6WRGte/m15EitMB9lGrgk-sJIEh/N4_QsJw-7lX_85DN_P3c/hvaQq6mxRKFr4Q6WRGte -extld=gcc $WORK/b001/_pkg_.a
///usr/local/go/pkg/tool/linux_amd64/buildid -w $WORK/b001/exe/a.out # internal
//mv $WORK/b001/exe/a.out link
//rm -r $WORK/b001/
