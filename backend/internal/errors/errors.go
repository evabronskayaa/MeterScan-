package errors

type SimpleError string

func (e SimpleError) Error() string {
	return string(e)
}
