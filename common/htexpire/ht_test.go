package htexpire

import (
	"log"
	"testing"
)

// TestNew simple test hash table
func TestNew(t *testing.T) {
	ht := New()
	var list = []string{"a", "b", "c", "x", "y", "z"}
	for _, n := range list {
		if err := ht.Set(n); err != nil {
			log.Println(err)
		}
	}
	if err := ht.Remove("y"); err != nil {
		log.Println(err)
	}
	log.Println(ht.Get("y"))
}
