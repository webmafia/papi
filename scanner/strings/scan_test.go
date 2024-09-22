package strings

import "testing"

func BenchmarkScan(b *testing.B) {
	f := NewFactory()

	b.ResetTimer()

	for range b.N {
		var i int

		if err := ScanString(f, &i, "123"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkScanPrepared(b *testing.B) {
	f := NewFactory()
	var i int
	scan, err := GetScanner(f, &i)

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
