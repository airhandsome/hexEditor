<template>
  <div class="hex-editor">
    <!-- 主工具栏 -->
    <ToolBar 
      @file-open="handleFileOpen"
      @file-save="handleSave"
      @find="handleFind"
      @goto="handleGoTo"
      :can-save="data.isModified"
    />
    
    <!-- 编辑区域 -->
    <div class="editor-container">
      <!-- 左侧文件结构面板 -->
      <FilePanel 
        v-if="showFilePanel"
        :file-name="data.fileName"
        :file-size="data.fileSize"
      />
      
      <!-- 主编辑区域 -->
      <div class="editor-main">
        <!-- 编辑器头部 -->
        <div class="editor-header">
          <div class="address-header">Address</div>
          <div class="hex-header">
            <span v-for="i in 16" :key="i">{{ (i-1).toString(16).toUpperCase().padStart(2, '0') }}</span>
          </div>
          <div class="ascii-header">ASCII</div>
        </div>
        
        <!-- 编辑器内容 -->
        <div class="editor-content" ref="editorContent" @scroll="handleScroll">
          <div v-for="row in visibleRows" :key="row" class="editor-row">
            <!-- 地址列 -->
            <div class="address-column">{{ formatAddress(getRowAddress(row)) }}</div>
            
            <!-- 十六进制列 -->
            <div class="hex-column">
              <span v-for="col in 16" 
                    :key="col" 
                    class="hex-byte"
                    :class="{
                      selected: isByteSelected(row, col-1),
                      active: isByteActive(row, col-1),
                      modified: isByteModified(row, col-1)
                    }"
                    @click="handleByteClick(row, col-1, $event)">
                {{ formatByte(getByteAt(row, col-1)) }}
              </span>
            </div>
            
            <!-- ASCII列 -->
            <div class="ascii-column">
              <span v-for="col in 16" 
                    :key="col"
                    class="ascii-char"
                    :class="{
                      selected: isByteSelected(row, col-1),
                      active: isByteActive(row, col-1),
                      modified: isByteModified(row, col-1)
                    }"
                    @click="handleByteClick(row, col-1, $event)">
                {{ formatAscii(getByteAt(row, col-1)) }}
              </span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 右侧数据检查器面板 -->
      <InspectorPanel 
        v-if="showInspector"
        :current-byte="getCurrentByte()"
        :offset="data.currentPosition"
      />
    </div>
    
    <!-- 状态栏 -->
    <StatusBar 
      :position="data.currentPosition"
      :selection="data.selection"
      :file-size="data.fileSize"
      :modified="data.isModified"
    />
    
    <!-- 对话框 -->
    <FindDialog 
      v-if="showFindDialog"
      @close="showFindDialog = false"
      @find="handleFindNext"
    />
    
    <GotoDialog
      v-if="showGotoDialog"
      @close="showGotoDialog = false"
      @goto="handleGotoOffset"
    />
    
    <!-- 更新使用 getCurrentByte 的地方 -->
    <div class="position-info" v-if="data.hexData.length">
      <span>位置: {{ formatAddress(data.currentPosition) }}</span>
      <span>字节: {{ formatByte(getCurrentByte()) }}</span>
    </div>
  </div>
</template>

<style scoped>
.hex-editor {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #1E1E1E;
  color: #D4D4D4;
  font-family: 'Consolas', monospace;
}

.editor-container {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.editor-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: #1E1E1E;
}

.editor-header {
  display: flex;
  padding: 4px 0;
  background-color: #2D2D2D;
  border-bottom: 1px solid #3D3D3D;
  font-family: 'Consolas', monospace;
  position: sticky;
  top: 0;
  z-index: 1;
}

.address-header {
  width: 100px;
  padding: 0 8px;
  color: #569CD6;
}

.hex-header {
  flex: 2;
  display: flex;
  padding: 0 8px;
}

.hex-header span {
  width: 24px;
  text-align: center;
  color: #569CD6;
}

.ascii-header {
  width: 160px;
  padding: 0 8px;
  color: #569CD6;
  border-left: 1px solid #3D3D3D;
}

.editor-content {
  flex: 1;
  overflow-y: auto;
  font-family: 'Consolas', monospace;
}

.editor-row {
  display: flex;
  height: 20px;
  line-height: 20px;
  white-space: nowrap;
}

.address-column {
  width: 100px;
  padding: 0 8px;
  color: #569CD6;
}

.hex-column {
  flex: 2;
  padding: 0 8px;
  display: flex;
}

.hex-byte {
  width: 24px;
  text-align: center;
  cursor: pointer;
  user-select: none;
}

