package model

type DuplicateURLError struct {
	err error
}

func (e *DuplicateURLError) Error() string {
	return "Duplicate url error"
} 

func NewDuplicateURLError(err error) error {
	return &DuplicateURLError{err: err}
}