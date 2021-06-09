package minipoker

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
)

type Rank int

const (
	Two Rank = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

type Suit int

const (
	Spade Suit = iota
	Club
	Diamond
	Heart
)

type MiniPokerReels []int

func reels(r []int) MiniPokerReels {
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return r
}

func rank(card int) Rank {
	return Rank(card / 4)
}

func suit(card int) Suit {
	return Suit(card % 4)
}

// true nếu là sảnh rồng
func (r MiniPokerReels) isDragonHead() bool {
	return r.isStraightFlush() && rank(r[0]) == Ten
}

// true nếu là dây đồng chất
func (r MiniPokerReels) isStraightFlush() bool {
	return r.isFlush() && r.isSequence()
}

// true nếu đồng chất
func (r MiniPokerReels) isFlush() bool {
	for i := 0; i < 4; i++ {
		if suit(r[i]) != suit(r[i+1]) {
			return false
		}
	}
	return true
}

// true nếu là 1 dây
func (r MiniPokerReels) isSequence() bool {
	for i := 0; i < 4; i++ {
		if rank(r[i])+1 != rank(r[i+1]) {
			return false
		}
	}
	return true
}

// true nếu có tứ quý
func (r MiniPokerReels) isQuads() bool {
	return rank(r[0]) == rank(r[3]) || rank(r[1]) == rank(r[4])
}

// true nếu 1 tam và 1 đôi
func (r MiniPokerReels) isTripsAndDubs() bool {
	return (rank(r[0]) == rank(r[2]) && rank(r[3]) == rank(r[4])) ||
		(rank(r[0]) == rank(r[1]) && rank(r[2]) == rank(r[4]))
}

// true nếu có 1 bộ tam
func (r MiniPokerReels) isTrips() bool {
	return rank(r[0]) == rank(r[2]) || rank(r[1]) == rank(r[3]) || rank(r[2]) == rank(r[4])
}

// true nếu có 2 đôi
func (r MiniPokerReels) isDoubleDubs() bool {
	return (rank(r[0]) == rank(r[1]) && rank(r[2]) == rank(r[3])) ||
		(rank(r[1]) == rank(r[2]) && rank(r[3]) == rank(r[4])) ||
		(rank(r[0]) == rank(r[1]) && rank(r[3]) == rank(r[4]))
}

// true nếu là đôi lớn hơn J
func (r MiniPokerReels) isJDubs() bool {
	if d := r.getDubs(); d != -1 && d >= Jack {
		return true
	}
	return false
}

// true nếu đôi nhỏ hơn 10
func (r MiniPokerReels) isTenDubs() bool {
	if d := r.getDubs(); d != - 1 && d <= Ten {
		return true
	}
	return false
}

func (r MiniPokerReels) getDubs() Rank {
	for i := 0; i < 4; i++ {
		if rank(r[i]) == rank(r[i+1]) {
			return rank(r[i])
		}
	}
	return -1
}

// true nếu là 1 dây
func isSequence(reels MiniPokerReels) bool {
	for i := 0; i < 4; i++ {
		if rank(reels[i])+1 != rank(reels[i+1]) {
			return false
		}
	}
	return true
}

func LoadJsonConf(config interface{}, path string) {
	flag.Parse()
	// load config
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	if err := json.Unmarshal(file, config); err != nil {
		log.Panic(err)
	}
}

type MiniPokerConf struct {
	StraightFlush    int     `json:"straight_flush"`
	Quads            int     `json:"quads"`
	TripsAndDubs     int     `json:"trips_and_dubs"`
	Flush            int     `json:"flush"`
	Sequence         int     `json:"sequence"`
	Trips            int     `json:"trips"`
	DoubleDubs       int     `json:"double_dubs"`
	JDubs            float64 `json:"j_dubs"`
	TenDubs          float64 `json:"ten_dubs"`
	JackpotHouseEdge float64 `json:"jackpot_house_edge"`
}

func startMinipoker() {
	flag.Parse()
	var conf *MiniPokerConf
	conf = &MiniPokerConf{}
	LoadJsonConf(conf, "./conf.json")

	straightFlushCount := 0
	flushCount := 0
	sequenceCount := 0
	quadsCount := 0
	tripsAndDubsCount := 0
	tripsCount := 0
	doubleDubsCount := 0
	jDubsCount := 0
	tenDubsCount := 0
	var total float64
	dragonHeadCount := 0
	for card1 := 0; card1 < 48; card1++ {
		for card2 := card1 + 1; card2 < 49; card2++ {
			for card3 := card2 + 1; card3 < 50; card3++ {
				for card4 := card3 + 1; card4 < 51; card4++ {
					for card5 := card4 + 1; card5 < 52; card5++ {
						total++
						r := reels([]int{card1, card2, card3, card4, card5})
						if r.isDragonHead() {
							dragonHeadCount++
						}
						if r.isStraightFlush() {
							straightFlushCount++
						} else if r.isQuads() {
							quadsCount++
						} else if r.isTripsAndDubs() {
							tripsAndDubsCount++
						} else if r.isFlush() {
							flushCount++
						} else if r.isSequence() {
							sequenceCount++
						} else if r.isTrips() {
							tripsCount++
						} else if r.isDoubleDubs() {
							doubleDubsCount++
						} else if r.isJDubs() {
							jDubsCount++
						} else if r.isTenDubs() {
							tenDubsCount++
						}
					}
				}
			}
		}
	}
	println(fmt.Sprintf("Jackpot: %d, tỉ lệ ăn: %f%%", dragonHeadCount, float64(dragonHeadCount*100)/total))
	println(fmt.Sprintf("Thùng phá sảnh: %d, tỉ lệ ăn: %f%%", straightFlushCount, float64(straightFlushCount*100)/total))
	println(fmt.Sprintf("Tứ quý: %d, tỉ lệ ăn: %f%%", quadsCount, float64(quadsCount*100)/total))
	println(fmt.Sprintf("1 Tam và 1 Đôi: %d, tỉ lệ ăn: %f%%", tripsAndDubsCount, float64(tripsAndDubsCount*100)/total))
	println(fmt.Sprintf("Đồng chất: %d, tỉ lệ ăn: %f%%", flushCount, float64(flushCount*100)/total))
	println(fmt.Sprintf("Dây 5: %d, tỉ lệ ăn: %f%%", sequenceCount, float64(sequenceCount*100)/total))
	println(fmt.Sprintf("1 Tam: %d, tỉ lệ ăn: %f%%", tripsCount, float64(tripsCount*100)/total))
	println(fmt.Sprintf("2 Đôi: %d, tỉ lệ ăn: %f%%", doubleDubsCount, float64(doubleDubsCount*100)/total))
	println(fmt.Sprintf("1 Đôi >= J: %d, tỉ lệ ăn: %f%%", jDubsCount, float64(jDubsCount*100)/total))
	println(fmt.Sprintf("1 Đôi <= 10: %d, tỉ lệ ăn: %f%%", tenDubsCount, float64(tenDubsCount*100)/total))
	println(fmt.Sprintf("Tổng trường hợp: %d", int64(total)))
	println(fmt.Sprintf("Xác xuất ăn: %f%%", (
		float64(straightFlushCount*(conf.StraightFlush)+
			quadsCount*(conf.Quads)+
			tripsAndDubsCount*(conf.TripsAndDubs)+
			flushCount*(conf.Flush)+
			sequenceCount*(conf.Sequence)+
			tripsCount*(conf.Trips)+
			doubleDubsCount*(conf.DoubleDubs))+
			float64(jDubsCount)*(conf.JDubs)+
			float64(tenDubsCount)*(conf.TenDubs))*100/(total-total*conf.JackpotHouseEdge)))
}
