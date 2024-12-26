package response

type BaseError struct {
	Code int
	Err  error
	Data interface{}
}

func (err BaseError) Error() string {
	return err.Err.Error()
}
