package jparser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var arr = []byte(`[{"jsonrpc":"2.0","method":"sum","params":[1],"id":1}, {"jsonrpc":"2.0","method":"sum","params":[1, 2],"id":2}, {"jsonrpc":"2.0","method":"sum","params":[1, 2, 3, 4],"id":2}]`)

func TestArrayLength(t *testing.T) {
	l := ArrayLength(arr)

	require.Equal(t, 3, l)
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
	expected := []byte(`{"jsonrpc":"2.0","method":"sum","params":[1, 2],"id":2}`)

	require.Equal(t, expected, element)
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
	require.Equal(t, true, isArr)

	isArr2 := IsArray([]byte(`{"jsonrpc":"2.0","method":"sum","params":[1, 2],"id":2}`))
	require.Equal(t, false, isArr2)
}
