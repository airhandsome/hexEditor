package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type HexEditor struct {
	FilePath    string
	FileSize    int64
	Buffer      []byte
	IsModified  bool
	PageSize    int
	CurrentPage int
	UndoStack   []Operation // 撤销栈
	RedoStack   []Operation // 重做栈
	Selection   Selection   // 选择区域
}

type Selection struct {
	Start int64
	End   int64
}

type Operation struct {
	Offset  int64
	OldData []byte
	NewData []byte
}

// 基础方法
type FileOperations interface {
	OpenFile(path string) error
	SaveFile() error
	ReadPage(pageNum int) ([]byte, error)
	WriteByte(offset int64, data byte) error
	WriteBytes(offset int64, data []byte) error
}

// 添加一个临时文件存储路径
var tempDir string

func init() {
	// 初始化临时目录
	tempDir = filepath.Join(os.TempDir(), "hexeditor")
	os.MkdirAll(tempDir, 0755)
}

// OpenFileWithData 使用文件数据打开文件
func (h *HexEditor) OpenFileWithData(filename string, data []byte) error {
	log.Printf("Opening file with data: %s, size: %d", filename, len(data))

	// 创建临时文件
	tempPath := filepath.Join(tempDir, filename)

	// 写入数据到临时文件
	err := os.WriteFile(tempPath, data, 0644)
	if err != nil {
		log.Printf("Error writing temp file: %v", err)
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// 打开临时文件
	file, err := os.Open(tempPath)
	if err != nil {
		log.Printf("Error opening temp file: %v", err)
		return fmt.Errorf("failed to open temp file: %w", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		return fmt.Errorf("failed to get file info: %w", err)
	}

	log.Printf("File size: %d bytes", fileInfo.Size())

	h.FilePath = tempPath
	h.FileSize = fileInfo.Size()
	h.IsModified = false
	h.CurrentPage = 0

	// 初始化或清空撤销/重做栈
	h.UndoStack = make([]Operation, 0)
	h.RedoStack = make([]Operation, 0)

	log.Printf("File opened successfully at: %s", tempPath)
	return nil
}

// 添加清理方法
func (h *HexEditor) Cleanup() error {
	if h.FilePath != "" && strings.HasPrefix(h.FilePath, tempDir) {
		return os.Remove(h.FilePath)
	}
	return nil
}

// 打开文件
func (h *HexEditor) OpenFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	h.FilePath = path
	h.FileSize = fileInfo.Size()
	h.IsModified = false

	return nil
}

// 读取指定页的数据
func (h *HexEditor) ReadPage(pageNum int) ([]byte, error) {
	file, err := os.Open(h.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	offset := int64(pageNum * h.PageSize)
	buffer := make([]byte, h.PageSize)

	// 设置读取位置
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	// 读取数据
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buffer[:n], nil
}

// 写入单个字节
func (h *HexEditor) WriteByte(offset int64, data byte) error {
	file, err := os.OpenFile(h.FilePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte{data})
	if err != nil {
		return err
	}

	h.IsModified = true
	return nil
}

// 保存文件
func (h *HexEditor) SaveFile() error {
	if !h.IsModified {
		return nil
	}
	// 实现文件保存逻辑
	h.IsModified = false
	return nil
}

// 搜索接口
func (h *HexEditor) Search(pattern []byte, startOffset int64) ([]int64, error) {
	file, err := os.Open(h.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []int64
	buffer := make([]byte, 4096) // 读取缓冲区
	offset := startOffset

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// 在当前缓冲区中查找
		for i := 0; i <= n-len(pattern); i++ {
			match := true
			for j := 0; j < len(pattern); j++ {
				if buffer[i+j] != pattern[j] {
					match = false
					break
				}
			}
			if match {
				results = append(results, offset+int64(i))
			}
		}
		offset += int64(n)
	}

	return results, nil
}

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (h *HexEditor) handleError(err error) ErrorResponse {
	if err != nil {
		log.Printf("Error: %v", err)
		return ErrorResponse{
			Error:   true,
			Message: err.Error(),
		}
	}
	return ErrorResponse{Error: false}
}

// 复制选中内容
func (h *HexEditor) CopySelection() ([]byte, error) {
	if h.Selection.Start == h.Selection.End {
		return nil, nil
	}

	file, err := os.Open(h.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	length := h.Selection.End - h.Selection.Start
	data := make([]byte, length)

	_, err = file.Seek(h.Selection.Start, 0)
	if err != nil {
		return nil, err
	}

	_, err = file.Read(data)
	return data, err
}

// 粘贴内容
func (h *HexEditor) PasteBytes(offset int64, data []byte) error {
	file, err := os.OpenFile(h.FilePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 保存旧数据用于撤销
	oldData := make([]byte, len(data))
	_, err = file.ReadAt(oldData, offset)
	if err != nil && err != io.EOF {
		return err
	}

	// 记录操作
	op := Operation{
		Offset:  offset,
		OldData: oldData,
		NewData: data,
	}
	h.UndoStack = append(h.UndoStack, op)
	h.RedoStack = nil

	// 写入新数据
	_, err = file.WriteAt(data, offset)
	if err != nil {
		return err
	}

	h.IsModified = true
	return nil
}

// 撤销操作
func (h *HexEditor) Undo() error {
	if len(h.UndoStack) == 0 {
		return nil
	}

	op := h.UndoStack[len(h.UndoStack)-1]
	h.UndoStack = h.UndoStack[:len(h.UndoStack)-1]

	// 执行撤销
	err := h.writeBytes(op.Offset, op.OldData)
	if err != nil {
		return err
	}

	h.RedoStack = append(h.RedoStack, op)
	return nil
}

// 重做操作
func (h *HexEditor) Redo() error {
	if len(h.RedoStack) == 0 {
		return nil
	}

	op := h.RedoStack[len(h.RedoStack)-1]
	h.RedoStack = h.RedoStack[:len(h.RedoStack)-1]

	// 执行重做
	err := h.writeBytes(op.Offset, op.NewData)
	if err != nil {
		return err
	}

	h.UndoStack = append(h.UndoStack, op)
	return nil
}

// 查找下一个
func (h *HexEditor) FindNext(pattern []byte, startOffset int64) (int64, error) {
	results, err := h.Search(pattern, startOffset)
	if err != nil {
		return -1, err
	}

	if len(results) > 0 {
		return results[0], nil
	}
	return -1, nil
}

// 替换选中内容
func (h *HexEditor) Replace(offset int64, oldData []byte, newData []byte) error {
	return h.PasteBytes(offset, newData)
}

// writeBytes 写入字节数组到指定偏移位置
func (h *HexEditor) writeBytes(offset int64, data []byte) error {
	// 参数检查
	if offset < 0 {
		return fmt.Errorf("invalid offset: %d", offset)
	}
	if len(data) == 0 {
		return nil
	}

	// 打开文件
	file, err := os.OpenFile(h.FilePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 检查写入范围是否超出文件大小
	if offset+int64(len(data)) > h.FileSize {
		return fmt.Errorf("write operation exceeds file size")
	}

	// 定位到指定位置
	_, err = file.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("failed to seek to position %d: %w", offset, err)
	}

	// 写入数据
	n, err := file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	// 检查是否完全写入
	if n != len(data) {
		return fmt.Errorf("incomplete write: wrote %d of %d bytes", n, len(data))
	}

	// 同步文件到磁盘
	err = file.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	// 更新修改状态
	h.IsModified = true

	// 如果写入的位置在当前缓存页中，更新缓存
	pageStart := int64(h.CurrentPage * h.PageSize)
	pageEnd := pageStart + int64(h.PageSize)
	if offset >= pageStart && offset < pageEnd {
		copyStart := offset - pageStart
		copyEnd := copyStart + int64(len(data))
		if copyEnd > int64(len(h.Buffer)) {
			copyEnd = int64(len(h.Buffer))
		}
		copy(h.Buffer[copyStart:copyEnd], data)
	}

	return nil
}

// 批量写入优化版本，用于大量数据写入
func (h *HexEditor) WriteBytesBuffered(offset int64, data []byte) error {
	const bufferSize = 4096 // 4KB 缓冲区

	if len(data) <= bufferSize {
		return h.writeBytes(offset, data)
	}

	// 对大数据进行分块写入
	for i := 0; i < len(data); i += bufferSize {
		end := i + bufferSize
		if end > len(data) {
			end = len(data)
		}

		chunk := data[i:end]
		err := h.writeBytes(offset+int64(i), chunk)
		if err != nil {
			return fmt.Errorf("failed to write chunk at offset %d: %w", offset+int64(i), err)
		}
	}

	return nil
}

// 安全写入方法，包含错误恢复机制
func (h *HexEditor) WriteBytesWithBackup(offset int64, data []byte) error {
	// 先备份原有数据
	file, err := os.OpenFile(h.FilePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	backup := make([]byte, len(data))
	_, err = file.ReadAt(backup, offset)
	if err != nil && err != io.EOF {
		return err
	}

	// 尝试写入新数据
	err = h.writeBytes(offset, data)
	if err != nil {
		// 发生错误时恢复原有数据
		restoreErr := h.writeBytes(offset, backup)
		if restoreErr != nil {
			return fmt.Errorf("write failed: %v, restore failed: %v", err, restoreErr)
		}
		return err
	}

	// 记录操作到撤销栈
	op := Operation{
		Offset:  offset,
		OldData: backup,
		NewData: data,
	}
	h.UndoStack = append(h.UndoStack, op)
	// 清空重做栈
	h.RedoStack = nil

	return nil
}
