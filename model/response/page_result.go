package response

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	PageNo   int         `json:"pageNo"`
	PageSize int         `json:"pageSize"`
}
