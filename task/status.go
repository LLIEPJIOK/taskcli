package task

type status int

const (
	ToDo status = iota
	InProgress
	Done
)

var (
	stringRepresentation = [3]string{"To do", "In progress", "Done"}
)

func (s status) String() string {
	return stringRepresentation[s]
}
