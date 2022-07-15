package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

/*
	./direct_rw_adt_fmt.go 中 fmt.Fscanf系列函数的运作本质是扫描和解析读出
	的文本字符串, 这导致其数据还原能力有局限:
	无法将从文件中读取的数据直接整体填充到抽象数据类型实例中, 只能逐个字段填充

	在数据还原方面, 二进制编码有着先天优势, 可以进行整体填充(TODO: binary 包)

*/

type Player struct {
	Name   [20]byte
	Age    int16
	Gender [6]byte
}

func directWriteADTToFile(path string, players []Player) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("open file error: %s", err)
	}

	for _, player := range players {
		// TODO: BigEndian, LittleEndian
		err = binary.Write(f, binary.BigEndian, &player)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var players [3]Player
	copy(players[0].Name[:], []byte("Tommy"))
	players[0].Age = 18
	copy(players[0].Gender[:], []byte("male"))

	copy(players[1].Name[:], []byte("Lucy"))
	players[1].Age = 17
	copy(players[1].Gender[:], []byte("female"))

	copy(players[2].Name[:], []byte("George"))
	players[2].Age = 19
	copy(players[2].Gender[:], []byte("male"))

	err := directWriteADTToFile("players.dat", players[:])
	if err != nil {
		fmt.Println("write file error: ", err)
		return
	}

	f, err := os.Open("players.dat")
	if err != nil {
		fmt.Println("open file error: ", err)
		return
	}

	var player Player
	for {
		err := binary.Read(f, binary.BigEndian, &player)
		if err == io.EOF {
			fmt.Println("read meet EOF")
			return
		}

		if err != nil {
			fmt.Println("read file error: ", err)
			return
		}
		fmt.Printf("%s %d %s\n", player.Name, player.Age, player.Gender)
	}

	/*
		Player类型定义中的字段都变成了导出字段(首字母大写)
		Player类型定义中的各个字段的类型都采用了定长类型, 比如string换成
		字节数组(这给结构体实例的初始化带来了一些麻烦), int换成了int16(int
		这个类型在不同CPU架构下长度可能不同), 这是binary包对直接操作的抽象
		数据类型的约束;
	*/
}
