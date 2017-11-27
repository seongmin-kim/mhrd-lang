package lang

type Part struct {
	id     string
	module string
}

type PartStatement struct {
	Parts []Part
}

func NewPartStatement() *PartStatement {
	return &PartStatement{
		Parts: make([]Part, 0),
	}
}

func (stmt *PartStatement) add(part Part) {
	stmt.Parts = append(stmt.Parts, part)
}
