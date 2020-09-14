package commands

type Command interface {
	Init([]string) error
	Name() string
	Run() error
}
