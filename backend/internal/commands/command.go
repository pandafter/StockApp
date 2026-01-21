package commands

type Command interface {
	Execute() (interface{}, error)
}
