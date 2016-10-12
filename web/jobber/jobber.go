package jobber

//Handler job handler
type Handler func(args []byte) error

//Jobber background job
type Jobber interface {
	Push(queue string, args interface{}) error
	Register(queue string, handler Handler)
	Start() error
}
