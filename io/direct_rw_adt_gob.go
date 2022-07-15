package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

/*
	虽然binary包实现了抽象数据类型实例的直接读写, 但只支持采用定长表示的
	抽象数据类型, 限制了其应用范围; Go标准库提供了更为通用的选择: gob包,
	gob包支持对任意抽象数据类型实例的直接读写, 唯一的约束是自定义结构体
	类型中的字段至少有一个是导出的(字段名首字母大写); (TODO: 至少一个导出?)

	gob包也是Go标准库提供的一个序列化/反序列化方案, 和JSON、XML等序列化/反
	序列化方案不同, 它的API直接支持读写实现了io.Reader和io.Writer接口的实例;
	TODO: gob 包
*/

type Player struct {
	Name   string
	Age    int
	Gender string
}

func directWriteADTToFile(path string, players []Player) error {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("open file error:", err)
		return err
	}

	enc := gob.NewEncoder(f)
	for _, player := range players {
		err = enc.Encode(player)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
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

	f, err := os.Open("players.dat")
	if err != nil {
		fmt.Println("open file error: ", err)
		return
	}

	var player Player
	dec := gob.NewDecoder(f)
	for {
		err := dec.Decode(&player)
		if err != nil {
			fmt.Println("read meet EOF")
			return
		}

		if err != nil {
			fmt.Println("read file error: ", err)
			return
		}
		fmt.Printf("%v\n", player)
	}
}
