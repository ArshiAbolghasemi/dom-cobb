package mongodb

type Pager struct {
	Page       uint
	Size       uint
	Total      uint
	TotalPages uint
}

func (p *Pager) GetLimit() uint {
	return p.Size
}

func (p *Pager) GetOffset() uint {
	return (p.Page - 1) * p.Size
}

func (p *Pager) SetTotal(total uint) {
	p.Total = total
	limit := p.GetLimit()
	if p.GetLimit() > 0 {
		p.TotalPages = (total + limit - 1) / limit
	}
}
