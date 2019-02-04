package skip

import (
	"testing"
	"strconv"

	"github.com/zeyadyasser/autocom/complete"
	"github.com/zeyadyasser/autocom/complete/tst"
)

const short_key = "Short Key"
const long_key = "Very Very Very Very Very Very Very Very Very Long Movie Name"

func benchSet1ShortKey(b *testing.B, E *SkipEngine) {
	b.ReportAllocs()
	end := 100000
	for n := 0; n < b.N; n++ {
		E.Set(strconv.Itoa(end) + short_key, nil)
		end++
	}
}

func benchTopN10ShortKey(b *testing.B, E *SkipEngine) {
	res := make(complete.Map)
	res["boo"] = "foo"
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tmp, _ := E.TopN(short_key, 10)
		res = tmp
	}
}

func BenchmarkTSTSet1ShortKey(b *testing.B) {
	opts := Options{}
	factory := tst.NewTSTCompleter
	E := NewSkipEngine(opts, factory)
	benchSet1ShortKey(b, E)
}


func BenchmarkTSTSet1LongKey(b *testing.B) {
	opts := Options{}
	factory := tst.NewTSTCompleter
	E := NewSkipEngine(opts, factory)
	benchSet1ShortKey(b, E)
}

func BenchmarkTSTTopN10ShortKey(b *testing.B) {
	opts := Options{}
	factory := tst.NewTSTCompleter
	E := NewSkipEngine(opts, factory)
	for n := 0; n < 1000; n++ {
		E.Set(strconv.Itoa(n), nil)
		E.Set(short_key + strconv.Itoa(n), nil)
	}
	benchTopN10ShortKey(b, E)
}
/*
func BenchmarkTrieSet1ShortKey(b *testing.B) {
	opts := Options{}
	E := NewSkipEngine(opts, nil)
	benchSet1ShortKey(b, E)
}
*/
