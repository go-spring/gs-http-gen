package jsonutil

import (
	"encoding/json/jsontext"
	"strings"
	"testing"
)

func BenchmarkDecodeJSON(b *testing.B) {
	strNumber := "3"
	strNumbers := "[3,9]"
	r := strings.NewReader("")
	d := jsontext.NewDecoder(r)

	b.Run("DecodeInt", func(b *testing.B) {
		for b.Loop() {
			r.Reset(strNumber)
			if _, err := DecodeInt[int8](d); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DecodeIntV2", func(b *testing.B) {
		for b.Loop() {
			r.Reset(strNumber)
			if _, err := DecodeIntV2[int8](d); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DecodeFloat", func(b *testing.B) {
		for b.Loop() {
			r.Reset(strNumber)
			if _, err := DecodeFloat[float32](d); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DecodeFloatV2", func(b *testing.B) {
		for b.Loop() {
			r.Reset(strNumber)
			if _, err := DecodeFloatV2[float32](d); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DecodeInts", func(b *testing.B) {
		for b.Loop() {
			r.Reset(strNumbers)
			if _, err := DecodeInts[int16](d); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DecodeIntsV2", func(b *testing.B) {
		for b.Loop() {
			r.Reset(strNumbers)
			if _, err := DecodeIntsV2[int16](d); err != nil {
				b.Fatal(err)
			}
		}
	})
}
