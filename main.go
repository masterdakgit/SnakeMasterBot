package main

import (
	"SnakeMasterBot/neuron"
	"encoding/json"
	"fmt"
	"log"
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
	for n := 0; n < 5; n++ {
		gameBot("BotN" + strconv.Itoa(n) + "r" + strconv.Itoa(rand.Intn(100)))
	}
	fmt.Scanln()
}
func gameBot(name string) {
	go func() {
		resp, session := autorizations(name)
		time.Sleep(200 * time.Millisecond)

		var data respData = respData{}
		var SnakeBot, SnakeBotLast [300]neuron.Snake
		var World neuron.World

		for n := 0; n < 300; n++ {
			SnakeBotLast[n].Body = make([]neuron.Cell, 1)
			SnakeBot[n].NeuroNetCreate()
			SnakeBot[n].NeuroNet = SnakeBot[0].NeuroNet
		}

		for {

			for resp.Body == nil {
				time.Sleep(50 * time.Millisecond)
				var err error
				resp, err = http.Get(host + "/game/?user=" + name + "&session=" + session)
				if err != nil {
					log.Println(err)
				}
			}

			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&data)
			if err != nil {
				log.Println(err)
			}

			move := ""

			if data.Answer == "Name is busy." {
				return
			}

			World.Area = data.Data.Area
			World.LenX = len(data.Data.Area)
			if World.LenX == 0 {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			World.LenY = len(data.Data.Area[0])

			for n := range data.Data.Snakes {
				SnakeBot[n].Body = data.Data.Snakes[n].Body
				SnakeBot[n].Energe = data.Data.Snakes[n].Energe

				if SnakeBot[n].Energe == SnakeBotLast[n].Energe {
					continue
				}

				standing := false
				if SnakeBot[n].Body[0] == SnakeBotLast[n].Body[0] {
					standing = true
				}

				if standing {
					SnakeBot[n].NeuroCorrect(&World, 0.1)
				} else if len(SnakeBot[n].Body) > len(SnakeBotLast[n].Body) {
					SnakeBot[n].NeuroCorrect(&World, 0.99)
				} else if len(SnakeBot[n].Body) < len(SnakeBotLast[n].Body) {
					SnakeBot[n].NeuroCorrect(&World, 0.2)
				} else if data.Data.Snakes[n].Energe == 1 && len(SnakeBot[n].Body) == 1 {
					SnakeBot[n].NeuroCorrect(&World, 0.01)
				} else {
					SnakeBot[n].NeuroCorrect(&World, 0.5)
				}

				SnakeBot[n].NeuroSetIn(&World)
				m := SnakeBot[n].NeuroWay(&World)
				rs := lrud[m : m+1]

				/*
					if len(SnakeBot[n].Body) > 8 {
						rs = "/"
					}
				*/
				move += rs

				SnakeBotLast[n].Energe = SnakeBot[n].Energe
				copy(SnakeBotLast[n].Body, SnakeBot[n].Body)
			}

			resp.Body.Close()

			resp, err = http.Get(host + "/game/?user=" + name + "&session=" + session + "&move=" + move)
			if err != nil {
				log.Println(err)
			}

			time.Sleep(50 * time.Millisecond)
		}
	}()
}

func autorizations(name string) (resp *http.Response, s string) {
	time.Sleep(1 * time.Second)
	resp, err := http.Get(host + "/game/?user=" + name)
	if err != nil {
		log.Println(err)
	}

	var data respData = respData{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(data.Answer)

	s = data.Session

	resp, err = http.Get(host + "/game/?user=" + name + "&session=" + s)
	if err != nil {
		log.Println(err)
	}
	return
}
