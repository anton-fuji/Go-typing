package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/tjarratt/babble"
)

func main() {
	var (
		ch_rcv = myInput(os.Stdin)
		t      = 60
		tick   = time.NewTicker(30 * time.Second)
		done   = make(chan bool)
		n      = 0
	)
	// 文字をランダムで生成
	babbler := babble.NewBabbler()
	babbler.Count = 1

	fmt.Printf("Start the typing game. Time limit is %d seconds. Let's Start!\n", t)

	// 30秒経過後の残り時間表示のためのゴルーチンを開始
	go func() {
		for {
			select {
			case <-tick.C:
				t -= 30
				if t > 0 {
					fmt.Printf("%d seconds left\n", t)
				}
			case <-done:
				return
			}
		}
	}()

OuterLoop:
	for {
		q := babbler.Babble()
		fmt.Println(q)

		select {
		case <-time.After(time.Duration(t) * time.Second):
			fmt.Printf("Finished! Your score is %d points! Good job\n", n)
			done <- true
			break OuterLoop
		case x := <-ch_rcv:
			x = strings.TrimSpace(x)
			if x == q {
				fmt.Println("OK!")
				n += 1
			} else {
				fmt.Println("NG")
			}
		}
	}
}

/*
標準入力をforの無限ループで受付け、
戻り値を受信専チャネルに送る
*/
func myInput(r io.Reader) <-chan string {
	ch := make(chan string)
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			ch <- s.Text()
		}
		if err := s.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Read Error", err)
		}
		close(ch)
	}()
	return ch
}
