package floats

import (
	"math"

	"gonum.org/v1/gonum/floats"
)

type Slice []float64

func NewSlice(a ...float64) Slice {
	return Slice(a)
}

func (s *Slice) Push(v float64) {
	*s = append(*s, v)
}

func (s *Slice) Update(v float64) {
	*s = append(*s, v)
}

func (s *Slice) Pop(i int64) (v float64) {
	v = (*s)[i]
	*s = append((*s)[:i], (*s)[i+1:]...)
	return v
}

func (s Slice) Max() float64 {
	return floats.Max(s)
}

func (s Slice) Min() float64 {
	return floats.Min(s)
}

func (s Slice) Sum() (sum float64) {
	return floats.Sum(s)
}

func (s Slice) Mean() (mean float64) {
	length := len(s)
	if length == 0 {
		panic("zero length slice")
	}
	return s.Sum() / float64(length)
}

func (s Slice) Tail(size int) Slice {
	length := len(s)
	if length <= size {
		win := make(Slice, length)
		copy(win, s)
		return win
	}

	win := make(Slice, size)
	copy(win, s[length-size:])
	return win
}

func (collection Slice) Reverse() Slice {
	length := len(collection)
	half := length / 2

	for i := 0; i < half; i = i + 1 {
		j := length - 1 - i
		collection[i], collection[j] = collection[j], collection[i]
	}

	return collection
}

func (s Slice) Diff() (values Slice) {
	for i, v := range s {
		if i == 0 {
			values.Push(0)
			continue
		}
		values.Push(v - s[i-1])
	}
	return values
}

func (s Slice) PositiveValuesOrZero() (values Slice) {
	for _, v := range s {
		values.Push(math.Max(v, 0))
	}
	return values
}

func (s Slice) NegativeValuesOrZero() (values Slice) {
	for _, v := range s {
		values.Push(math.Min(v, 0))
	}
	return values
}

func (s Slice) Abs() (values Slice) {
	for _, v := range s {
		values.Push(math.Abs(v))
	}
	return values
}

func (s Slice) MulScalar(x float64) (values Slice) {
	for _, v := range s {
		values.Push(v * x)
	}
	return values
}

func (s Slice) DivScalar(x float64) (values Slice) {
	for _, v := range s {
		values.Push(v / x)
	}
	return values
}

func (s Slice) Mul(other Slice) (values Slice) {
	if len(s) != len(other) {
		panic("slice lengths do not match")
	}

	for i, v := range s {
		values.Push(v * other[i])
	}

	return values
}

func (s Slice) Dot(other Slice) float64 {
	return floats.Dot(s, other)
}

func (s Slice) Normalize() Slice {
	return s.DivScalar(s.Sum())
}

func (s *Slice) Last() float64 {
	length := len(*s)
	if length > 0 {
		return (*s)[length-1]
	}
	return 0.0
}

func (s *Slice) Index(i int) float64 {
	length := len(*s)
	if length-i <= 0 || i < 0 {
		return 0.0
	}
	return (*s)[length-i-1]
}

func (s *Slice) Length() int {
	return len(*s)
}

func (s *Slice) Len() int {
	return len(*s)
}

func (s Slice) Addr() *Slice {
	return &s
}
func (s Slice) CrossOver(t Slice) bool {
	if s.Length() == 0 {
		return false
	}
	return s.Index(0)-t.Index(0) > 0 && s.Index(1)-t.Index(1) < 0

}

func (s Slice) CrossUnder(t Slice) bool {
	if s.Length() == 0 {
		return false
	}
	return s.Index(0)-t.Index(0) < 0 && s.Index(1)-t.Index(1) > 0

}
func (s Slice) String(t Slice) float64 {

	return s.Index(0)

}
