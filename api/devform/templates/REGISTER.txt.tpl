在线表单「[[.TableName]]」生成代码 — 手动合并说明
================================================

1) 后端目录 gin-naive-admin-server
   - 将压缩包内 backend/api/[[.PackageName]]/ 拷贝到项目的 api/[[.PackageName]]/

2) 注册数据库表（initialize/gorm/init_gorm.go 的 AutoMigrate 列表中）追加一行：
   [[.PackageImport]].[[.EntityName]]{},

   并确保文件顶部 import：
   "gin-admin-server/api/[[.PackageName]]"

3) 注册路由（initialize/router/router_group.go）
   import 增加：
   "gin-admin-server/api/[[.PackageName]]"

   SetupPrivateRouters 中增加一行：
   [[.PackageName]].[[.EntityName]]Router.Init[[.EntityName]]Router(PrivateGroup)

4) 重新运行或编译后端。

5) 前端目录 gin-naive-admin-web
   - 将 frontend/src/views/[[.ViewPath]]/ 拷贝到 src/views/[[.ViewPath]]/

6) 在系统「菜单管理」中新增菜单：路径与 component 对应 views 路径，例如：
   component: /[[.ViewPath]]/index.vue
   （与项目内其他菜单填写方式一致）

7) 为角色勾选新菜单权限后刷新页面即可访问。
