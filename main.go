package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"math"

	"github.com/sethgrid/curse"
)

func getRandomNumber(c chan int) {
	num := rand.Intn(100)

	wait := rand.Intn(10)
	//fmt.Println("wait", wait, num)
	time.Sleep(time.Second * time.Duration(wait))
	c <- num
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	randomNum := make(chan int)
	const maxSpawns = 100

	cur, err := curse.New()
	if err != nil {
		log.Fatal(err)
	}

	cur.SetColor(curse.RED).SetBackgroundColor(curse.BLACK)
	fmt.Println("Channel Spawner")
	cur.SetColor(curse.WHITE)

	fmt.Println("Spawning")
	i := 1
	for i <= maxSpawns {
		go getRandomNumber(randomNum)

		cur.MoveUp(1).EraseCurrentLine()
		fmt.Println("Spawn", i)

		i = i + 1
	}

	fmt.Printf("Launched all %d spawns", maxSpawns)

	var num int
	min := math.MaxInt32
	max := 0
	avg := float64(0)
	i = 1
	fmt.Println("Waiting...")
	for i <= maxSpawns {
		num = <-randomNum

		if num > max {
			max = num
		}

		if num < min {
			min = num
		}

		avg = ((avg * float64(i-1)) + float64(num)) / float64(i)

		cur.MoveUp(1).EraseCurrentLine()
		fmt.Printf("num: %d min: %d max: %d avg: %s [%d/%d returned]\n", num, min, max, strconv.FormatFloat(avg, 'f', 0, 64), i, maxSpawns)
		i = i + 1
	}
	fmt.Println("Done")
	cur.SetDefaultStyle()
}
