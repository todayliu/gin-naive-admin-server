package user

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"gin-admin-server/api/dict"
	"gin-admin-server/global"
	"gin-admin-server/model/response"
	"gin-admin-server/utils"
	"gin-admin-server/utils/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	// 单次导出行数上限（防止超大 CSV 占满内存）；与列表分页无关，导出的是「当前筛选下」的全量
	exportMaxRows = 50000
	importMaxRows  = 500
	importMaxBytes = 2 << 20 // 2MB
	defaultImportPW = "123456"
)

// userImportExportCSVHeader 用户导入/导出 CSV 表头（须与 ImportUsers 解析列一致）
func userImportExportCSVHeader() []string {
	return []string{"账号", "用户姓名", "昵称", "手机号", "邮箱", "性别", "状态", "部门名称", "角色名称", "职务名称", "密码"}
}

// 与前端字典类型编码一致（字典管理）
const (
	dictTypeCodeSex    = "sex"
	dictTypeCodeStatus = "status"
)

// dictValueToLabelMap type_code -> (字典 value 字符串 -> 中文 label)
func dictValueToLabelMap(typeCode string) map[string]string {
	var rows []dict.SysDictData
	if err := global.GNA_DB.Where("type_code = ? AND status = ?", typeCode, 1).Order("sort").Find(&rows).Error; err != nil {
		global.GNA_LOG.Warn("加载字典数据失败", zap.String("typeCode", typeCode), zap.Error(err))
		return nil
	}
	m := make(map[string]string, len(rows))
	for _, r := range rows {
		v := strings.TrimSpace(r.Value)
		if v != "" {
			m[v] = strings.TrimSpace(r.Label)
		}
	}
	return m
}

func exportDictLabel(m map[string]string, raw uint) string {
	if m == nil {
		return strconv.FormatUint(uint64(raw), 10)
	}
	key := strconv.FormatUint(uint64(raw), 10)
	if lab, ok := m[key]; ok && lab != "" {
		return lab
	}
	return key
}

// parseUintOrDictLabel 优先按数字解析；否则按字典 label 查 value（与导出中文列兼容）
func parseUintOrDictLabel(s string, typeCode string) (uint64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty")
	}
	if v, err := strconv.ParseUint(s, 10, 32); err == nil {
		return v, nil
	}
	var row dict.SysDictData
	err := global.GNA_DB.Where("type_code = ? AND status = ? AND label = ?", typeCode, 1, s).First(&row).Error
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(row.Value), 10, 32)
}

