package zap_log

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Cutter struct {
	level    string        // 日志级别(debug, info, warn, error, dpanic, panic, fatal)
	format   string        // 时间格式(2006-01-02)
	Director string        // 日志文件夹
	file     *os.File      // 文件句柄
	mutex    *sync.RWMutex // 读写锁

	// 轮转相关字段
	maxSize     int64  // 单个日志文件最大大小(字节)
	maxBackups  int    // 保留的旧日志文件最大数量
	maxAge      int    // 旧日志文件保留的最大天数
	compress    bool   // 是否压缩旧日志
	currentSize int64  // 当前文件大小
	currentDate string // 当前日期
}

type CutterOption func(*Cutter)

// WithCutterFormat 设置时间格式
func WithCutterFormat(format string) CutterOption {
	return func(c *Cutter) {
		c.format = format
	}
}

// WithMaxSize 设置单个日志文件最大大小(MB)
func WithMaxSize(sizeMB int) CutterOption {
	return func(c *Cutter) {
		c.maxSize = int64(sizeMB) * 1024 * 1024
	}
}

// WithMaxBackups 设置保留的旧日志文件最大数量
func WithMaxBackups(maxBackups int) CutterOption {
	return func(c *Cutter) {
		c.maxBackups = maxBackups
	}
}

// WithMaxAge 设置旧日志文件保留的最大天数
func WithMaxAge(maxAge int) CutterOption {
	return func(c *Cutter) {
		c.maxAge = maxAge
	}
}

// WithCompress 设置是否压缩旧日志
func WithCompress(compress bool) CutterOption {
	return func(c *Cutter) {
		c.compress = compress
	}
}

func NewCutter(director string, level string, options ...CutterOption) *Cutter {
	rotate := &Cutter{
		level:      level,
		Director:   director,
		mutex:      new(sync.RWMutex),
		maxSize:    100 * 1024 * 1024, // 默认100MB
		maxBackups: 0,
		maxAge:     0,
		compress:   false,
	}
	for i := 0; i < len(options); i++ {
		options[i](rotate)
	}

	return rotate
}

func (c *Cutter) Write(bytes []byte) (n int, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 提取business字段
	var business string
	if strings.Contains(string(bytes), "business") {
		business, bytes = c.extractBusiness(bytes)
	}

	// 生成文件路径
	format := time.Now().Format(c.format)
	filename := c.buildFilename(format, business)

	// 检查是否需要轮转
	if err = c.rotate(filename, format); err != nil {
		return 0, err
	}

	// 写入日志
	n, err = c.file.Write(bytes)
	if err != nil {
		return n, err
	}

	c.currentSize += int64(n)
	return n, nil
}

// rotate 检查并执行日志轮转
func (c *Cutter) rotate(filename, currentDate string) error {
	needRotate := false

	// 检查日期是否变化
	if c.currentDate != "" && c.currentDate != currentDate {
		needRotate = true
	}

	// 检查文件大小
	fileInfo, err := os.Stat(filename)
	if err == nil {
		if c.maxSize > 0 && fileInfo.Size() >= c.maxSize {
			needRotate = true
		}
	}

	// 执行轮转
	if needRotate && c.file != nil {
		_ = c.file.Close()
		c.file = nil

		// 重命名并可能压缩旧文件
		if err := c.rotateFile(filename); err != nil {
			return err
		}
	}

	// 打开或创建文件
	if c.file == nil {
		if err := c.openFile(filename); err != nil {
			return err
		}
		c.currentDate = currentDate

		if fileInfo, err := c.file.Stat(); err == nil {
			c.currentSize = fileInfo.Size()
		} else {
			c.currentSize = 0
		}
	}

	// 异步清理旧日志
	go c.cleanup(filepath.Dir(filename))

	return nil
}

