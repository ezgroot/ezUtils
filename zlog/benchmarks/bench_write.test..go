package benchmarks

import (
	"fmt"
	"os"
	"testing"
)

// go test -bench=. -benchmem -run=none bench_write_test.go

func Benchmark_Write(b *testing.B) {
	b.Run("std", func(b *testing.B) {
		var err error
		os.Stdout, err = os.Open("/dev/null")
		if err != nil {
			fmt.Printf("open /dev/null error = %s\n", err)
			return
		}

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				os.Stdout.Write([]byte("111\n"))
			}

		})
	})

	b.Run("fmt", func(b *testing.B) {
		var err error
		os.Stdout, err = os.Open("/dev/null")
		if err != nil {
			fmt.Printf("open /dev/null error = %s\n", err)
			return
		}

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fmt.Printf("222\n")
			}
		})
	})

}
