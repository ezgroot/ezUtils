package benchmarks

import (
	"fmt"
	"os"
	"testing"

	"github.com/ezgroot/ezUtils/zlog"
	"github.com/golang/glog"
	"go.uber.org/zap"
)

// go test -bench=. -benchmem -run=none bench_test.go

func Benchmark_Static(b *testing.B) {
	b.Run("zap", func(b *testing.B) {
		logger, _ := zap.NewProduction()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info("333\n")
			}
		})
	})

	b.Run("glog", func(b *testing.B) {
		var err error
		os.Stdout, err = os.Open("/dev/null")
		if err != nil {
			fmt.Printf("open /dev/null error = %s\n", err)
			return
		}
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				glog.Info("444\n")
			}
		})
	})

	b.Run("zlog", func(b *testing.B) {
		var err error
		os.Stdout, err = os.Open("/dev/null")
		if err != nil {
			fmt.Printf("open /dev/null error = %s\n", err)
			return
		}
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				zlog.Debug("555", "555")
			}
		})
	})
}
