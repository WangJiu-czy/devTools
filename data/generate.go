package data

import (
	"math/rand"
	"time"
)

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []rune("123456789qwertyuiopasfghjklzxcvbnm")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i, _ := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
