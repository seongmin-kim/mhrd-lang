package lang

type Range struct {
	Start int
	End   int
}

func newRange(start, end int) Range {
	return Range{
		Start: start,
		End:   end,
	}
}

func (r Range) count() int {
	return r.Start - r.End + 1
}
