package main

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"math"

	"github.com/sethgrid/curse"
)

var ActivatableMethods []string = []string{
	"IncrementF",
	"DecrementF",
	"IncrementM",
	"DecrementM",
	"IncrementT",
	"DecrementT",
}

type Creature struct {
	Species   string
	Activated []string
	F         int
	M         int
	T         int
	HP        int
}

func (c *Creature) IncrementF() {
	c.F = c.F + 1
}
func (c *Creature) DecrementF() {
	c.F = c.F - 1
}
func (c *Creature) IncrementM() {
	c.M = c.M + 1
}
func (c *Creature) DecrementM() {
	c.M = c.M - 1
}
func (c *Creature) IncrementT() {
	c.T = c.T + 1
}
func (c *Creature) DecrementT() {
	c.T = c.T - 1
}
func (c *Creature) IsActivated(methStr string) bool {
	for _, item := range c.Activated {
		if item == methStr {
			return true
		}
	}
	return false
}

func NewCapability() Creature {
	c := Creature{
		randomName(6),
		[]string{},
		rand.Intn(100),
		rand.Intn(100),
		rand.Intn(100),
		1000,
	}
	amount := rand.Intn(len(ActivatableMethods) + 1)
	indexBag := rand.Perm(6)
	// choose to grab a certain amount of methods to activate
	for i := 1; i <= amount; i++ {
		c.Activated = append(c.Activated, ActivatableMethods[indexBag[i-1]])
	}
	return c
}

type Environment struct {
	F int
	M int
	T int
}

func (env *Environment) Shift() {
	env.F = env.F + (rand.Intn(3) - 1) // -1 to 1
	env.M = env.M + (rand.Intn(3) - 1)
	env.T = env.T + (rand.Intn(3) - 1)
}

func randomName(strlen int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func getRandom(c chan Creature) {
	newC := NewCapability()

	wait := rand.Intn(3)
	//fmt.Println("wait", wait, num)
	time.Sleep(time.Second * time.Duration(wait))
	c <- newC
}

func tick(env *Environment, capa Creature) <-chan Creature {
	out := make(chan Creature, 1)

	go func() {

		//fmt.Printf("tick requested with %+v\n", capa)
		// F
		diff := env.F - capa.F
		if diff > 0 { // capa was less
			if rand.Intn(2) == 0 && capa.IsActivated("IncrementF") { // 50:50
				f := reflect.ValueOf(&capa).MethodByName("IncrementF")
				fmt.Printf("+F %+v\n", capa)
				f.Call(nil)
			}
		} else if diff < 0 { // capa was more
			if rand.Intn(2) == 0 && capa.IsActivated("DecrementF") { // 50:50
				f := reflect.ValueOf(&capa).MethodByName("DecrementF")
				fmt.Printf("-F %+v\n", capa)
				f.Call(nil)
			}
		}
		delta := int(math.Abs(float64(env.F - capa.F)))
		capa.HP = capa.HP - delta

		// M
		diff = env.M - capa.M
		if diff > 0 { // capa was less
			if rand.Intn(2) == 0 && capa.IsActivated("IncrementM") { // 50:50
				f := reflect.ValueOf(&capa).MethodByName("IncrementM")
				fmt.Printf("+M %+v\n", capa)
				f.Call(nil)
			}
		} else if diff < 0 { // capa was more
			if rand.Intn(2) == 0 && capa.IsActivated("DecrementM") { // 50:50
				f := reflect.ValueOf(&capa).MethodByName("DecrementM")
				fmt.Printf("-M %+v\n", capa)
				f.Call(nil)
			}
		}
		delta = int(math.Abs(float64(env.M - capa.M)))
		capa.HP = capa.HP - delta

		// T
		diff = env.T - capa.T
		if diff > 0 { // capa was less
			if rand.Intn(2) == 0 && capa.IsActivated("IncrementT") { // 50:50
				f := reflect.ValueOf(&capa).MethodByName("IncrementT")
				fmt.Printf("+T %+v\n", capa)
				f.Call(nil)
			}
		} else if diff < 0 { // capa was more
			if rand.Intn(2) == 0 && capa.IsActivated("DecrementT") { // 50:50
				f := reflect.ValueOf(&capa).MethodByName("DecrementT")
				fmt.Printf("-T %+v\n", capa)
				f.Call(nil)
			}
		}
		delta = int(math.Abs(float64(env.T - capa.T)))
		capa.HP = capa.HP - delta

		out <- capa
	}()
	return out
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	aCap := make(chan Creature)
	const maxSpawns = 10

	cur, err := curse.New()
	if err != nil {
		log.Fatal(err)
	}

	cur.SetColor(curse.RED).SetBackgroundColor(curse.BLACK)
	fmt.Println("Channel Spawner")
	cur.SetColor(curse.WHITE)

	// Setup
	fmt.Println("Spawning")
	i := 1
	for i <= maxSpawns {
		go getRandom(aCap)

		cur.MoveUp(1).EraseCurrentLine()
		fmt.Println("Spawn", i)

		i = i + 1
	}

	fmt.Printf("Launched all %d spawns", maxSpawns)

	caps := []Creature{}
	min := math.MaxInt32
	max := 0
	avg := float64(0)
	i = 1
	fmt.Println("Waiting...")
	for i <= maxSpawns {
		capa := <-aCap
		caps = append(caps, capa)
		methods := len(capa.Activated)

		if methods > max {
			max = methods
		}

		if methods < min {
			min = methods
		}

		avg = ((avg * float64(i-1)) + float64(methods)) / float64(i)

		cur.MoveUp(1).EraseCurrentLine()
		fmt.Printf("num: %d min: %d max: %d avg: %s [%d/%d returned]\n", methods, min, max, strconv.FormatFloat(avg, 'f', 0, 64), i, maxSpawns)
		i = i + 1
	}
	fmt.Println("Done")
	close(aCap)
	fmt.Printf("Creatures: %+v\n", caps)

	fmt.Println("Making Environment")
	env := Environment{
		rand.Intn(100),
		rand.Intn(100),
		rand.Intn(100),
	}
	fmt.Printf("ENV: %+v\n", env)

	// run
	year := 1
	for len(caps) > 1 {
		fmt.Printf("\n=== year %d ===\n", year)
		retCaps := []Creature{}
		for _, capa := range caps {
			retCaps = append(retCaps, <-tick(&env, capa))
		}
		fmt.Println("Sent all ticks")

		tempCaps := []Creature{}
		for _, capa := range retCaps {
			//fmt.Printf("Got back result: %+v\n", capa)
			if capa.HP > 0 {
				fmt.Print("!")
				tempCaps = append(tempCaps, capa)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")

		fmt.Printf("Before Creatures: %+v\n", caps)
		fmt.Printf("Env: %+v\n", env)
		fmt.Printf("After  Creatures: %+v\n", tempCaps)

		caps = tempCaps
		env.Shift()
		year = year + 1
	}

	fmt.Printf("Last Creature Standing: %+v\n", caps)
	fmt.Printf("Final ENV: %+v\n", env)

	cur.SetDefaultStyle()
}
