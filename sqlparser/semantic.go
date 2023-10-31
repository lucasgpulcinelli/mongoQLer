package sqlparser

import "fmt"

type Column struct {
	Name          string
	GroupFunction string
}

type Comparision struct {
	Left  string
	Right string
	Op    string
}

type Statement struct {
	SelectColumn []Column

	FromTable string

	JoinTable    string
	JoinFromAttr string
	JoinToAttr   string

	Where     []Comparision
	BooleanOp string
}

func (stmt *Statement) ToMongoFind() string {
	return fmt.Sprintln(*stmt)
}
