package classic

import (
	"../../goslot"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Model struct {
	conf     *goslot.Conf
	paylines [][]int
	paytable [][]int
}

func NewModel(conf *goslot.Conf, paylines [][]int, paytable [][]int) *Model {
	if paylines == nil || len(paylines) == 0 {
		panic("invalid pay table")
	}

	for i := range paylines {
		if paylines[i] == nil || len(paylines[i]) != conf.ColsSize {
			panic(fmt.Sprintf("invalid pay lines or row size at %d is not %d", i, conf.ColsSize))
		}
		for j := range paylines[i] {
			if paylines[i][j] < 0 || paylines[i][j] >= conf.RowsSize {
				panic(fmt.Sprintf("invalid pay lines value, must be positive and less than %d", conf.RowsSize))
			}
		}
	}

	if paytable == nil || len(paytable) != conf.ColsSize+1 {
		panic(fmt.Sprintf("invalid pay table or paytable size (paytable size = number of columns + 1)"))
	}

	for i := range paytable {
		if paytable[i] == nil || len(paytable[i]) != len(conf.Symbols) {
			panic(fmt.Sprintf("invalid pay table at %d (size must equals number of symbols)", i))
		}
	}
	return &Model{
		conf:     conf,
		paylines: paylines,
		paytable: paytable,
	}
}

func (m *Model) Win(machine *goslot.SlotMachine) int {
	win := 0
	for _, payLine := range m.paylines {
		// lấy line tương ứng với payline này
		line := make([]int, m.conf.ColsSize)
		for i := 0; i < m.conf.ColsSize; i++ {
			line[i] = machine.Reels()[i][(machine.Stops()[i]+payLine[i])%m.conf.ReelSize]
		}

		// lấy biểu tượng đầu tiên (từ trái qua phải) khác WILD
		symbol := line[0]
		for i := 0; i < len(line); i++ {
			if m.conf.Types[symbol] != goslot.WILD {
				break
			}
			symbol = line[i]
		}

		// thay tất cả các WILD thành biểu tượng tìm được
		for i := 0; i < len(line); i++ {
			if m.conf.Types[line[i]] == goslot.WILD {
				line[i] = symbol
			}
		}

		// đếm từ trái qua phải xem có bao nhiêu symbol liên tiếp
		counter := 0
		for i := 0; i < len(line); i++ {
			if line[i] == symbol {
				counter++
			} else {
				break
			}
		}
		// tính tiền số lượng symbol đó
		win += m.paytable[counter][symbol]
	}
	return win
}

func (m *Model) Jackpot(machine *goslot.SlotMachine) bool {
	Loop: for _, payLine := range m.paylines {
		for i := 0; i < m.conf.ColsSize; i++ {
			if m.conf.Types[machine.Reels()[i][(machine.Stops()[i]+payLine[i])%m.conf.ReelSize]] != goslot.WILD {
				continue Loop
			}
		}
		return true
	}
	return false
}

func (m *Model) IsInvalid(machine *goslot.SlotMachine) bool {
	for i := 0; i < m.conf.ColsSize; i++ {
		counter := make([]int, len(m.conf.Symbols))
		for j := 0; j < m.conf.ReelSize; j++ {
			counter[machine.Reels()[i][j]]++
		}
		for j, count := range counter {
			if count == 0 {
				return true
			}
			if m.conf.Types[j] == goslot.WILD && count < 2 {
				return true
			}
		}
	}
	return false
}

func (m *Model) Scatters(machine *goslot.SlotMachine) int {
	return 0
}

func (m *Model) Bonus(machine *goslot.SlotMachine) int {
	return 0
}

// RTP, Jackpot
func (m *Model) Result(machine *goslot.SlotMachine) []float64 {
	result := make([]float64, 2)
	result[0] += float64(m.Win(machine)) / float64(len(m.paylines))
	if m.Jackpot(machine) {
		result[1] += 1
	}
	return result
}

var paylines = [][]int{
	{0, 0, 0},
	{1, 1, 1},
	{2, 2, 2},
	{0, 1, 0},
	{2, 1, 2},
	{1, 0, 1},
	{1, 2, 1},
	{0, 2, 0},
	{2, 0, 2},
	{0, 1, 2},
	{2, 1, 0},
	{0, 0, 1},
	{1, 1, 2},
	{1, 1, 0},
	{2, 2, 1},
	{1, 0, 0},
	{2, 1, 1},
	{0, 1, 1},
	{1, 2, 2},
	{0, 2, 1},
}

var paytable = [][]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{1000, 300, 100, 50, 20, 10, 5, 0},
}

var conf = &goslot.Conf{
	ColsSize:                3,
	ReelSize:                40,
	RowsSize:                3,
	NumberOfNodes:           5,
	LocalPopulationSize:     10,
	LocalOptimizationEpochs: 20,
	NumberOfLifeCircle:      20,
	Targets:                 []float64{0.8, 0.00001},
	Symbols:                 []string{"A", "B", "C", "D", "E", "F", "G", "WILD"},
	Types: []goslot.SymbolType{goslot.REGULAR, goslot.REGULAR,
		goslot.REGULAR, goslot.REGULAR, goslot.REGULAR, goslot.REGULAR,
		goslot.REGULAR, goslot.WILD},
	OutputFile:                 fmt.Sprintf("model-classic-%s.txt", now()),
}

func Start() {
	conf.Validate()
	model := NewModel(conf, paylines, paytable)
	gen := goslot.NewGenerator(conf, model)
	gen.Start()
	data := []byte(goslot.ChromosomeString(gen.GetBestChromosome(), conf.Symbols))
	if err := gen.WriteFile(data); err != nil {
		panic(err)
	}
}

func Gen() {
	rand.Seed(time.Now().UnixNano())
	conf.Validate()
	model := NewModel(conf, paylines, paytable)
	for {
		machine := goslot.NewMachine(conf, model)
		ga := goslot.NewGeneticAlgorithm(conf)
		ga.RandomReels(machine)
		m := machine.Compute(ga.GetRandomChromosome().Reels())
		var sum float64
		var jackpot float64
		var counter = 0
		var max float64
		list := []int64{}
		for key, value := range m {
			if value[0] > 1 {
				counter++
			}
			if value[0] > 1 && rand.Intn(30) != 0 {
				continue
			}
			if value[1] > 0 && rand.Intn(100) != 0 {
				continue
			}
			if value[0] > max {
				max = value[0]
			}
			sum += value[0]
			jackpot += value[1]
			list = append(list, key)
		}
		sum = sum / float64(len(m))
		jackpot = jackpot / float64(len(m))
		eps1 := math.Abs(conf.Targets[0] - sum)
		eps2 := math.Abs(conf.Targets[1] - jackpot)
		if eps1 <= 0.01 {
			println(goslot.ChromosomeString(ga.GetRandomChromosome(), conf.Symbols))
			println(fmt.Sprintf("%f", machine.Evaluate(ga.GetRandomChromosome().Reels())))
			println(fmt.Sprintf("tỉ lệ ăn: %f", sum))
			println(fmt.Sprintf("tỉ lệ ăn jackpot: %f", jackpot))
			println(fmt.Sprintf("số case ăn lớn hơn 1: %d, số case tổng: %d", counter, len(m)))
			println(fmt.Sprintf("ăn lớn nhất: %f",  max))
			println(fmt.Sprintf("eps: %f %f", eps1, eps2))
			println(fmt.Sprintf("size: %d", len(list)))
			if  eps2 <= 0.00001 && len(list) > 30000 {
				break
			}
		}
	}

}

func now() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d-%02d-%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}
