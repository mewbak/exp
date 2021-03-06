package niterator

import (
	"testing"

	"github.com/egonelbre/exp/niterator/basic"
	"github.com/egonelbre/exp/niterator/instruct"
	"github.com/egonelbre/exp/niterator/onearr"
	"github.com/egonelbre/exp/niterator/onearrpremul"
	"github.com/egonelbre/exp/niterator/onearrrev"
	"github.com/egonelbre/exp/niterator/onearrrevadvance"
	"github.com/egonelbre/exp/niterator/onearrrevspecialize"
	"github.com/egonelbre/exp/niterator/onearrrevspecializeadvance"
	"github.com/egonelbre/exp/niterator/ordone"
	"github.com/egonelbre/exp/niterator/premul"
	"github.com/egonelbre/exp/niterator/shape"
	"github.com/egonelbre/exp/niterator/specialize"
	"github.com/egonelbre/exp/niterator/unroll"
	"github.com/egonelbre/exp/niterator/unrollinreverse"
	"github.com/egonelbre/exp/niterator/unrollinreverseadvance"
	"github.com/egonelbre/exp/niterator/unrollinreversebool"
	"github.com/egonelbre/exp/niterator/unrollinreversehardcode"
	"github.com/egonelbre/exp/niterator/unrollinreverseswitch"
	"github.com/egonelbre/exp/niterator/unrollinverse"
	"github.com/egonelbre/exp/niterator/unrollreverse"
)

const skipUseless = true

type TestSize struct {
	Name string
	AP   *shape.AP
}

var (
	indexOf100x = 5
	testSizes   = []TestSize{
		0: {"3x20x40x24", shape.New(3, 20, 40, 24)},
		1: {"24x20x40x3", shape.New(24, 20, 40, 3)},
		2: {"24x20x40x30", shape.New(24, 20, 40, 30)},
		3: {"5x50x100x150", shape.New(5, 50, 100, 150)},
		4: {"150x100x50x5", shape.New(150, 100, 50, 5)},
		5: {"100x100x100x100", shape.New(100, 100, 100, 100)},
	}

	verifyIndex = make([][]int, len(testSizes))
)

func init() {
	for api, testsize := range testSizes {
		ap := testsize.AP
		verify := make([]int, ap.TotalSize())

		it := basic.NewIterator(ap)
		i := 0
		for index, err := it.Next(); err == nil; index, err = it.Next() {
			i++
			verify[index] = i
		}
		for _, v := range verify {
			if v == 0 {
				panic("invalid iterator")
			}
		}

		verifyIndex[api] = verify
	}
}

func testIterator(t *testing.T, newiterator func(ap *shape.AP) Iterator) {
	t.Helper()
	for api, testsize := range testSizes {
		t.Run(testsize.Name, func(t *testing.T) {
			verify := verifyIndex[api]
			ap := testsize.AP
			count := 0
			i := 0
			it := newiterator(ap)
			for index, err := it.Next(); err == nil; index, err = it.Next() {
				count++
				i++
				if verify[index] != i {
					t.Fatalf("invalid at %d", index)
				}
			}

			if count != ap.TotalSize() {
				t.Fatalf("invalid count: got %d expected %d", count, ap.TotalSize())
			}
		})
	}
}

func TestBasic(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return basic.NewIterator(ap) })
}

func BenchmarkBasic(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := basic.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestInstruct(t *testing.T) {
	if skipUseless {
		t.Skip("skipping useless approach")
	}
	testIterator(t, func(ap *shape.AP) Iterator { return instruct.NewIterator(ap) })
}

func BenchmarkInstruct(b *testing.B) {
	if skipUseless {
		b.Skip("skipping useless approach")
	}
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := instruct.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestOrdone(t *testing.T) {
	if skipUseless {
		t.Skip("skipping useless approach")
	}
	testIterator(t, func(ap *shape.AP) Iterator { return ordone.NewIterator(ap) })
}

func BenchmarkOrdone(b *testing.B) {
	if skipUseless {
		b.Skip("skipping useless approach")
	}
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := ordone.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestPremul(t *testing.T) {
	if skipUseless {
		t.Skip("skipping useless approach")
	}
	testIterator(t, func(ap *shape.AP) Iterator { return premul.NewIterator(ap) })
}

func BenchmarkPremul(b *testing.B) {
	if skipUseless {
		b.Skip("skipping useless approach")
	}
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := premul.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestOnearr(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return onearr.NewIterator(ap) })
}

func BenchmarkOnearr(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := onearr.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestOnearrRev(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return onearrrev.NewIterator(ap) })
}

func BenchmarkOnearrRev(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := onearrrev.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestOnearrRevAdvance(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return onearrrevadvance.NewIterator(ap) })
}

