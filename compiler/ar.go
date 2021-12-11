package main

/*
 * ar 文件是一种非常简单的打包文件格式, 广泛用于 Linux 的静态链接库中, 文件
 * 以字符串 "!\n" 开头, 随后是 60 字节的文件头部(包含文件名, 修改时间等信息),
 * 之后是文件内容; 因为 ar 文件格式简单, 所以 go 编译器直接在函数中实验了 ar
 * 打包过程($GOROOT/src/cmd/compile/internal/gc/obj.go);
 * startArchiveEntry 用于预留 ar 文件头信息的位置(60字节),
 * finishArchiveEntry 用于写入文件头信息, 因为文件头信息中包含文件大小,
 * 在写入完成之前文件大小未知, 所以分两步完成.(TODO: 阅读实现代码)
 *
 */

// func dumpobj1(outfile string, mode int) {
// 	bout, err := bio.Create(outfile)
// 	if err != nil {
// 		flusherrors()
// 		fmt.Printf("can't create %s: %v\n", outfile, err)
// 		errorexit()
// 	}
// 	defer bout.Close()
// 	bout.WriteString("!<arch>\n")

// 	if mode&modeCompilerObj != 0 {
// 		start := startArchiveEntry(bout)
// 		dumpCompilerObj(bout)
// 		finishArchiveEntry(bout, start, "__.PKGDEF")
// 	}
// 	if mode&modeLinkerObj != 0 {
// 		start := startArchiveEntry(bout)
// 		dumpLinkerObj(bout)
// 		finishArchiveEntry(bout, start, "_go_.o")
// 	}
// }

// func startArchiveEntry(bout *bio.Writer) int64 {
// 	var arhdr [ArhdrSize]byte
// 	bout.Write(arhdr[:])
// 	return bout.Offset()
// }

// func finishArchiveEntry(bout *bio.Writer, start int64, name string) {
// 	bout.Flush()
// 	size := bout.Offset() - start
// 	if size&1 != 0 {
// 		bout.WriteByte(0)
// 	}
// 	bout.MustSeek(start-ArhdrSize, 0)

// 	var arhdr [ArhdrSize]byte
// 	formathdr(arhdr[:], name, size)
// 	bout.Write(arhdr[:])
// 	bout.Flush()
// 	bout.MustSeek(start+size+(size&1), 0)
// }
