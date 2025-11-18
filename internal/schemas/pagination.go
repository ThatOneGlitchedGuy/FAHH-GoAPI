package schemas

type PageMeta struct {
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

type Page struct {
	Meta  PageMeta    `json:"meta"`
	Items interface{} `json:"items"`
}