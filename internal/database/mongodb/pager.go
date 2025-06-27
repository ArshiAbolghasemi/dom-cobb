package mongodb

type Pager struct {
	Limit      uint
	Offset     uint
	Total      uint
	TotalPages uint
}

func (p *Pager) SetTotal(total uint) {
	p.Total = total
	if p.Limit > 0 {
		p.TotalPages = (total + p.Limit - 1) / p.Limit
	}
}
