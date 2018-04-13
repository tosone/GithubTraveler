package htexpire

import (
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	ht := New()
	var list = []string{"a", "b", "c", "x", "y", "z"}
	for _, n := range list {
		ht.Set(n)
	}
	ht.Remove("y")
	log.Println(ht.Get("y"))
}