.ascii-column {
  width: 160px;
  padding: 0 8px;
  border-left: 1px solid #3D3D3D;
  display: flex;
}

.ascii-char {
  width: 10px;
  text-align: center;
  cursor: pointer;
  user-select: none;
}

.selected {
  background-color: #264F78;
}

.active {
  background-color: #3D3D3D;
}

.modified {
  color: #CE9178;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 14px;
  height: 14px;
}

::-webkit-scrollbar-track {
  background: #1E1E1E;
}

::-webkit-scrollbar-thumb {
  background: #424242;
  border: 3px solid #1E1E1E;
  border-radius: 7px;
}

::-webkit-scrollbar-thumb:hover {
  background: #4D4D4D;
}

::-webkit-scrollbar-corner {
  background: #1E1E1E;
}
</style> 

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { OpenFileWithData, ReadPage, WriteByte, SaveFile } from '../../wailsjs/go/main/HexEditor'
import ToolBar from './ToolBar.vue'
import FilePanel from './FilePanel.vue'
import InspectorPanel from './InspectorPanel.vue'
import StatusBar from './StatusBar.vue'
import FindDialog from './FindDialog.vue'
import GotoDialog from './GotoDialog.vue'

// 状态管理
const data = reactive({
  fileName: '',
  fileSize: 0,
  hexData: [],
  isModified: false,
  currentPosition: 0,
  selection: {
    start: -1,
    end: -1
  },
  modifiedBytes: new Set(),
  startRow: 0,
  visibleRows: 100,
  bytesPerRow: 16,
  undoStack: [],
  redoStack: []
})

// UI状态
const showFilePanel = ref(true)
const showInspector = ref(true)
const showFindDialog = ref(false)
const showGotoDialog = ref(false)
const editorContent = ref(null)

// 计算属性
const visibleRows = computed(() => {
  const totalRows = Math.ceil(data.hexData.length / data.bytesPerRow)
  const rows = []
  for (let i = data.startRow; i < Math.min(data.startRow + data.visibleRows, totalRows); i++) {
    rows.push(i)
  }
  return rows
})

const currentByte = computed(() => {
  return data.hexData[data.currentPosition]
})

// 方法定义
function getCurrentByte() {
  return data.hexData[data.currentPosition]
}

function getByteAt(row, col) {
  const offset = row * data.bytesPerRow + col
  return data.hexData[offset]
}

function isByteSelected(row, col) {
  const offset = row * data.bytesPerRow + col
  if (data.selection.start === -1 || data.selection.end === -1) return false
  
  const start = Math.min(data.selection.start, data.selection.end)
  const end = Math.max(data.selection.start, data.selection.end)
  return offset >= start && offset <= end
}

function isByteActive(row, col) {
  const offset = row * data.bytesPerRow + col
  return offset === data.currentPosition
}

function isByteModified(row, col) {
  const offset = row * data.bytesPerRow + col
  return data.modifiedBytes.has(offset)
}

function getRowAddress(row) {
  return row * data.bytesPerRow
}

function formatAddress(addr) {
  return addr.toString(16).padStart(8, '0').toUpperCase()
}

function formatByte(byte) {
  return byte !== undefined ? byte.toString(16).padStart(2, '0').toUpperCase() : '  '
}

function formatAscii(byte) {
  if (byte === undefined) return ' '
  return byte >= 32 && byte <= 126 ? String.fromCharCode(byte) : '.'
}

// 文件处理函数
async function handleFileOpen() {
  try {
    const [fileHandle] = await window.showOpenFilePicker()
    const file = await fileHandle.getFile()
    const buffer = await file.arrayBuffer()
    const bytes = new Uint8Array(buffer)
    
    // 更新状态
    data.fileName = file.name
    data.fileSize = file.size
    data.hexData = Array.from(bytes)
    data.currentPosition = 0
    data.selection = { start: -1, end: -1 }
    data.isModified = false
    data.modifiedBytes.clear()
    data.startRow = 0
    
    // 调用后端打开文件
    await window.go.main.HexEditor.OpenFileWithData(file.name, Array.from(bytes))
  } catch (err) {
    console.error('打开文件失败:', err)
  }
}

async function handleSave() {
  if (!data.isModified) return
  
  try {
    await window.go.main.HexEditor.SaveFile()
    data.isModified = false
    data.modifiedBytes.clear()
  } catch (err) {
    console.error('保存文件失败:', err)
  }
}

