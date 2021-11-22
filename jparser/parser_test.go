package jparser

import (
	"reflect"
	"testing"
)

var arr = []byte(`[{"jsonrpc":"2.0","method":"sum","params":[1],"id":1}, {"jsonrpc":"2.0","method":"sum","params":[1, 2],"id":2}, {"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":2}]`)

func TestArrayLength(t *testing.T) {
	l := ArrayLength(arr)

	if !reflect.DeepEqual(3, l) {
		t.Errorf("Unexpected result. Expected %v. Got %v", 3, l)
		t.FailNow()
	}
}

func BenchmarkArrayLength(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ArrayLength(arr)
	}
}

func TestArrayElement(t *testing.T) {
	element := ArrayElement(arr, 1)
	expected := `{"jsonrpc":"2.0","method":"sum","params":[1, 2],"id":2}`

	if expected != string(element) {
		t.Errorf("Unexpected result. Expected %v. Got %v", expected, string(element))
		t.FailNow()
	}
}

func BenchmarkArrayElement(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ArrayElement(arr, 2)
	}
}

func TestIsArray(t *testing.T) {
	isArr := IsArray([]byte(` [1, 2]`))
	if !reflect.DeepEqual(true, isArr) {
		t.Errorf("Unexpected result. Expected %v. Got %v", true, isArr)
		t.FailNow()
	}

	isArr2 := IsArray([]byte(`{"jsonrpc":"2.0","method":"sum","params":[1, 2],"id":2}`))
	if !reflect.DeepEqual(false, isArr2) {
		t.Errorf("Unexpected result. Expected %v. Got %v", false, isArr2)
		t.FailNow()
	}
}
