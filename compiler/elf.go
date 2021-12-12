package main

/*
 * ./link.go 编译生成的可执行文件 ./link 是 ELF 格式的文件
 * ELF 文件解析
 * 在 windows 下编译后的 go 文本文件最终会生成以 .exe 为后缀的 PE 格式的可执行
 * 文件, 在 linux 或类 Unix 操作系统下, 会生成 ELF 格式的可执行文件; 除机器码
 * 外, 在可执行文件中还可能包含调试信息, 动态链接库信息, 符号表信息;
 * ELF(Executable and Linkable Format) 是类 UNIX 操作系统下最常见的可执行且可
 * 链接的文件格式; readelf, objdump 等工具可以查看 ELF 文件的头信息.
 *
 */

// >> readelf -h ./link  // 通过 readelf 查看 elf 文件的头信息
// ELF Header:
//  Magic:   7f 45 4c 46 02 01 01 00 00 00 00 00 00 00 00 00
//  Class:                             ELF64
//  Data:                              2's complement, little endian
//  Version:                           1 (current)
//  OS/ABI:                            UNIX - System V
//  ABI Version:                       0
//  Type:                              EXEC (Executable file)
//  Machine:                           Advanced Micro Devices X86-64
//  Version:                           0x1
//  Entry point address:               0x464820
//  Start of program headers:          64 (bytes into file)
//  Start of section headers:          456 (bytes into file)
//  Flags:                             0x0
//  Size of this header:               64 (bytes)
//  Size of program headers:           56 (bytes)
//  Number of program headers:         7
//  Size of section headers:           64 (bytes)
//  Number of section headers:         25
//  Section header string table index: 3

// TODO:
// ELF 包含多个 segment 与 section; 标准库 debug/elf 包中提供了调试 ELF 的 API
import (
	"debug/elf"
	"log"
)

func main() {
	f, err := elf.Open("link")
	if err != nil {
		log.Fatal(err)
	}

	for _, section := range f.Sections {
		log.Println(section)
	}
}

// >> readelf -S ./link // 通过 readelf 查看 ELF 文件中 section 信息
// There are 25 section headers, starting at offset 0x1c8:

// Section Headers:
//   [Nr] Name              Type             Address           Offset
//        Size              EntSize          Flags  Link  Info  Align
//   [ 0]                   NULL             0000000000000000  00000000
//        0000000000000000  0000000000000000           0     0     0
//   [ 1] .text             PROGBITS         0000000000401000  00001000
//        000000000009810a  0000000000000000  AX       0     0     32
//   [ 2] .rodata           PROGBITS         000000000049a000  0009a000
//        0000000000043d66  0000000000000000   A       0     0     32
//   [ 3] .shstrtab         STRTAB           0000000000000000  000ddd80
//        00000000000001bc  0000000000000000           0     0     1
//   [ 4] .typelink         PROGBITS         00000000004ddf40  000ddf40
//        0000000000000734  0000000000000000   A       0     0     32
//   [ 5] .itablink         PROGBITS         00000000004de678  000de678
//        0000000000000050  0000000000000000   A       0     0     8
//   [ 6] .gosymtab         PROGBITS         00000000004de6c8  000de6c8
//        0000000000000000  0000000000000000   A       0     0     1
//   [ 7] .gopclntab        PROGBITS         00000000004de6e0  000de6e0
//        000000000005f945  0000000000000000   A       0     0     32
//   [ 8] .go.buildinfo     PROGBITS         000000000053f000  0013f000
//        0000000000000020  0000000000000000  WA       0     0     16
//   [ 9] .noptrdata        PROGBITS         000000000053f020  0013f020
//        000000000000e4a0  0000000000000000  WA       0     0     32
//   [10] .data             PROGBITS         000000000054d4c0  0014d4c0
//        0000000000007490  0000000000000000  WA       0     0     32
//   [11] .bss              NOBITS           0000000000554960  00154960
//        000000000002f8f0  0000000000000000  WA       0     0     32
//   [12] .noptrbss         NOBITS           0000000000584260  00184260
//        0000000000002e88  0000000000000000  WA       0     0     32
//   [13] .zdebug_abbrev    PROGBITS         0000000000588000  00155000
//        0000000000000119  0000000000000000           0     0     1
//   [14] .zdebug_line      PROGBITS         0000000000588119  00155119
//        000000000001bcc1  0000000000000000           0     0     1
//   [15] .zdebug_frame     PROGBITS         00000000005a3dda  00170dda
//        0000000000006277  0000000000000000           0     0     1
//   [16] .zdebug_pubnames  PROGBITS         00000000005aa051  00177051
//        00000000000014c0  0000000000000000           0     0     1
//   [17] .zdebug_pubtypes  PROGBITS         00000000005ab511  00178511
//        00000000000034d0  0000000000000000           0     0     1
//   [18] .debug_gdb_script PROGBITS         00000000005ae9e1  0017b9e1
//        000000000000002a  0000000000000000           0     0     1
//   [19] .zdebug_info      PROGBITS         00000000005aea0b  0017ba0b
//        0000000000033793  0000000000000000           0     0     1
//   [20] .zdebug_loc       PROGBITS         00000000005e219e  001af19e
//        0000000000016818  0000000000000000           0     0     1
//   [21] .zdebug_ranges    PROGBITS         00000000005f89b6  001c59b6
//        0000000000008cf3  0000000000000000           0     0     1
//   [22] .note.go.buildid  NOTE             0000000000400f9c  00000f9c
//        0000000000000064  0000000000000000   A       0     0     4
//   [23] .symtab           SYMTAB           0000000000000000  001cf000
//        00000000000112e0  0000000000000018          24   421     8
//   [24] .strtab           STRTAB           0000000000000000  001e02e0
//        000000000001097d  0000000000000000           0     0     1
// Key to Flags:
//   W (write), A (alloc), X (execute), M (merge), S (strings), I (info),
//   L (link order), O (extra OS processing required), G (group), T (TLS),
//   C (compressed), x (unknown), o (OS specific), E (exclude),
//   l (large), p (processor specific)

