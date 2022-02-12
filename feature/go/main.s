"".myFunction STEXT nosplit size=49 args=0x20 locals=0x0
	0x0000 00000 (main.go:3)	TEXT	"".myFunction(SB), NOSPLIT|ABIInternal, $0-32
	0x0000 00000 (main.go:3)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (main.go:3)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (main.go:3)	MOVQ	$0, "".~r2+24(SP) # 初始化第一个返回值
	0x0009 00009 (main.go:3)	MOVQ	$0, "".~r3+32(SP) # 初始化第二个返回值
	0x0012 00018 (main.go:4)	MOVQ	"".a+8(SP), AX
	0x0017 00023 (main.go:4)	ADDQ	"".b+16(SP), AX
	0x001c 00028 (main.go:4)	MOVQ	AX, "".~r2+24(SP)
	0x0021 00033 (main.go:4)	MOVQ	"".a+8(SP), AX
	0x0026 00038 (main.go:4)	SUBQ	"".b+16(SP), AX
	0x002b 00043 (main.go:4)	MOVQ	AX, "".~r3+32(SP)
	0x0030 00048 (main.go:4)	RET
	0x0000 48 c7 44 24 18 00 00 00 00 48 c7 44 24 20 00 00  H.D$.....H.D$ ..
	0x0010 00 00 48 8b 44 24 08 48 03 44 24 10 48 89 44 24  ..H.D$.H.D$.H.D$
	0x0020 18 48 8b 44 24 08 48 2b 44 24 10 48 89 44 24 20  .H.D$.H+D$.H.D$ 
	0x0030 c3                                               .
"".main STEXT size=71 args=0x0 locals=0x28
	0x0000 00000 (main.go:7)	TEXT	"".main(SB), ABIInternal, $40-0
	0x0000 00000 (main.go:7)	MOVQ	(TLS), CX
	0x0009 00009 (main.go:7)	CMPQ	SP, 16(CX)
	0x000d 00013 (main.go:7)	PCDATA	$0, $-2
	0x000d 00013 (main.go:7)	JLS	64
	0x000f 00015 (main.go:7)	PCDATA	$0, $-1
	0x000f 00015 (main.go:7)	SUBQ	$40, SP # 分配40字节的栈空间(基址地址: 8Byte, 两个返回值: 16 字节, 两个参数: 16字节)
	0x0013 00019 (main.go:7)	MOVQ	BP, 32(SP) # 将基址地址存储到栈上, main函数的栈基址地址占8字节
	0x0018 00024 (main.go:7)	LEAQ	32(SP), BP
	0x001d 00029 (main.go:7)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x001d 00029 (main.go:7)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x001d 00029 (main.go:8)	MOVQ	$66, (SP)  # 第一个参数(golang 64位cpu, int 占64位)
	0x0025 00037 (main.go:8)	MOVQ	$77, 8(SP) # 第二个参数
												   # SP+32 -- BP: 8Byte
									               # SP+16 -- SP+32: 两个返回值
												   # SP    -- SP+16: 两个参数
	0x002e 00046 (main.go:8)	PCDATA	$1, $0
	0x002e 00046 (main.go:8)	CALL	"".myFunction(SB) # TODO: SB?
	
								# 当 myFunction 函数返回后, main 函数会通过
								# 以下指令来恢复栈基址地址并销毁已经失去作用
								# 的 40 字节栈内存----------------------------
	0x0033 00051 (main.go:9)	MOVQ	32(SP), BP
	0x0038 00056 (main.go:9)	ADDQ	$40, SP
	0x003c 00060 (main.go:9)	RET
								# --------------------------------------------
	0x003d 00061 (main.go:9)	NOP
	0x003d 00061 (main.go:7)	PCDATA	$1, $-1
	0x003d 00061 (main.go:7)	PCDATA	$0, $-2
	0x003d 00061 (main.go:7)	NOP
	0x0040 00064 (main.go:7)	CALL	runtime.morestack_noctxt(SB)
	0x0045 00069 (main.go:7)	PCDATA	$0, $-1
	0x0045 00069 (main.go:7)	JMP	0
	0x0000 64 48 8b 0c 25 00 00 00 00 48 3b 61 10 76 31 48  dH..%....H;a.v1H
	0x0010 83 ec 28 48 89 6c 24 20 48 8d 6c 24 20 48 c7 04  ..(H.l$ H.l$ H..
	0x0020 24 42 00 00 00 48 c7 44 24 08 4d 00 00 00 e8 00  $B...H.D$.M.....
	0x0030 00 00 00 48 8b 6c 24 20 48 83 c4 28 c3 0f 1f 00  ...H.l$ H..(....
	0x0040 e8 00 00 00 00 eb b9                             .......
	rel 5+4 t=17 TLS+0
	rel 47+4 t=8 "".myFunction+0
	rel 65+4 t=8 runtime.morestack_noctxt+0
go.cuinfo.packagename. SDWARFINFO dupok size=0
	0x0000 6d 61 69 6e                                      main
""..inittask SNOPTRDATA size=24
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 00 00 00 00 00 00 00 00                          ........
gclocals·33cdeccccebe80329f1fdbee7f5874cb SRODATA dupok size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
