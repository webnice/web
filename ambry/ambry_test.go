package ambry

import (
	"testing"
)

var (
	k  = `A69F2F4D-0B8F-43D3-8E1B-1EA90766FEB5`
	v1 = `ac5110c701777822d999d0be2fd44a945dde58b46398d2c12ada878c5271e167770305532a6d1aa05cf16c575edf3ae99fd3f55db814369062106c81a9b443f3`
	v2 = `205c23c53fc4cf5ea6bf7a352fc4c7cf91ac001dcbfbebbe97b0c59cdf053096ceaf8e41e392c7953b379d0c6f6094c56471aa253f2b3102e4f904cc8d9ad156`
	v3 = `5f4f7401e818ccc9af6d6c6c19eff4ce0696925bfb2a839c7bea211dda029816f48bc2417a9420b401450b6d6092dd4273f2a2d6a9364a88566b175d2de72576`
)

func TestSet(t *testing.T) {
	var prm = New()

	if prm.Has(k) {
		t.Errorf("key %s exist", k)
	}
	prm.Set(k, v1)
	if !prm.Has(k) {
		t.Errorf("key is %s exist", k)
	}
	if prm.Get(k) != v1 {
		t.Errorf("value is not equal")
	}
}

func TestGet(t *testing.T) {
	var prm = New()

	if prm.Get(k) != nil {
		t.Errorf("Get return incorrect value")
	}
	prm.Set(k, v1)
	if prm.Get(k) != v1 {
		t.Errorf("value is not equal")
	}
}

func TestAdd(t *testing.T) {
	var prm Interface
	var all []interface{}

	prm = New()

	all = prm.GetAll(k)
	if len(all) != 0 {
		t.Errorf("Error in GetAll")
	}
	prm.Add(k, v1)
	prm.Add(k, v2)
	prm.Add(k, v3)
	if prm.Get(k) != v1 {
		t.Errorf("value is not equal")
	}
	all = prm.GetAll(k)
	if len(all) != 3 {
		t.Errorf("incorrect length")
	}
	if all[0] != v1 || all[1] != v2 || all[2] != v3 {
		t.Errorf("incorrect value")
	}
}

func TestKeys(t *testing.T) {
	var prm Interface
	var arr []interface{}
	var oks []bool

	prm = New()
	arr = prm.Keys()
	if len(arr) != 0 {
		t.Errorf("Error in Keys")
	}

	prm.Add("k1", v1)
	prm.Add("k2", v2)
	prm.Set("k3", v3)
	arr = prm.Keys()
	if len(arr) != 3 {
		t.Errorf("Error in Keys")
	}
	for i := range arr {
		if arr[i] == "k1" {
			oks = append(oks, true)
		}
		if arr[i] == "k2" {
			oks = append(oks, true)
		}
		if arr[i] == "k3" {
			oks = append(oks, true)
		}
	}
	if len(oks) != 3 {
		t.Errorf("Error in Keys")
	}
}

func TestDel(t *testing.T) {
	var prm = New()

	if prm.Del(k) != nil {
		t.Errorf("Error in Del")
	}
	prm.Set(k, v1)
	if prm.Get(k) != v1 {
		t.Errorf("Error in Get")
	}
	if prm.Del(k) != v1 {
		t.Errorf("Error in Del")
	}
	if prm.Get(k) != nil {
		t.Errorf("Error in Get")
	}
}
