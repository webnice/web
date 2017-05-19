package param

import (
	"testing"
)

const KEY = "foo"
const VAL = "bar"

var res interface{}

func BenchmarkGet(b *testing.B) {
	var prm Interface
	var str string

	prm = New()
	prm.Set(KEY, VAL)
	for n := 0; n < b.N; n++ {
		if x := prm.Get(KEY); x != "" {
			str = x
		}
	}
	res = str
}

func BenchmarkSet(b *testing.B) {
	var prm Interface

	prm = New()
	for n := 0; n < b.N; n++ {
		prm.Set(KEY, VAL)
	}
}

func BenchmarkHas(b *testing.B) {
	var prm Interface
	var has bool

	prm = New()
	for n := 0; n < b.N; n++ {
		res = prm.Has(KEY)
	}
	res = has
}

func BenchmarkDel(b *testing.B) {
	var prm Interface

	prm = New()
	prm.Set(KEY, VAL)
	for n := 0; n < b.N; n++ {
		prm.Del(KEY)
	}
}
