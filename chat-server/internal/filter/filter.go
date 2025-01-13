package filter

type Filter struct {
	Page     int
	PageSize int
}

func (filter Filter) Limit() int {
	return filter.PageSize
}

func (filter Filter) Offset() int {
	return (filter.Page - 1) * filter.PageSize
}
