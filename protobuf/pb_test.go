package main

import (
	"testing"
)

func BenchmarkMarshalJson(b *testing.B) {
	b.ResetTimer()
	var a []byte

	for i := 0; i < b.N; i++ {
		c, err := marshalJson()
		a = c
		if err != nil {
			b.Fatal(err)
		}
	}
	b.ReportMetric(float64(len(a)), "json")
}

func BenchmarkMarshalPb(b *testing.B) {
	b.ResetTimer()
	var a []byte
	for i := 0; i < b.N; i++ {
		c, err := marshalPb()
		if err != nil {
			b.Fatal(err)
		}
		a = c
	}
	b.ReportMetric(float64(len(a)), "pb")
}
