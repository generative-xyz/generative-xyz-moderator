package structure

type CategoryData struct {
	Name string `json:"name"`
}

type UpdateCategoryData struct {
	ID *string `json:"id"`
	Name *string `json:"name"`
}

type FilterCategories struct {
	BaseFilters
}


