package jsonutil

import (
	"encoding/json/jsontext"
	"strings"
	"testing"
)

func BenchmarkDecodeJSON(b *testing.B) {
	strNumber := "3"
	r := strings.NewReader(strNumber)
	d := jsontext.NewDecoder(r)

	b.Run("DecodeInt", func(b *testing.B) {
		for b.Loop() {
			if _, err := DecodeInt[int8](d); err != nil {
				b.Fatal(err)
			}
			r.Reset(strNumber)
		}
	})

	b.Run("DecodeIntV2", func(b *testing.B) {
		for b.Loop() {
			if _, err := DecodeIntV2[int8](d); err != nil {
				b.Fatal(err)
			}
			r.Reset(strNumber)
		}
	})

	b.Run("DecodeFloat", func(b *testing.B) {
		for b.Loop() {
			if _, err := DecodeFloat[float32](d); err != nil {
				b.Fatal(err)
			}
			r.Reset(strNumber)
		}
	})

	b.Run("DecodeFloatV2", func(b *testing.B) {
		for b.Loop() {
			if _, err := DecodeFloatV2[float32](d); err != nil {
				b.Fatal(err)
			}
			r.Reset(strNumber)
		}
	})
}
