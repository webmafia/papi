package value

import "testing"

func BenchmarkScan(b *testing.B) {
	b.ResetTimer()

	for range b.N {
		var i int

		if err := ScanString(&i, "123"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkScanPrepared(b *testing.B) {
	var i int
	scan, err := GetScanner(&i)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for range b.N {
		if err := scan(&i, "123"); err != nil {
			b.Fatal(err)
		}
	}
}
