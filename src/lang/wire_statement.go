package lang

type Wire struct {
	Source      string
	SourceOut   string
	SourceRange Range
	Target      string
	TargetIn    string
	TargetRange Range
}

type WireStatement struct {
	Wires []Wire
}

func NewWireStatement() *WireStatement {
	return &WireStatement{
		Wires: make([]Wire, 0),
	}
}

func (stmt *WireStatement) add(wire Wire) {
	stmt.Wires = append(stmt.Wires, wire)
}
