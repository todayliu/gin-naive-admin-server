package dict

import "github.com/gin-gonic/gin"

type _dictRouter struct{}

var DictRouter = new(_dictRouter)

// InitDictRouter 初始化字典模块路由
func (r *_dictRouter) InitDictRouter(Router *gin.RouterGroup) {
	dictRouter := Router.Group("dict")
	{
		dictRouter.GET("type/list", DictService.GetDictTypeList)
		dictRouter.POST("type/add", DictService.AddDictType)
		dictRouter.PUT("type/edit", DictService.EditDictType)
		dictRouter.DELETE("type/delete/:id", DictService.DeleteDictType)

		dictRouter.GET("data/list", DictService.GetDictDataList)
		dictRouter.POST("data/add", DictService.AddDictData)
		dictRouter.PUT("data/edit", DictService.EditDictData)
		dictRouter.DELETE("data/delete/:id", DictService.DeleteDictData)

		dictRouter.GET("data/:typeCode", DictService.GetDictByType)
	}
}