func BenchmarkOnearrRevAdvance(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := onearrrevadvance.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestOnearrRevSpecialize(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return onearrrevspecialize.NewIterator(ap) })
}

func BenchmarkOnearrRevSpecialize(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := onearrrevspecialize.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestOnearrRevSpecializeAdvance(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return onearrrevspecializeadvance.NewIterator(ap) })
}

func BenchmarkOnearrRevSpecializeAdvance(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := onearrrevspecializeadvance.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestOnearrPremul(t *testing.T) {
	if skipUseless {
		t.Skip("skipping useless approach")
	}
	testIterator(t, func(ap *shape.AP) Iterator { return onearrpremul.NewIterator(ap) })
}

func BenchmarkOnearrPremul(b *testing.B) {
	if skipUseless {
		b.Skip("skipping useless approach")
	}
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := onearrpremul.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestSpecialize(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return specialize.NewIterator(ap) })
}

func BenchmarkSpecialize(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := specialize.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestUnroll(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return unroll.NewIterator(ap) })
}

func BenchmarkUnroll(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := unroll.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestUnrollReverse(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return unrollreverse.NewIterator(ap) })
}

func BenchmarkUnrollReverse(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := unrollreverse.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestUnrollInverse(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return unrollinverse.NewIterator(ap) })
}

func BenchmarkUnrollInverse(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := unrollinverse.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestUnrollInReverse(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return unrollinreverse.NewIterator(ap) })
}

func BenchmarkUnrollInReverse(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := unrollinreverse.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestUnrollInReverseAdvance(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return unrollinreverseadvance.NewIterator(ap) })
}

func BenchmarkUnrollInReverseAdvance(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := unrollinreverseadvance.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func TestUnrollInReverseSwitch(t *testing.T) {
	testIterator(t, func(ap *shape.AP) Iterator { return unrollinreverseswitch.NewIterator(ap) })
}

func BenchmarkUnrollInReverseSwitch(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := unrollinreverseswitch.NewIterator(ap)
				total := 0
				for index, err := it.Next(); err == nil; index, err = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

func BenchmarkUnrollInReverseBool(b *testing.B) {
	for api, testsize := range testSizes {
		b.Run(testsize.Name, func(b *testing.B) {
			verify := verifyIndex[api]
			ap := testsize.AP
			b.SetBytes(int64(ap.TotalSize()))
			for i := 0; i < b.N; i++ {
				it := unrollinreversebool.NewIterator(ap)
				total := 0
				for index, done := it.Next(); !done; index, done = it.Next() {
					total += verify[index] * index
				}
				_ = total
			}
		})
	}
}

// special
// func TestUnrollInReverseHardcode(t *testing.T) {
// 	testIterator(t, unrollinreversehardcode.NewIterator(shape.New(100, 100, 100, 100)))
// }

func BenchmarkUnrollInReverseHardcode(b *testing.B) {
	b.Run("100x100x100x100", func(b *testing.B) {
		verify := verifyIndex[indexOf100x]
		ap := shape.New(100, 100, 100, 100)
		b.SetBytes(int64(ap.TotalSize()))
		for i := 0; i < b.N; i++ {
			it := unrollinreversehardcode.NewIterator(ap)
			total := 0
			for index, err := it.Next(); err == nil; index, err = it.Next() {
				total += verify[index] * index
			}
			_ = total
		}
	})
}
