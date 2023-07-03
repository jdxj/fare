package fare

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestMain(t *testing.M) {
	SetLoggerLevel(zapcore.DebugLevel)

	os.Exit(t.Run())
}

func TestCalcFare(t *testing.T) {
	f, err := os.Open("data.csv")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	defer f.Close()

	w, err := os.OpenFile("data_calc.csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	defer func() {
		_ = w.Sync()
		_ = w.Close()
	}()

	err = CalcFare(w, f)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}

func TestFormatFloat(t *testing.T) {
	res := strconv.FormatFloat(3.333, 'f', 2, 64)
	fmt.Println(res)
}

func TestSlice(t *testing.T) {
	dst := make([]int, 0, 10)
	fmt.Println(len(dst[:5]))
}

func TestFloat(t *testing.T) {
	var a, b float64
	fmt.Println(a == b)
}

func TestCalcFare2(t *testing.T) {
	tests := []struct {
		name   string
		weight float64
		want   int64
	}{
		{
			name:   "case1",
			weight: 0,
			want:   0,
		},
		{
			name:   "case2",
			weight: -1,
			want:   0,
		},
		{
			name:   "case3",
			weight: 0.1,
			want:   5,
		},
		{
			name:   "case4",
			weight: 1,
			want:   5,
		},
		{
			name:   "case5",
			weight: 1.01,
			want:   7,
		},
		{
			name:   "case6",
			weight: 2.1,
			want:   9,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := calcFare(test.weight)
			if got != test.want {
				t.Fatalf("name: %s, want: %d, got: %d",
					test.name, test.want, got)
			}
		})
	}
}
