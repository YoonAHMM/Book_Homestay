package Randx

import (
	"math/rand"
	"time"
)

func GenerateRandomFixedRange(max, len int64) (int64,int64) {
    // 初始化随机数生成器
    rand.Seed(time.Now().UnixNano())

	var right,left int64 
	
	right=0
	for right < len {
		right = rand.Int63n(max)
	}

	left = right + 1 - len

	return left,right
} 