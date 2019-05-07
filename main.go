package main

import (
	"SnakeMasterBot/neuron"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	host = "http://84.201.140.232:8080"
	lrud = "lrud/"
)

type respData struct {
	Answer  string
	Session string
	Data    struct {
		Area   [][]int
		Snakes []struct {
			Num    int
			Body   []neuron.Cell
			Energe int
			Dead   bool
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for n := 0; n < 1; n++ {
		gameBot("masterdak" + strconv.Itoa(n) + "r" + strconv.Itoa(rand.Intn(100)))
	}
	fmt.Scanln()
}
func gameBot(name string) {
	go func() {
		resp, session := autorizations(name)
		time.Sleep(200 * time.Millisecond)

		var data respData = respData{}
		var SnakeBot, SnakeBotLast neuron.Snake
		var World neuron.World
		SnakeBot.NeuroNetCreate()

		SnakeBotLast.Body = make([]neuron.Cell, 1)

		for {

			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&data)
			if err != nil {
				panic(err)
			}

			move := ""

			if data.Answer == "Name is busy." {
				return
			}

			World.Area = data.Data.Area
			World.LenX = len(data.Data.Area)
			World.LenY = len(data.Data.Area[0])

			for n := range data.Data.Snakes {
				SnakeBot.Body = data.Data.Snakes[n].Body
				SnakeBot.Energe = data.Data.Snakes[n].Energe

				if SnakeBot.Energe == SnakeBotLast.Energe {
					continue
				}

				standing := false
				if SnakeBot.Body[0] == SnakeBotLast.Body[0] {
					standing = true
				}

				if standing {
					SnakeBot.NeuroCorrect(&World, 0.1)
				} else {
					if len(SnakeBot.Body) > len(SnakeBotLast.Body) {
						SnakeBot.NeuroCorrect(&World, 0.99)
					} else {
						SnakeBot.NeuroCorrect(&World, 0.5)
					}
				}

				SnakeBot.NeuroSetIn(&World)
				m := SnakeBot.NeuroWay(&World)
				rs := lrud[m : m+1]

				move += rs

				SnakeBotLast.Energe = SnakeBot.Energe
				copy(SnakeBotLast.Body, SnakeBot.Body)
			}

			resp.Body.Close()

			resp, err = http.Get(host + "/game/?user=" + name + "&session=" + session + "&move=" + move)
			if err != nil {
				panic(err)
			}

			time.Sleep(50 * time.Millisecond)
		}
	}()
}

func autorizations(name string) (resp *http.Response, s string) {
	time.Sleep(1 * time.Second)
	resp, err := http.Get(host + "/game/?user=" + name)
	if err != nil {
		panic(err)
	}

	var data respData = respData{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	fmt.Println(data.Answer)

	s = data.Session

	resp, err = http.Get(host + "/game/?user=" + name + "&session=" + s)
	if err != nil {
		panic(err)
	}
	return
}
