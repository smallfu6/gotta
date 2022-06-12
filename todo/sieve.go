package main

import "fmt"

/*
	素数筛
	本例采用埃拉托斯特尼素数筛算法: 先用最小的素数2去筛, 把2的倍数筛除; 下一个
	未筛除的数就是素数(3), 再用这个素数3区筛, 筛除3的倍数... 这样不断重复下去,
	直到筛完为止;
*/

func Generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i
		fmt.Println("gen", i)

	}
}

func Filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in
		if i%prime != 0 {
			out <- i
		}
	}
}

func main() {
	ch := make(chan int)
	go Generate(ch)
	for i := 0; i < 10; i++ {
		prime := <-ch
		fmt.Println(prime)
		ch1 := make(chan int)
		go Filter(ch, ch1, prime)
		ch = ch1
	}
}
