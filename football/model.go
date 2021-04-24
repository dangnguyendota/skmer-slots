package football

import (
	"fmt"
	"github.com/dangnguyendota/goslot"
	"time"
)

type Model struct {
	conf     *goslot.Conf
	paylines [][]int
	paytable [][]int
	wild     []int
}

func NewModel(conf *goslot.Conf, paylines [][]int, paytable [][]int, wild []int) *Model {
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
		wild:     wild,
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
		wild := 0
		for i := 0; i < len(line); i++ {
			if m.conf.Types[line[i]] == goslot.WILD {
				line[i] = symbol
				wild++
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
		win += m.wild[wild]
	}
	return win
}

func (m *Model) Jackpot(machine *goslot.SlotMachine) bool {
Loop:
	for _, payLine := range m.paylines {
		for i := 0; i < m.conf.ColsSize; i++ {
			if m.conf.Types[machine.Reels()[i][(machine.Stops()[i]+payLine[i])%m.conf.ReelSize]] != goslot.WILD {
				continue Loop
			}
		}
		return true
	}
	return false
}

func (m *Model) Scatters(machine *goslot.SlotMachine) int {
	return 0
}

func (m *Model) Bonus(machine *goslot.SlotMachine) int {
	counter := 0
	for i := 0; i < len(machine.Reels()); i++ {
		for j := 0; j < m.conf.RowsSize; j++ {
			if m.conf.Types[machine.Reels()[i][(machine.Stops()[i]+j)%len(machine.Reels()[i])]] == goslot.BONUS {
				counter++
			}
		}
	}
	return counter
}

func (m *Model) IsInvalid(machine *goslot.SlotMachine) bool {
Loop1:
	for i := 0; i < m.conf.ColsSize; i++ {
		for j := 0; j < m.conf.ReelSize; j++ {
			if m.conf.Types[machine.Reels()[i][j]] == goslot.WILD {
				continue Loop1
			}
		}
		return true
	}
Loop2:
	for i := 0; i < m.conf.ColsSize; i++ {
		for j := 0; j < m.conf.ReelSize; j++ {
			if m.conf.Types[machine.Reels()[i][j]] == goslot.BONUS {
				continue Loop2
			}
		}
		return true
	}
	return false
}

// RTP, Jackpot, 3 Free spins, 4 Free Spins, 5 Free spins
func (m *Model) Result(machine *goslot.SlotMachine) []float64 {
	if m.IsInvalid(machine) {
		return []float64{goslot.InvalidReelsPenalty, goslot.InvalidReelsPenalty,
			goslot.InvalidReelsPenalty, goslot.InvalidReelsPenalty, goslot.InvalidReelsPenalty}
	}
	result := make([]float64, 5)
	result[0] += float64(m.Win(machine)) / float64(len(m.paylines))
	if m.Jackpot(machine) {
		result[1] += 1
	}
	bonus := m.Bonus(machine)
	switch bonus {
	case 0:
	case 1:
	case 2:
	case 3:
		result[2]++
	case 4:
		result[3]++
	case 5:
		result[4]++
	default:
		// nếu nhiều hơn 5 bonus trên 1 màn hình thì penalty
		result[0] += goslot.InvalidReelsPenalty
	}
	return result
}

var paylines = [][]int{
	{0, 0, 0, 0, 0},
	{1, 1, 1, 1, 1},
	{2, 2, 2, 2, 2},
	{1, 2, 2, 2, 1},
	{1, 0, 1, 2, 1},
	{0, 0, 2, 0, 0},
	{2, 2, 0, 2, 2},
	{0, 1, 2, 1, 0},
	{0, 1, 1, 1, 0},
	{0, 1, 0, 1, 0},
	{0, 2, 2, 2, 0},
	{2, 0, 0, 0, 2},
	{2, 1, 0, 1, 2},
	{0, 0, 1, 2, 2},
	{2, 1, 2, 1, 2},
	{1, 2, 0, 2, 1},
	{2, 1, 1, 1, 2},
	{1, 2, 1, 0, 1},
	{1, 1, 2, 1, 1},
	{0, 2, 0, 2, 0},
	{2, 0, 2, 0, 2},
	{1, 0, 0, 0, 1},
	{2, 2, 1, 0, 0},
	{1, 1, 0, 1, 1},
	{1, 0, 2, 0, 1},
}

var paytable = [][]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{15, 10, 10, 7, 7, 6, 5, 5, 0, 0},
	{30, 20, 20, 15, 15, 12, 10, 10, 0, 0},
	{75, 50, 50, 25, 25, 20, 15, 15, 0, 0},
}

// tiền thưởng khi ăn [0, 1, 2, 3, 4, 5] WILD trên 1 pay line
var wild = []int{0, 0, 0, 25, 40, 0}

func Start() {
	conf := &goslot.Conf{
		ColsSize:                5,
		ReelSize:                60,
		RowsSize:                3,
		NumberOfNodes:           5,
		LocalPopulationSize:     3,
		LocalOptimizationEpochs: 5,
		NumberOfLifeCircle:      11,
		Targets:                 []float64{0.7, 0.0001, 0.02, 0.01, 0.005},
		Symbols:                 []string{"A", "B", "C", "D", "E", "F", "G", "H", "WILD", "FREE SPIN"},
		Types: []goslot.SymbolType{
			goslot.REGULAR, goslot.REGULAR, goslot.REGULAR,
			goslot.REGULAR, goslot.REGULAR, goslot.REGULAR,
			goslot.REGULAR, goslot.REGULAR, goslot.WILD, goslot.BONUS},
		OutputFile: fmt.Sprintf("model-football-%s.txt", now()),
	}
	conf.Validate()
	model := NewModel(conf, paylines, paytable, wild)
	gen := goslot.NewGenerator(conf, model)
	gen.Start()
	data := []byte(goslot.ChromosomeString(gen.GetBestChromosome(), conf.Symbols))
	if err := gen.WriteFile(data); err != nil {
		panic(err)
	}
}

func now() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d-%02d-%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}
