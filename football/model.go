package football

import (
	"../../goslot"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math"
	"math/rand"
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
	for i := 0; i < m.conf.ColsSize; i++ {
		for j := 0; j < m.conf.RowsSize; j++ {
			if m.conf.Types[machine.Reels()[i][(machine.Stops()[i]+j)%len(machine.Reels()[i])]] == goslot.BONUS {
				counter++
			}
		}
	}
	return counter
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

// RTP, Jackpot, 3 Free spins, 4 Free Spins, 5 Free spins
func (m *Model) Result(machine *goslot.SlotMachine) []float64 {
	result := make([]float64, 3)
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
		result[2] += 10
	case 4:
		result[2] += 15
	case 5:
		result[2] += 25
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
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{15, 10, 10, 7, 7, 6, 5, 5, 0},
	{30, 20, 20, 15, 15, 12, 10, 10, 0},
	{75, 50, 50, 25, 25, 20, 15, 15, 0},
}

var conf = &goslot.Conf{
	ColsSize:                5,
	ReelSize:                20,
	RowsSize:                3,
	NumberOfNodes:           5,
	LocalPopulationSize:     10,
	LocalOptimizationEpochs: 20,
	NumberOfLifeCircle:      11,
	Targets:                 []float64{0.9, 0.00001, 0.02, 0.01, 0.005},
	Symbols:                 []string{"A", "B", "C", "D", "E", "F", "G", "H", "WILD"},
	Types: []goslot.SymbolType{
		goslot.REGULAR, goslot.REGULAR, goslot.REGULAR,
		goslot.REGULAR, goslot.REGULAR, goslot.REGULAR,
		goslot.REGULAR, goslot.REGULAR, goslot.WILD},
	OutputFile: fmt.Sprintf("model-football-%s.txt", now()),
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

type Result struct {
	Id       uuid.UUID `json:"id"`
	RTP      float64   `json:"rtp"`
	Jackpot  float64   `json:"jackpot"`
	Bound    float64   `json:"bound"`
	ReelSize int       `json:"reel_size"`
	Code     string    `json:"code"`
	List     []int64   `json:"list"`
	Blocked  []int64   `json:"blocked"`
}

func Gen() {
	rand.Seed(time.Now().UnixNano())
	conf.Validate()
	model := NewModel(conf, paylines, paytable)
	for {
		var bound float64 = 10
		machine := goslot.NewMachine(conf, model)
		ga := goslot.NewGeneticAlgorithm(conf)
		ga.RandomReels(machine, false)
		m := machine.Compute(ga.GetRandomChromosome().Reels())
		var rtp float64
		var jackpot float64
		var freespins float64

		var counter = 0
		var zeroCounter = 0
		var oneCounter = 0
		var max float64
		list := []int64{}
		blocked := []int64{}
		for key, value := range m {
			if value[0] > bound {
				continue
			}
			if value[1] > 0 && rand.Intn(100) != 0 {
				if value[0] <= bound {
					blocked = append(blocked, key)
				}
				continue
			}
			if value[0] > max {
				max = value[0]
			}
			if value[0] > 1 {
				oneCounter++
			}
			if value[0] == 0 {
				zeroCounter++
			}
			rtp += value[0]
			jackpot += value[1]
			freespins += value[2]
			counter++
		}
		rtp = rtp / float64(counter)
		jackpot = jackpot / float64(counter)
		freespins = freespins / float64(counter)
		if jackpot == 0 {
			continue
		}
		//if freespins == 0 {
		//	continue
		//}
		eps1 := math.Abs(conf.Targets[0] - rtp)
		eps2 := math.Abs(conf.Targets[1] - jackpot)

		println(goslot.ChromosomeString(ga.GetRandomChromosome(), conf.Symbols))
		//println(fmt.Sprintf("%f", machine.Evaluate(ga.GetRandomChromosome().Stops())))
		println(fmt.Sprintf("tỉ lệ ăn (RTP): %f", rtp))
		println(fmt.Sprintf("tỉ lệ ăn jackpot (Jackpot): %f", jackpot))
		println(fmt.Sprintf("tỉ lệ ăn free spins: %f", freespins))
		println(fmt.Sprintf("số case tổng: %d", len(m)))
		println(fmt.Sprintf("số case lấy ra: %d ", counter))
		println(fmt.Sprintf("ăn lớn nhất: %f", max))
		println(fmt.Sprintf("eps: %f %f", eps1, eps2))
		println(fmt.Sprintf("size: %d", len(list)))
		println(fmt.Sprintf("blocked: %d", len(blocked)))
		println(fmt.Sprintf("zero: %d", zeroCounter))
		println(fmt.Sprintf("one: %d", oneCounter))
		println(rtp <= 0.9 && jackpot <= 0.0001)
		if rtp <= 0.9 && jackpot <= 0.0001 {
			result := &Result{
				Id:      uuid.New(),
				RTP:     rtp,
				Jackpot: jackpot,
				Bound:   bound,
				Code:    ga.GetRandomChromosome().Code(conf.Symbols),
				List:    list,
				Blocked: blocked,
			}
			s, err := json.Marshal(result)
			if err != nil {
				panic(err)
			}
			filename := "/home/dangnguyendota/Desktop/backup/code/skmer/skmer-server/skmer-slots/result/football-" + result.Id.String() + ".json"
			if err := WriteFile(filename, s); err != nil {
				panic(err)
			}
			println("write to file")
		}
	}
}

func WriteFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func now() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d-%02d-%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}