// segment 包含多个 section, 它描述程序如何映射到内存中, 如哪些 section 需要
// 导入内存, 采取只读模式还是读写模式, 内存对齐大小等; 以下是 section 与
// segment 的对应关系.
// >> readelf -lW link
// Elf file type is EXEC (Executable file)
// Entry point 0x464820
// There are 7 program headers, starting at offset 64

// Program Headers:
//   Type           Offset   VirtAddr           PhysAddr           FileSiz  MemSiz   Flg Align
//   PHDR           0x000040 0x0000000000400040 0x0000000000400040 0x000188 0x000188 R   0x1000
//   NOTE           0x000f9c 0x0000000000400f9c 0x0000000000400f9c 0x000064 0x000064 R   0x4
//   LOAD           0x000000 0x0000000000400000 0x0000000000400000 0x09910a 0x09910a R E 0x1000
//   LOAD           0x09a000 0x000000000049a000 0x000000000049a000 0x0a4025 0x0a4025 R   0x1000
//   LOAD           0x13f000 0x000000000053f000 0x000000000053f000 0x015960 0x0480e8 RW  0x1000
//   GNU_STACK      0x000000 0x0000000000000000 0x0000000000000000 0x000000 0x000000 RW  0x8
//   LOOS+0x5041580 0x000000 0x0000000000000000 0x0000000000000000 0x000000 0x000000     0x8

//  Section to Segment mapping:
//   Segment Sections...
//    00
//    01     .note.go.buildid
//    02     .text .note.go.buildid
//    03     .rodata .typelink .itablink .gosymtab .gopclntab
//    04     .go.buildinfo .noptrdata .data .bss .noptrbss
//    05
//    06

/*
 * TODO
 * 并非所有的 section 都需要导入内存, 当 Type 为 LOAD 时, 代表 section 需要被
 * 导入内存, 后面的 Flg 代表内存的读写模式; 包含 text 的代码区代表可以被读和
 * 执行, 包含 .data 与 .bss 的全局变量可以被读写; 其中, 为了满足垃圾回收的需要
 * 还区分了是否包含指针的区域; 包含 .rodata 常量数据的区域代表只读区, 其中
 * .itablink 为与 go 语言接口相关的全局符号表(TODO), .gopclntab 包含程序计数器
 * PC与源代码行的对应关系.
 */

// 一个 hello, world 程序一共包含了 25 个 section, 可看出并不是所有的 section
// 都需要导入内存; 同时, 该程序包含单独存储调试信息的区域, 如 .note.go.buildid
// 包含 go 程序唯一的 ID, 可通过 objdump 工具在 .note.go.buildid 中查找每个 go
// 程序唯一的 ID.(TODO)
// >> objdump -s -j ./link
// link:     file format elf64-x86-64

// Contents of section .note.go.buildid:
//  400f9c 04000000 53000000 04000000 476f0000  ....S.......Go..
//  400fac 61697038 71306e30 656c4159 4e386857  aip8q0n0elAYN8hW
//  400fbc 4a516d6d 2f786f67 46414342 3969324d  JQmm/xogFACB9i2M
//  400fcc 6f4d7257 62774544 612f4e34 5f51734a  oMrWbwEDa/N4_QsJ
//  400fdc 772d376c 585f3835 444e5f50 33632f33  w-7lX_85DN_P3c/3
//  400fec 39396558 34746c55 374b5359 36434a36  99eX4tlU7KSY6CJ6
//  400ffc 74343700                             t47.

// .go.buildinfo section 包含 go 程序的构建信息, go version 命令会查找该区域的
// 信息获取 go 语言版本号
// >> go version ./link
// link: go1.15.6
