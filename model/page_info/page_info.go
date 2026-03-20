package page_info

// PageInfo 分页请求通用参数
type PageInfo struct {
	PageNo   int `json:"pageNo" form:"pageNo" binding:"required" message:"pageNo不能为空"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"required" message:"pageSize不能为空"`
}
