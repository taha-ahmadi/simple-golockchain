package golockchain

import (
	"crypto/sha256"
	"fmt"
)

func GenerateMask(zeros int) []byte{
	full, half := zeros/2, zeros%2
	var mask []byte
	for i := 0; i < full; i++ {
		mask = append(mask, 0)
	}

	if half > 0 {
		mask = append(mask, 0xf)
	}

	return mask
}

// ValidDifficulty checks the difficulty of a given hash
func ValidDifficulty(mask, hash []byte) bool{
	for i := range mask {
		if hash[i] > mask[i]{
			return false
		}
	}
	return true
}

// GenerateHash takes data and retruns the hash of data
func GenerateHash(data ...interface{}) []byte {
	hasher := sha256.New()

	fmt.Fprint(hasher, data...)

	return hasher.Sum(nil)
}

// DifficultHash is a simple mining algorithm
func DifficultHash(mask []byte, data ...interface{}) ([]byte, int32) {
	ln := len(data)
	data = append(data, nil)
	var i int32
	for { 
		data[ln] = i
		hash := GenerateHash(data...)
		if ValidDifficulty(mask, hash){
			return hash, i
		}
		i++
	}
}
