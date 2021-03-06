package neuron

import (
	"SnakeMasters/server"
	"fmt"
	"math"
	"neuron/nr"
)

const (
	viewRange = 4
	viewLen   = 1 + viewRange*2
	lenMomory = viewRange * 2
	dirWay    = 5
)

type World struct {
	LenY, LenX int
	Area       [][]int
}

type Snake struct {
	Body     []Cell
	Divider  int
	Energe   int
	Standing int
	nCorrect float64
	memory   memory
	NeuroNet nr.NeuroNet
}

type Cell struct {
	X int
	Y int
}

type memory struct {
	data [lenMomory][]nr.Neuron
	way  [lenMomory]int
	pos  int
}

func (s *Snake) NeuroNetCreate() {
	s.Divider = 8
	s.nCorrect = 0.2

	neuroLayer := make([]int, 2)
	neuroLayer[0] = viewLen * viewLen
	neuroLayer[1] = dirWay

	for n := range s.memory.data {
		s.memory.data[n] = make([]nr.Neuron, viewLen*viewLen)
	}

	s.NeuroNet.CreateLayer(neuroLayer)
}

func (s *Snake) NeuroSetIn(w *World) {
	x := s.Body[0].X
	y := s.Body[0].Y
	x0 := x - viewRange
	x1 := x + viewRange
	y0 := y - viewRange
	y1 := y + viewRange

	dOut := float64(0.01)

	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if y < 0 || y >= w.LenY || x < 0 || x >= w.LenX {
				dOut = -0.5 //Выход за край карты
			} else {
				dOut, _ = s.dataToOut(w, w.Area[x][y])
			}

			dx := x - x0
			dy := y - y0

			n := dy*viewLen + dx
			s.NeuroNet.Layers[0][n].Out = dOut
		}
	}

	copy(s.memory.data[s.memory.pos], s.NeuroNet.Layers[0])
}

func (s *Snake) NeuroWay(w *World) int {
	s.NeuroSetIn(w)
	s.NeuroNet.Calc()
	mo := s.NeuroNet.MaxOutputNumber(0)
	s.memory.way[s.memory.pos] = mo
	s.memory.pos = (s.memory.pos + 1) % lenMomory
	return mo
}

func (s *Snake) NeuroCorrect(w *World, a float64) {
	ans := make([]float64, dirWay)
	way := 0

	pow := float64(0)
	for pos := s.memory.pos + lenMomory - 1; pos >= s.memory.pos; pos-- {
		p := pos % lenMomory

		s.NeuroNet.NCorrect = 0.5 / math.Pow(2, pow)
		pow++

		s.NeuroNet.Layers[0] = s.memory.data[p]
		s.NeuroNet.Calc()

		for n := 0; n < dirWay; n++ {
			ans[n] = s.NeuroNet.Layers[len(s.NeuroNet.Layers)-1][n].Out
		}

		way = s.memory.way[p]

		ans[way] = a

		s.NeuroNet.SetAnswers(ans)
		s.NeuroNet.Correct()
	}
}

func (s *Snake) dataToOut(w *World, data int) (d float64, str string) {
	switch data {
	case server.ElWall:
		return -0.5, "# "
	case server.ElEmpty:
		return 0.01, ". "
	case server.ElEat:
		return 0.99, "* "
	case server.ElHead:
		return -0.6, "@ "
	case server.ElBody:
		return -0.55, "o "
	}

	panic("dataToOut: Пустое значение.")
	return 0, "  " //
}

func printData(layer []nr.Neuron) {
	fmt.Println()
	for n := range layer {
		if n%viewLen == 0 {
			fmt.Println()
		}
		switch layer[n].Out {
		case -0.5:
			fmt.Print("# ")
		case 0.01:
			fmt.Print(". ")
		case 0.99:
			fmt.Print("* ")
		case -0.6:
			fmt.Print("@ ")
		case -0.55:
			fmt.Print("o ")
		}
	}
}
