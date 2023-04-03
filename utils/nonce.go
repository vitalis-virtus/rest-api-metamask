package utils

import (
	"math/big"
	"sync"
)

var (
	max  *big.Int
	once sync.Once
)

func GetOnce() (string, error) {
	once.Do()
}
