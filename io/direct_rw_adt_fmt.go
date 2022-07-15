package main

import (
	"fmt"
	"io"
	"os"
)

/*
	fmt.Fprint系列函数支持将任意抽象数据类型实例按照format参数的格式写入
	某io.Writer实例, 之后通过fmt.Fscan系列函数可以将从io.Reader中读取的数据还原
*/

type Player struct {
	name   string
	age    int
	gender string
}

func (p Player) String() string {
	return fmt.Sprintf("%s %d %s", p.name, p.age, p.gender)
}

func directWriteADTToFile(path string, players []Player) error {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("open file error:", err)
		return err
	}

	defer func() {
		// 在关闭文件句柄前执行了一次Sync操作, 该操作会完成数据落盘
		// 将尚处于内存中的数据写入磁盘
		f.Sync()
		f.Close()
	}()

	for _, player := range players {
		_, err := fmt.Fprintf(f, "%s\n", player)
		if err != nil {
			return err
		}
	}

	return nil
}

func mainForWrite() {
	var players = []Player{
		{"Tommy", 18, "male"},
		{"Lucy", 17, "female"},
		{"George", 19, "male"},
	}

	err := directWriteADTToFile("players.dat", players)
	if err != nil {
		fmt.Println("write file error: ", err)
		return
	}
	// 由于Player结构体类型实现了String方法, 当使用fmt.Fprintf向文件
	// (io.Writer实例)写入数据时(通过%s), Player类型的String方法便会
	// 被调用(TODO: 源码实现), 因此每行的数据格式均为自定义的name age gender
}

func main() {
	f, err := os.Open("players.dat")
	if err != nil {
		fmt.Println("open file error: ", err)
		return
	}
	defer f.Close()

	var player Player
	for {
		// _, err := fmt.Fscanf(f, "%s %d %s", &player)
		// fmt.Fscanf仅支持扫描原生类型或底层类型为原生类型的数据
		_, err := fmt.Fscanf(f, "%s %d %s", &player.name, &player.age, &player.gender)
		if err == io.EOF {
			fmt.Println("read meet EOF")
			return
		}

		if err != nil {
			fmt.Println("read file error: ", err) // read file error:  can't scan type: *main.Player
			/*
				不支持对*main.Player类型进行扫描
				fmt.Fscanf仅支持扫描原生类型或底层类型为原生类型的数据
			*/
			return
		}

		fmt.Printf("%s %d %s\n", player.name, player.age, player.gender)
	}
}

/*
	格式化输出/输入函数Printf系列和Scanf系列本质上是当io.Writer为
	os.Stdout或io.Reader为os.Stdin时的特例

	// Printf formats according to a format specifier and writes to standard output.
	// It returns the number of bytes written and any write error encountered.
	func Printf(format string, a ...any) (n int, err error) {
		return Fprintf(os.Stdout, format, a...)
	}

	// Scanf scans text read from standard input, storing successive
	// space-separated values into successive arguments as determined by
	// the format. It returns the number of items successfully scanned.
	// If that is less than the number of arguments, err will report why.
	// Newlines in the input must match newlines in the format.
	// The one exception: the verb %c always scans the next rune in the
	// input, even if it is a space (or tab etc.) or newline.
	func Scanf(format string, a ...any) (n int, err error) {
		return Fscanf(os.Stdin, format, a...)
	}

*/
