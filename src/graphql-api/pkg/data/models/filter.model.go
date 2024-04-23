package models

type FilterModel struct {
	Table string `json:"table"`
	Query string `json:"query"`
	SearchText string `json:"searchText"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	OrderBy string `json:"orderBy"`
	FilterCondition string `json:"filterCondition"`
}
