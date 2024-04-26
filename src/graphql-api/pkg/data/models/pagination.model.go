package models

type PaginationModel struct {
	Page        int
	PageSize    int
	TotalPages  int
	TotalItems  int
	HasNext     bool
	HasPrevious bool
}
