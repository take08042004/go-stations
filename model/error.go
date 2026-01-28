package model

type ErrNotFound struct{}

func (e ErrNotFound) Error() string {
	return "TODO was not found"
}