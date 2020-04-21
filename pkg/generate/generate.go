package generate

import (
	"crypto/rand"
	"math/big"
)

const (
	alphaNum = "0123456789abcdefghijklmnopqrstuvwxyz"
)

func GetRandomString(n int) (string, error) {
	buffer := make([]byte, n)
	max := big.NewInt(int64(len(alphaNum)))

	for i := 0; i < n; i++ {
		index, err := randomInt(max)
		if err != nil {
			return "", err
		}
		buffer[i] = alphaNum[index]
	}

	return string(buffer), nil
}

func randomInt(max *big.Int) (int, error) {
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return int(r.Int64()), nil
}
