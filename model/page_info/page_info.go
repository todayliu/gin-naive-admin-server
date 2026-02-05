package page_info

type PageInfo struct {
	PageNo   int `json:"pageNo" form:"pageNo" binding:"required" message:"pageNo不能为空"`
	PageSize int `json:"pageSize" from:"pageSize" binding:"required" message:"pageSize不能为空"`
}