// ExportUsers 导出当前筛选条件下的用户为 CSV（UTF-8 BOM，与导入模板列一致）
func (us *_userService) ExportUsers(c *gin.Context) {
	var filters UserListFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		errMessage := validator.GetValidatorErrorMessage(err, filters)
		response.FailWithMessage(errMessage, c)
		return
	}
	db := buildUserListQuery(c, &filters)
	var list []SysUser
	if err := db.Order("create_time desc").Limit(exportMaxRows).Find(&list).Error; err != nil {
		global.GNA_LOG.Error("导出用户失败", zap.Error(err))
		response.FailWithMessage("导出用户失败", c)
		return
	}

	ids := make([]uint, 0, len(list))
	for i := range list {
		ids = append(ids, list[i].ID)
	}

	roleByUser := map[uint][]string{}
	if len(ids) > 0 {
		var rrows []struct {
			SysUserID uint   `gorm:"column:sys_user_id"`
			Name      string `gorm:"column:name"`
		}
		_ = global.GNA_DB.Table("sys_user_role ur").
			Select("ur.sys_user_id, r.name").
			Joins("JOIN sys_role r ON r.id = ur.sys_role_id").
			Where("ur.sys_user_id IN ?", ids).
			Order("r.name").
			Scan(&rrows).Error
		for _, row := range rrows {
			roleByUser[row.SysUserID] = append(roleByUser[row.SysUserID], row.Name)
		}
	}

	deptByUser := map[uint][]uint{}
	if len(ids) > 0 {
		var drows []struct {
			SysUserID       uint `gorm:"column:sys_user_id"`
			SysDepartmentID uint `gorm:"column:sys_department_id"`
		}
		_ = global.GNA_DB.Table("sys_user_department").
			Select("sys_user_id, sys_department_id").
			Where("sys_user_id IN ?", ids).
			Order("sys_department_id").
			Scan(&drows).Error
		for _, row := range drows {
			deptByUser[row.SysUserID] = append(deptByUser[row.SysUserID], row.SysDepartmentID)
		}
	}

	posByUser := map[uint][]uint{}
	if len(ids) > 0 {
		var prows []struct {
			SysUserID     uint `gorm:"column:sys_user_id"`
			SysJobLevelID uint `gorm:"column:sys_job_level_id"`
		}
		_ = global.GNA_DB.Table("sys_user_job_level").
			Select("sys_user_id, sys_job_level_id").
			Where("sys_user_id IN ?", ids).
			Order("sys_job_level_id").
			Scan(&prows).Error
		for _, row := range prows {
			posByUser[row.SysUserID] = append(posByUser[row.SysUserID], row.SysJobLevelID)
		}
	}

	sexLabelMap := dictValueToLabelMap(dictTypeCodeSex)
	statusLabelMap := dictValueToLabelMap(dictTypeCodeStatus)
	deptNameByID := loadDeptNameMap()
	jobLevelNameByID := loadJobLevelNameMap()

	buf := &bytes.Buffer{}
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(buf)
	if err := w.Write(userImportExportCSVHeader()); err != nil {
		global.GNA_LOG.Error("写入 CSV 表头失败", zap.Error(err))
		response.FailWithMessage("导出用户失败", c)
		return
	}
	for _, u := range list {
		depts := deptByUser[u.ID]
		if len(depts) == 0 && u.DepartmentId > 0 {
			depts = []uint{u.DepartmentId}
		}
		roles := roleByUser[u.ID]
		positions := posByUser[u.ID]
		if len(positions) == 0 && u.JobLevelID > 0 {
			positions = []uint{u.JobLevelID}
		}
		row := []string{
			u.Account,
			u.UName,
			u.UNickname,
			u.UMobile,
			u.UEmail,
			exportDictLabel(sexLabelMap, u.Gender),
			exportDictLabel(statusLabelMap, u.Status),
			joinNamesByUintIDs(depts, deptNameByID, ";"),
			strings.Join(roles, ";"),
			joinNamesByUintIDs(positions, jobLevelNameByID, ";"),
			"",
		}
		if err := w.Write(row); err != nil {
			global.GNA_LOG.Error("写入 CSV 行失败", zap.Error(err))
			response.FailWithMessage("导出用户失败", c)
			return
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		global.GNA_LOG.Error("刷新 CSV 失败", zap.Error(err))
		response.FailWithMessage("导出用户失败", c)
		return
	}

	c.Header("Content-Disposition", `attachment; filename="users_export.csv"`)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// DownloadImportTemplate 仅含表头的 CSV 模板（UTF-8 BOM）
func (us *_userService) DownloadImportTemplate(c *gin.Context) {
	var buf bytes.Buffer
	buf.WriteString("\xef\xbb\xbf")
	w := csv.NewWriter(&buf)
	if err := w.Write(userImportExportCSVHeader()); err != nil {
		global.GNA_LOG.Error("写入导入模板失败", zap.Error(err))
		response.FailWithMessage("生成模板失败", c)
		return
	}
	w.Flush()
	if err := w.Error(); err != nil {
		global.GNA_LOG.Error("刷新导入模板失败", zap.Error(err))
		response.FailWithMessage("生成模板失败", c)
		return
	}
	c.Header("Content-Disposition", `attachment; filename="user_import_template.csv"`)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

func loadDeptNameMap() map[uint]string {
	var rows []struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	_ = global.GNA_DB.Table("sys_department").Select("id, name").Find(&rows).Error
	m := make(map[uint]string, len(rows))
	for _, r := range rows {
		m[r.ID] = strings.TrimSpace(r.Name)
	}
	return m
}

func loadJobLevelNameMap() map[uint]string {
	var rows []struct {
		ID        uint   `gorm:"column:id"`
		LevelName string `gorm:"column:level_name"`
	}
	_ = global.GNA_DB.Table("sys_job_level").Select("id, level_name").Find(&rows).Error
	m := make(map[uint]string, len(rows))
	for _, r := range rows {
		m[r.ID] = strings.TrimSpace(r.LevelName)
	}
	return m
}

func joinNamesByUintIDs(ids []uint, idToName map[uint]string, sep string) string {
	if len(ids) == 0 {
		return ""
	}
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		if name, ok := idToName[id]; ok && name != "" {
			parts = append(parts, name)
		} else {
			parts = append(parts, strconv.FormatUint(uint64(id), 10))
		}
	}
	return strings.Join(parts, sep)
}

// parseDepartmentIDsFromCell 支持「1;2」或「研发部;市场部」等
func parseDepartmentIDsFromCell(s string) ([]uint, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	parts := splitSemicolonOrComma(s)
	ids := make([]uint, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if v, err := strconv.ParseUint(p, 10, 32); err == nil {
			ids = append(ids, uint(v))
			continue
		}
		var row struct {
			ID uint `gorm:"column:id"`
		}
		err := global.GNA_DB.Table("sys_department").Select("id").Where("name = ?", p).First(&row).Error
		if err != nil || row.ID == 0 {
			return nil, fmt.Errorf("部门 ID 或名称无法识别: %s", p)
		}
		ids = append(ids, row.ID)
	}
	return ids, nil
}

// parseRoleIDsFromCell 支持角色 code、或角色名称（与导出列一致）
func parseRoleIDsFromCell(s string) ([]uint, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	parts := splitSemicolonOrComma(s)
	ids := make([]uint, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		var rid uint
		_ = global.GNA_DB.Table("sys_role").Where("code = ?", p).Select("id").Scan(&rid)
		if rid == 0 {
			_ = global.GNA_DB.Table("sys_role").Where("name = ?", p).Select("id").Scan(&rid)
		}
		if rid == 0 {
			return nil, fmt.Errorf("角色编码或名称不存在: %s", p)
		}
		ids = append(ids, rid)
	}
	return ids, nil
}

// parsePositionIDsFromCell 支持职务数字 ID 或 level_name
func parsePositionIDsFromCell(s string) ([]uint, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	parts := splitSemicolonOrComma(s)
	ids := make([]uint, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if v, err := strconv.ParseUint(p, 10, 32); err == nil {
			ids = append(ids, uint(v))
			continue
		}
		var row struct {
			ID uint `gorm:"column:id"`
		}
		err := global.GNA_DB.Table("sys_job_level").Select("id").Where("level_name = ?", p).First(&row).Error
		if err != nil || row.ID == 0 {
			return nil, fmt.Errorf("职务 ID 或名称无法识别: %s", p)
		}
		ids = append(ids, row.ID)
	}
	return ids, nil
}

// canonicalImportKey 将表头（中文或英文）规范为内部列名，与 get() 使用的 key 一致
func canonicalImportKey(header string) string {
	h := strings.TrimSpace(header)
	if h == "" {
		return ""
	}
	switch h {
	case "账号", "登录账号":
		return "account"
	case "用户姓名", "姓名":
		return "uname"
	case "昵称":
		return "unickname"
	case "手机号", "手机":
		return "umobile"
	case "邮箱", "电子邮件":
		return "uemail"
	case "性别":
		return "gender"
	case "状态":
		return "status"
	case "部门名称", "部门":
		return "departmentids"
	case "部门ID", "部门id", "部门IDs", "部门ids":
		return "departmentids"
	case "角色名称":
		return "rolecodes"
	case "角色编码", "角色":
		return "rolecodes"
	case "职务名称", "职务":
		return "positionids"
	case "职务ID", "职务id":
		return "positionids"
	case "密码":
		return "password"
	}
	lower := strings.ToLower(h)
	switch lower {
	case "account":
		return "account"
	case "uname", "username":
		return "uname"
	case "unickname", "nickname":
		return "unickname"
	case "umobile", "mobile", "phone":
		return "umobile"
	case "uemail", "email":
		return "uemail"
	case "gender":
		return "gender"
	case "status":
		return "status"
	case "departmentids", "department_id", "department_ids":
		return "departmentids"
	case "rolecodes", "role_codes":
		return "rolecodes"
	case "positionids", "position_ids":
		return "positionids"
	case "password":
		return "password"
	}
	return lower
}

func buildImportColumnIndex(header []string) map[string]int {
	col := make(map[string]int)
	for i, name := range header {
		key := canonicalImportKey(name)
		if key == "" {
			continue
		}
		col[key] = i
	}
	return col
}

// ImportUsersResult 导入结果明细
type ImportUsersResult struct {
	SuccessCount int      `json:"successCount"`
	SkipCount    int      `json:"skipCount"`
	FailCount    int      `json:"failCount"`
	Errors       []string `json:"errors"`
}

// ImportUsers 从 CSV 批量导入用户（列与导出一致；password 为空则使用默认密码并需符合最少 6 位）
func (us *_userService) ImportUsers(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("请选择要上传的 CSV 文件", c)
		return
	}
	if file.Size > importMaxBytes {
		response.FailWithMessage("文件过大，请不超过 2MB", c)
		return
	}
	f, err := file.Open()
	if err != nil {
		response.FailWithMessage("读取文件失败", c)
		return
	}
	defer f.Close()

	body, err := io.ReadAll(io.LimitReader(f, importMaxBytes+1))
	if err != nil {
		response.FailWithMessage("读取文件失败", c)
		return
	}
	if len(body) > importMaxBytes {
		response.FailWithMessage("文件过大，请不超过 2MB", c)
		return
	}
	// 去 BOM
	if len(body) >= 3 && body[0] == 0xef && body[1] == 0xbb && body[2] == 0xbf {
		body = body[3:]
	}

	r := csv.NewReader(bytes.NewReader(body))
	r.LazyQuotes = true
	records, err := r.ReadAll()
	if err != nil {
		response.FailWithMessage("CSV 解析失败: "+err.Error(), c)
		return
	}
	if len(records) < 2 {
		response.FailWithMessage("CSV 至少需要表头与一行数据", c)
		return
	}
	if len(records)-1 > importMaxRows {
		response.FailWithMessage(fmt.Sprintf("单次最多导入 %d 行", importMaxRows), c)
		return
	}

	header := records[0]
	col := buildImportColumnIndex(header)
	required := []string{"account", "uname", "umobile", "gender", "status"}
	for _, k := range required {
		if _, ok := col[k]; !ok {
			response.FailWithMessage("缺少必填列（需中文或英文表头）: "+k, c)
			return
		}
	}

	result := ImportUsersResult{Errors: make([]string, 0, 8)}
	const maxErr = 30

	for lineIdx, rec := range records[1:] {
		lineNo := lineIdx + 2
		get := func(keys ...string) string {
			for _, k := range keys {
				if i, ok := col[k]; ok && i < len(rec) {
					return strings.TrimSpace(rec[i])
				}
			}
			return ""
		}
		account := get("account")
		uName := get("uname")
		uNickname := get("unickname")
		uMobile := get("umobile")
		uEmail := get("uemail")
		genderStr := get("gender")
		statusStr := get("status")
		deptStr := get("departmentids")
		roleStr := get("rolecodes")
		posStr := get("positionids")
		password := get("password")

		if account == "" {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 账号不能为空", lineNo))
			}
			continue
		}
		if uName == "" || uMobile == "" {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 用户名称与手机号不能为空", lineNo))
			}
			continue
		}
		gender, err1 := parseUintOrDictLabel(genderStr, dictTypeCodeSex)
		status, err2 := parseUintOrDictLabel(statusStr, dictTypeCodeStatus)
		if err1 != nil || err2 != nil {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 性别/状态须为字典中的数字或标签", lineNo))
			}
			continue
		}

		var exist SysUser
		if err := global.GNA_DB.Where("account = ?", account).First(&exist).Error; err == nil {
			result.SkipCount++
			continue
		} else if err != gorm.ErrRecordNotFound {
			global.GNA_LOG.Error("导入检查账号失败", zap.Error(err))
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 数据库错误", lineNo))
			}
			continue
		}

		if password == "" {
			password = defaultImportPW
		}
		if len(password) < 6 {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 密码至少 6 位（可留空使用默认）", lineNo))
			}
			continue
		}

		deptIDs, errDept := parseDepartmentIDsFromCell(deptStr)
		if errDept != nil {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: %v", lineNo, errDept))
			}
			continue
		}
		posIDs, errPos := parsePositionIDsFromCell(posStr)
		if errPos != nil {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: %v", lineNo, errPos))
			}
			continue
		}
		roleIDs, errRoles := parseRoleIDsFromCell(roleStr)
		if errRoles != nil {
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: %v", lineNo, errRoles))
			}
			continue
		}

		primaryDept := uint(0)
		if len(deptIDs) > 0 {
			primaryDept = deptIDs[0]
		}
		primaryJL := uint(0)
		if len(posIDs) > 0 {
			primaryJL = posIDs[0]
		}

		user := SysUser{
			UUID:         utils.GenerateUUID(),
			Account:      account,
			Password:     hashPasswordForStorage(password),
			UName:        uName,
			UNickname:    uNickname,
			UMobile:      uMobile,
			UEmail:       uEmail,
			Gender:       uint(gender),
			Status:       uint(status),
			DepartmentId: primaryDept,
			JobLevelID:   primaryJL,
		}

		err := global.GNA_DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
			for _, did := range deptIDs {
				if did == 0 {
					continue
				}
				if err := tx.Create(&SysUserDepartment{SysUserID: user.ID, SysDepartmentID: did}).Error; err != nil {
					return err
				}
			}
			for _, rid := range roleIDs {
				if err := tx.Exec("INSERT INTO sys_user_role (sys_user_id, sys_role_id) VALUES (?, ?)", user.ID, rid).Error; err != nil {
					return err
				}
			}
			for _, jl := range posIDs {
				if jl == 0 {
					continue
				}
				if err := tx.Create(&SysUserJobLevel{SysUserID: user.ID, SysJobLevelID: jl}).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			global.GNA_LOG.Error("导入用户失败", zap.Error(err))
			result.FailCount++
			if len(result.Errors) < maxErr {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: %v", lineNo, err))
			}
			continue
		}
		result.SuccessCount++
	}

	msg := fmt.Sprintf("成功 %d，跳过(账号已存在) %d，失败 %d", result.SuccessCount, result.SkipCount, result.FailCount)
	response.OkWithDetailed(result, msg, c)
}

func splitSemicolonOrComma(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var parts []string
	if strings.Contains(s, ";") {
		parts = strings.Split(s, ";")
	} else {
		parts = strings.Split(s, ",")
	}
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
