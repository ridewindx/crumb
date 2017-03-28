package functional

type Iterator interface {
	HasNext() bool
	Next() (interface{}, error)
}
