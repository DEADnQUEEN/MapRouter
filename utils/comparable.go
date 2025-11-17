package utils

const (
	LESS    = -1
	EQUAL   = 0
	GREATER = 1
)

type Comparable interface {
	CompareValues(c Comparable) (int, error)
	ExactItem(c Comparable) (bool, error)
	String() string
}
