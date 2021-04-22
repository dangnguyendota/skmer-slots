package classic

import (
	"fmt"
	"github.com/dangnguyendota/goslot"
	"io/ioutil"
	"os"
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

func (m *Model) Scatters(machine *goslot.SlotMachine) int {
	return 0
}

func (m *Model) Bonus(machine *goslot.SlotMachine) int {
	return 0
}

func (m *Model) Result(machine *goslot.SlotMachine) []float64 {
	result := make([]float64, 2)
	result[0] += float64(m.Win(machine))
	Loop: for i := 0; i < m.conf.ColsSize; i++ {
		for j := 0; j < m.conf.RowsSize; j++ {
			if m.conf.Types[machine.Reels()[i][j]] != goslot.WILD {
				continue Loop
			}
		}
		result[1] += 1
		break
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

func Start() {
	conf := &goslot.Conf{
		ColsSize:                3,
		ReelSize:                40,
		RowsSize:                3,
		NumberOfNodes:           50,
		LocalPopulationSize:     37,
		LocalOptimizationEpochs: 100,
		NumberOfLifeCircle:      67,
		Targets:                 []float64{0.6, 0.000015},
		Symbols:                 []string{"A", "B", "C", "D", "E", "F", "G", "WILD"},
		Types: []goslot.SymbolType{goslot.REGULAR, goslot.REGULAR,
			goslot.REGULAR, goslot.REGULAR, goslot.REGULAR, goslot.REGULAR,
			goslot.REGULAR, goslot.WILD},
		OutputFile: fmt.Sprintf("model-classic-%s.txt", now()),
	}
	conf.Validate()
	model := NewModel(conf, paylines, paytable)
	gen := goslot.NewGenerator(conf, model)
	gen.Start()
	data := []byte(goslot.ChromosomeString(gen.GetBestChromosome(), conf.Symbols))
	if err := ioutil.WriteFile(conf.OutputFile, data, os.ModeAppend); err != nil {
		panic(err)
	}
}

func now() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}
