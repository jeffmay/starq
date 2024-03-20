package tsq

import "fmt"

type AnyRealNumber interface {
	int | int32 | int64 | uint | uint32 | uint64 | float32 | float64
}

type RealNumber struct {
	string
}

func (n RealNumber) String() string {
	return n.string
}

func RealNum[N AnyRealNumber](n N) RealNumber {
	return RealNumber{fmt.Sprint(n)}
}