// rotateFile 重命名并可能压缩当前日志文件
func (c *Cutter) rotateFile(filename string) error {
	timestamp := time.Now().Format("20060102-150405")
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	backupName := fmt.Sprintf("%s.%s%s", nameWithoutExt, timestamp, ext)

	// 重命名文件
	if err := os.Rename(filename, backupName); err != nil {
		return err
	}

	// 如果需要压缩，异步压缩
	if c.compress {
		go c.compressFile(backupName)
	}

	return nil
}

// compressFile 压缩日志文件
func (c *Cutter) compressFile(filename string) error {
	// 打开源文件
	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建压缩文件
	gzFilename := filename + ".gz"
	dst, err := os.Create(gzFilename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 创建gzip writer
	gzWriter := gzip.NewWriter(dst)
	defer gzWriter.Close()

	// 复制数据
	if _, err := io.Copy(gzWriter, src); err != nil {
		return err
	}

	// 删除原文件
	return os.Remove(filename)
}

// cleanup 清理过期的日志文件
func (c *Cutter) cleanup(dir string) {
	if c.maxBackups == 0 && c.maxAge == 0 {
		return
	}

	// 获取所有备份文件（包括.gz文件）
	pattern := filepath.Join(dir, "*.log.*")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	// 同时获取.gz文件
	gzPattern := filepath.Join(dir, "*.log.*.gz")
	gzFiles, err := filepath.Glob(gzPattern)
	if err == nil {
		files = append(files, gzFiles...)
	}

	// 获取文件信息并排序
	type fileWithInfo struct {
		path    string
		modTime time.Time
	}

	var backups []fileWithInfo
	for _, file := range files {
		info, err := os.Stat(file)
		if err == nil {
			backups = append(backups, fileWithInfo{
				path:    file,
				modTime: info.ModTime(),
			})
		}
	}

	// 按修改时间排序（最新的在前）
	for i := 0; i < len(backups)-1; i++ {
		for j := i + 1; j < len(backups); j++ {
			if backups[i].modTime.Before(backups[j].modTime) {
				backups[i], backups[j] = backups[j], backups[i]
			}
		}
	}

	// 按数量清理
	if c.maxBackups > 0 && len(backups) > c.maxBackups {
		for i := c.maxBackups; i < len(backups); i++ {
			_ = os.Remove(backups[i].path)
		}
		backups = backups[:c.maxBackups]
	}

	// 按时间清理
	if c.maxAge > 0 {
		cutoff := time.Now().AddDate(0, 0, -c.maxAge)
		for _, backup := range backups {
			if backup.modTime.Before(cutoff) {
				_ = os.Remove(backup.path)
			}
		}
	}
}

// extractBusiness 提取business字段
func (c *Cutter) extractBusiness(bytes []byte) (string, []byte) {
	var business string

	compile, err := regexp.Compile(`{"business":\s*"([^"]+)"}`)
	if err == nil && compile.Match(bytes) {
		finds := compile.FindSubmatch(bytes)
		if len(finds) > 1 {
			business = string(finds[1])
			bytes = compile.ReplaceAll(bytes, []byte(""))
			return business, bytes
		}
	}

	compile, err = regexp.Compile(`"business":\s*"([^"]+)"`)
	if err == nil && compile.Match(bytes) {
		finds := compile.FindSubmatch(bytes)
		if len(finds) > 1 {
			business = string(finds[1])
			bytes = compile.ReplaceAll(bytes, []byte(""))
		}
	}

	return business, bytes
}

// buildFilename 构建文件名
func (c *Cutter) buildFilename(format, business string) string {
	formats := make([]string, 0, 4)
	formats = append(formats, c.Director)
	if format != "" {
		formats = append(formats, format)
	}
	if business != "" {
		formats = append(formats, business)
	}
	formats = append(formats, c.level+".log")
	return filepath.Join(formats...)
}

// openFile 打开或创建日志文件
func (c *Cutter) openFile(filename string) error {
	dirname := filepath.Dir(filename)
	if err := os.MkdirAll(dirname, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	c.file = file
	return nil
}

// Close 关闭文件句柄
func (c *Cutter) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.file != nil {
		err := c.file.Close()
		c.file = nil
		return err
	}
	return nil
}
