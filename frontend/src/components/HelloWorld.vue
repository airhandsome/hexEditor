<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'

const data = reactive({
  currentPosition: 0,
  fileSize: 0,
  isModified: false,
  hexData: [],
  selection: {
    start: -1,
    end: -1
  },
  visibleRows: 100, // 可见行数
  startRow: 0,      // 当前显示的起始行
  bytesPerRow: 16   // 每行显示的字节数
})

// 计算可见行
const visibleRows = computed(() => {
  const totalRows = Math.ceil(data.hexData.length / data.bytesPerRow)
  const rows = []
  for (let i = data.startRow; i < data.startRow + data.visibleRows && i < totalRows; i++) {
    rows.push(i)
  }
  return rows
})

// 获取行地址
function getRowAddress(row) {
  return row * data.bytesPerRow
}

// 计算属性
const canUndo = computed(() => false) // TODO: 实现撤销栈检查
const canRedo = computed(() => false) // TODO: 实现重做栈检查

// 格式化函数
function formatAddress(addr) {
  return addr.toString(16).padStart(8, '0').toUpperCase()
}

function formatByte(byte) {
  return byte?.toString(16).padStart(2, '0').toUpperCase() ?? '  '
}

function formatAscii(byte) {
  if (byte === undefined) return ' '
  return byte >= 32 && byte <= 126 ? String.fromCharCode(byte) : '.'
}

function formatFileSize(size) {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`
  return `${(size / 1024 / 1024).toFixed(2)} MB`
}

function formatSelection() {
  if (data.selection.start === -1) return ''
  const length = data.selection.end - data.selection.start + 1
  return `${length} 字节`
}

// 事件处理函数
function handleScroll(event) {
  const element = event.target
  const scrollTop = element.scrollTop
  const rowHeight = 20 // 行高
  data.startRow = Math.floor(scrollTop / rowHeight)
}

function handleByteClick(row, col) {
  const offset = row * data.bytesPerRow + col
  data.selection.start = offset
  data.selection.end = offset
  data.currentPosition = offset
}

function handleByteDoubleClick(row, col) {
  // TODO: 实现编辑模式
  console.log('Double click at:', row, col)
}

function handleAsciiClick(row, col) {
  handleByteClick(row, col)
}

// 文件操作函数
async function handleFileOpen() {
  try {
    console.log('开始选择文件...')
    const [fileHandle] = await window.showOpenFilePicker()
    const file = await fileHandle.getFile()
    
    const buffer = await file.arrayBuffer()
    data.hexData = Array.from(new Uint8Array(buffer))
    data.fileSize = file.size
    data.currentPosition = 0
    data.selection.start = -1
    data.selection.end = -1
    data.isModified = false
    data.startRow = 0
    
    console.log('文件加载完成，大小:', data.fileSize)
  } catch (err) {
    console.error('文件打开错误:', err)
  }
}

onMounted(() => {
  // TODO: 初始化
})

onUnmounted(() => {
  // TODO: 清理
})
</script>