// 数据加载函数
async function loadPage(pageNum) {
  try {
    const pageData = await window.go.main.HexEditor.ReadPage(pageNum)
    data.hexData = pageData
    data.startRow = pageNum * (data.visibleRows - 1)
  } catch (err) {
    console.error('加载数据失败:', err)
  }
}

// 编辑功能
async function handleByteEdit(row, col, newValue) {
  const offset = row * data.bytesPerRow + col
  try {
    const value = parseInt(newValue, 16)
    if (isNaN(value) || value < 0 || value > 255) return
    
    await window.go.main.HexEditor.WriteByte(offset, value)
    data.hexData[offset] = value
    data.modifiedBytes.add(offset)
    data.isModified = true
    
    // 添加到撤销栈
    data.undoStack.push({
      offset,
      oldValue: data.hexData[offset],
      newValue: value
    })
    data.redoStack = []
  } catch (err) {
    console.error('编辑失败:', err)
  }
}

// 选择功能
function handleByteClick(row, col, event) {
  const offset = row * data.bytesPerRow + col
  
  if (event.shiftKey && data.selection.start !== -1) {
    // 扩展选择
    data.selection.end = offset
  } else {
    // 新选择
    data.selection = {
      start: offset,
      end: offset
    }
  }
  data.currentPosition = offset
}

// 滚动处理
function handleScroll(event) {
  const { scrollTop, clientHeight, scrollHeight } = event.target
  const rowHeight = 20 // 每行高度
  const currentRow = Math.floor(scrollTop / rowHeight)
  data.startRow = currentRow
}

// 查找功能
async function handleSearchResult(position) {
  if (position >= 0) {
    const page = Math.floor(position / (data.bytesPerRow * data.visibleRows))
    await loadPage(page)
    data.currentPosition = position
    scrollToPosition(position)
  }
}

function scrollToPosition(position) {
  const row = Math.floor(position / data.bytesPerRow)
  const rowElement = editorContent.value.children[row - data.startRow]
  if (rowElement) {
    rowElement.scrollIntoView({ block: 'center' })
  }
}

// 跳转功能
async function handleGoto(position) {
  const page = Math.floor(position / (data.bytesPerRow * data.visibleRows))
  await loadPage(page)
  data.currentPosition = position
  scrollToPosition(position)
}

// 撤销/重做功能
async function handleUndo() {
  try {
    await window.go.main.HexEditor.Undo()
    await loadPage(Math.floor(data.currentPosition / (data.bytesPerRow * data.visibleRows)))
  } catch (err) {
    console.error('撤销失败:', err)
  }
}

async function handleRedo() {
  try {
    await window.go.main.HexEditor.Redo()
    await loadPage(Math.floor(data.currentPosition / (data.bytesPerRow * data.visibleRows)))
  } catch (err) {
    console.error('重做失败:', err)
  }
}

// 键盘事件处理
function handleKeyDown(event) {
  if (!data.hexData.length) return
  
  const currentRow = Math.floor(data.currentPosition / data.bytesPerRow)
  const currentCol = data.currentPosition % data.bytesPerRow
  
  switch (event.key) {
    case 'ArrowLeft':
      if (data.currentPosition > 0) {
        data.currentPosition--
        ensurePositionVisible()
      }
      break
    case 'ArrowRight':
      if (data.currentPosition < data.hexData.length - 1) {
        data.currentPosition++
        ensurePositionVisible()
      }
      break
    case 'ArrowUp':
      if (currentRow > 0) {
        data.currentPosition -= data.bytesPerRow
        ensurePositionVisible()
      }
      break
    case 'ArrowDown':
      if (data.currentPosition + data.bytesPerRow < data.hexData.length) {
        data.currentPosition += data.bytesPerRow
        ensurePositionVisible()
      }
      break
  }
}

// 确保当前位置可见
function ensurePositionVisible() {
  const row = Math.floor(data.currentPosition / data.bytesPerRow)
  if (row < data.startRow) {
    data.startRow = row
  } else if (row >= data.startRow + data.visibleRows) {
    data.startRow = row - data.visibleRows + 1
  }
}

// 生命周期钩子
onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})

// 暴露给模板使用的方法和属性
defineExpose({
  handleFileOpen,
  handleSave,
  handleByteEdit,
  handleByteClick,
  handleScroll,
  handleSearchResult,
  handleGoto,
  handleUndo,
  handleRedo,
  getByteAt,
  getCurrentByte,
  isByteSelected,
  isByteActive,
  isByteModified,
  getRowAddress,
  formatAddress,
  formatByte,
  formatAscii,
  currentByte,
  data
})
</script> 