package dto

const defaultLimit = 50
const maxLimit = 200

type PaginationReq struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (p PaginationReq) LimitOffset() (limit, offset int) {
	limit = p.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	page := p.Page
	if page <= 0 {
		page = 1
	}
	offset = (page - 1) * limit
	return limit, offset
}
