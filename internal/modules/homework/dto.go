package homework

type CreateRequest struct {
	Name string `json:"name"`
}

type UpdateRequest struct {
	Name string `json:"name"`
}

type ListResponse struct {
	Items []Entity `json:"items"`
	Total int64    `json:"total"`
}
