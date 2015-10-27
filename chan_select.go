package main
import (
	"time"
	"fmt"
)

func sleep(n time.Duration) {
	<-time.After(time.Second * n)
}

func main() {
	c1 := make (chan string, 1)
	c2 := make (chan string, 1)

	go func() {
		for {
			c1 <- "from 1"
			sleep(2)
		}
	}()

	go func() {
		for {
			c1 <- "from 2"
			sleep(3)
		}
	}()

	go func() {
		for {
			select {
			case msg1 := <- c1:
				fmt.Println("Message 1", msg1)
			case msg2 := <- c2:
				fmt.Println("Message 2", msg2)
			case <- time.After(time.Second):
				fmt.Println("timeout")
//			default:
//				fmt.Println("nothing ready")
			}
		}
	}()

	var input string
	fmt.Scanln(&input)
}
