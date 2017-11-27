package lang

type Port struct {
	Id   string
	Pins int
}

type InputStatement struct {
	Ports []Port
}

func NewInputStatement() *InputStatement {
	return &InputStatement{
		Ports: make([]Port, 0),
	}
}

func (stmt *InputStatement) add(port Port) {
	stmt.Ports = append(stmt.Ports, port)
}

type OutputStatement struct {
	Ports []Port
}

func (stmt *OutputStatement) add(port Port) {
	stmt.Ports = append(stmt.Ports, port)
}

func NewOutputStatement() *OutputStatement {
	return &OutputStatement{
		Ports: make([]Port, 0),
	}
}
