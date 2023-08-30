package slice

import (
	"log"
	"testing"
)

func TestDivideBigSlice(t *testing.T) {
	result := DivideBigSlice([]string{}, 3)
	log.Println(result)
}
