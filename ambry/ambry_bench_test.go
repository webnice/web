package ambry

import "testing"

const (
	KEY = "foo"
	VAL = "bar"
)

var res interface{}

func BenchmarkGet(b *testing.B) {
	var str interface{}
	var prm = New()

	prm.Set(KEY, VAL)
	for n := 0; n < b.N; n++ {
		if x := prm.Get(KEY); x != "" {
			str = x
		}
	}
	res = str
}

func BenchmarkSet(b *testing.B) {
	var prm = New()

	for n := 0; n < b.N; n++ {
		prm.Set(KEY, VAL)
	}
}

func BenchmarkHas(b *testing.B) {
	var has bool
	var prm = New()

	for n := 0; n < b.N; n++ {
		res = prm.Has(KEY)
	}
	res = has
}

func BenchmarkDel(b *testing.B) {
	var prm = New()

	prm.Set(KEY, VAL)
	for n := 0; n < b.N; n++ {
		prm.Del(KEY)
	}
}
