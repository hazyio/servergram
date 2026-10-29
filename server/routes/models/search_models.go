package models

type SearchRequest struct {
	Query string `query:"query"  maxLength:"200" doc:"The search query"`
}
type SearchResponse[T any] struct {
	Body struct {
		Items []T `json:"items" required:"true" doc:"The search results"`
		Total int `json:"total" required:"true" doc:"The total number of search results"`
	}
}
type SSHKeySearchData struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Encrypted bool   `json:"encrypted"`
}

type GroupSearchData struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
