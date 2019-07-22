package data

type Creator interface {
	Create() (int, error)
}