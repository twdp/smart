package smart

import (
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/lists/arraylist"
	"testing"
)

func TestNewProcess(t *testing.T) {
	l := arraylist.New("1", "2", "3")
	l.Add("dd")
	a, _ := l.ToJSON()
	aa, _ := json.Marshal(l.Values()[:l.Size()])
	fmt.Println(string(a))
	fmt.Println(string(aa))
	aaa, _ := json.Marshal(l.Values())
	fmt.Println(string(aaa))

	aaaa, _ := json.Marshal(l)
	fmt.Println(string(aaaa))
}