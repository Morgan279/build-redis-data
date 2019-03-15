package tools

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	r *rand.Rand
)

type Rand struct {
	r *rand.Rand
}

func NewRand() *Rand {
	r := &Rand{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return r
}

// 包括传入的最大值, int值最小为1
func (r *Rand) RandInt(maxInt int) int {
	return r.r.Intn(maxInt) + 1
}

func (r *Rand) RandFloat(maxFloat int) float64 {
	return float64(maxFloat) * r.r.Float64()
}

func (r *Rand) RandRangeKey(prefix string, maxInt int) string {
	suffix := r.r.Intn(maxInt)
	return fmt.Sprintf("%s:%d", prefix, suffix)
}

func (r *Rand) RandString(maxLen int) string {
	length := r.RandInt(maxLen)
	buff := make([]byte, length)
	for i := 0; i < length; i++ {
		buff[i] = byte(r.r.Intn(26) + 'a')
	}
	return string(buff)
}

func (r *Rand) RandRuneString(maxLen int) string {
	buff := make([]rune, 0)
	length := r.RandInt(maxLen)
	for {
		if length < 0 {
			break
		}
		buff = append(buff, rune(r.r.Intn(0xFF)))
		length--
	}
	return string(buff)
}
