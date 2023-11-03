package sqlparser

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

func (stmt *Statement) IsAggregate() bool {
	if stmt.JoinTable != "" {
		return true
	}

	for _, col := range stmt.SelectColumn {
		if col.GroupFunction != "" {
			return true
		}
	}

	return false
}
