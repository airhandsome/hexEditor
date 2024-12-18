<template>
  <div class="dialog-overlay" @click="handleOverlayClick">
    <div class="goto-dialog" @click.stop>
      <div class="dialog-header">
        <h3>跳转到位置</h3>
        <button class="close-button" @click="$emit('close')">&times;</button>
      </div>

      <div class="dialog-content">
        <!-- 跳转类型选择 -->
        <div class="goto-type">
          <label class="radio-label">
            <input type="radio" v-model="gotoType" value="offset">
            <span>偏移量</span>
          </label>
          <label class="radio-label">
            <input type="radio" v-model="gotoType" value="page">
            <span>页码</span>
          </label>
          <label class="radio-label">
            <input type="radio" v-model="gotoType" value="percent">
            <span>百分比</span>
          </label>
        </div>

        <!-- 数值输入 -->
        <div class="goto-input">
          <div class="input-group">
            <label>{{ inputLabel }}:</label>
            <input 
              type="text"
              v-model="inputValue"
              :placeholder="inputPlaceholder"
              @keyup.enter="handleGoto"
              @input="validateInput"
            >
          </div>

          <!-- 数值格式选择 -->
          <div class="format-select" v-if="gotoType === 'offset'">
            <label>格式:</label>
            <select v-model="numberFormat">
              <option value="hex">十六进制</option>
              <option value="dec">十进制</option>
            </select>
          </div>

          <!-- 相对位置选择 -->
          <div class="position-select">
            <label>位置:</label>
            <select v-model="relativePosition">
              <option value="start">从文件开始</option>
              <option value="current">从当前位置</option>
              <option value="end">从文件末尾</option>
            </select>
          </div>
        </div>

        <!-- 错误提示 -->
        <div class="error-message" v-if="errorMessage">
          {{ errorMessage }}
        </div>

        <!-- 位置预览 -->
        <div class="position-preview" v-if="previewPosition !== null">
          <span>目标位置: {{ formatPosition(previewPosition) }}</span>
        </div>
      </div>

      <div class="dialog-footer">
        <button class="primary-button" 
                @click="handleGoto"
                :disabled="!isValidInput">
          跳转
        </button>
        <button class="secondary-button" @click="$emit('close')">
          取消
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

// Props
const props = defineProps({
  currentPosition: {
    type: Number,
    default: 0
  },
  fileSize: {
    type: Number,
    default: 0
  },
  pageSize: {
    type: Number,
    default: 4096
  }
})

// 状态
const gotoType = ref('offset')
const inputValue = ref('')
const numberFormat = ref('hex')
const relativePosition = ref('start')
const errorMessage = ref('')
const previewPosition = ref(null)

// 计算属性
const inputLabel = computed(() => {
  switch (gotoType.value) {
    case 'offset': return '偏移量'
    case 'page': return '页码'
    case 'percent': return '百分比'
  }
})

const inputPlaceholder = computed(() => {
  switch (gotoType.value) {
    case 'offset': 
      return numberFormat.value === 'hex' ? '输入十六进制值 (例如: 1A0F)' : '输入十进制值'
    case 'page': 
      return '输入页码 (1-' + Math.ceil(props.fileSize / props.pageSize) + ')'
    case 'percent': 
      return '输入百分比 (0-100)'
  }
})

const isValidInput = computed(() => {
  return inputValue.value && !errorMessage.value
})

// 监听输入变化
watch([gotoType, inputValue, numberFormat, relativePosition], updatePreview)

// 输入验证
function validateInput() {
  errorMessage.value = ''
  
  switch (gotoType.value) {
    case 'offset':
      validateOffset()
      break
    case 'page':
      validatePage()
      break
    case 'percent':
      validatePercent()
      break
  }
  
  updatePreview()
}

function validateOffset() {
  if (numberFormat.value === 'hex') {
    if (!/^[0-9A-Fa-f]*$/.test(inputValue.value)) {
      errorMessage.value = '请输入有效的十六进制值'
      return false
    }
  } else {
    if (!/^\d*$/.test(inputValue.value)) {
      errorMessage.value = '请输入有效的十进制值'
      return false
    }
  }
  return true
}

function validatePage() {
  const page = parseInt(inputValue.value)
  const maxPage = Math.ceil(props.fileSize / props.pageSize)
  
  if (isNaN(page) || page < 1 || page > maxPage) {
    errorMessage.value = `请输入有效的页码 (1-${maxPage})`
    return false
  }
  return true
}

function validatePercent() {
  const percent = parseFloat(inputValue.value)
  if (isNaN(percent) || percent < 0 || percent > 100) {
    errorMessage.value = '请输入有效的百分比值 (0-100)'
    return false
  }
  return true
}

// 更新预览位置
function updatePreview() {
  if (!inputValue.value) {
    previewPosition.value = null
    return
  }

  const targetPosition = calculateTargetPosition()
  if (targetPosition !== null) {
    previewPosition.value = targetPosition
  }
}

// 计算目标位置
function calculateTargetPosition() {
  let position = 0
  
  switch (gotoType.value) {
    case 'offset':
      position = numberFormat.value === 'hex' 
        ? parseInt(inputValue.value, 16)
        : parseInt(inputValue.value, 10)
      break
      
    case 'page':
      position = (parseInt(inputValue.value) - 1) * props.pageSize
      break
      
    case 'percent':
      position = Math.floor(props.fileSize * parseFloat(inputValue.value) / 100)
      break
  }

  if (isNaN(position)) return null

  // 处理相对位置
  switch (relativePosition.value) {
    case 'current':
      position += props.currentPosition
      break
    case 'end':
      position = props.fileSize - position
      break
  }

  // 确保位置在有效范围内
  return Math.max(0, Math.min(position, props.fileSize - 1))
}

// 格式化位置显示
function formatPosition(position) {
  return `0x${position.toString(16).toUpperCase()} (${position})`
}

// 处理跳转
function handleGoto() {
  if (!isValidInput.value) return
  
  const targetPosition = calculateTargetPosition()
  if (targetPosition !== null) {
    emit('goto', targetPosition)
    emit('close')
  }
}

// 处理遮罩层点击
function handleOverlayClick(event) {
  if (event.target === event.currentTarget) {
    emit('close')
  }
}

// 定义事件
const emit = defineEmits(['close', 'goto'])
</script>

<style scoped>
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.goto-dialog {
  background-color: #252526;
  border: 1px solid #3D3D3D;
  border-radius: 4px;
  width: 360px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background-color: #2D2D2D;
  border-bottom: 1px solid #3D3D3D;
}

.dialog-header h3 {
  margin: 0;
  color: #D4D4D4;
  font-size: 14px;
}

.close-button {
  background: none;
  border: none;
  color: #D4D4D4;
  font-size: 18px;
  cursor: pointer;
  padding: 0 4px;
}

.close-button:hover {
  color: #FFFFFF;
}

.dialog-content {
  padding: 16px;
}

.goto-type {
  margin-bottom: 16px;
}

.radio-label {
  display: inline-flex;
  align-items: center;
  margin-right: 16px;
  color: #D4D4D4;
  cursor: pointer;
}

.goto-input {
  margin-bottom: 16px;
}

.input-group {
  margin-bottom: 12px;
}

.input-group label,
.format-select label,
.position-select label {
  display: block;
  margin-bottom: 6px;
  color: #D4D4D4;
}

input[type="text"] {
  width: 100%;
  padding: 6px 8px;
  background-color: #3D3D3D;
  border: 1px solid #4D4D4D;
  border-radius: 3px;
  color: #D4D4D4;
  font-family: 'Consolas', monospace;
}

input[type="text"]:focus {
  outline: none;
  border-color: #007ACC;
}

select {
  padding: 6px 8px;
  background-color: #3D3D3D;
  border: 1px solid #4D4D4D;
  border-radius: 3px;
  color: #D4D4D4;
  width: 150px;
}

.error-message {
  color: #F48771;
  margin-top: 8px;
  font-size: 12px;
}

.position-preview {
  margin-top: 16px;
  padding: 8px;
  background-color: #2D2D2D;
  border-radius: 3px;
  color: #4EC9B0;
  font-family: 'Consolas', monospace;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px;
  background-color: #2D2D2D;
  border-top: 1px solid #3D3D3D;
}

.primary-button,
.secondary-button {
  padding: 6px 12px;
  border-radius: 3px;
  border: 1px solid transparent;
  cursor: pointer;
  font-size: 12px;
}

.primary-button {
  background-color: #007ACC;
  color: #FFFFFF;
}

.primary-button:hover:not(:disabled) {
  background-color: #1B8BD0;
}

.primary-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.secondary-button {
  background-color: #3D3D3D;
  color: #D4D4D4;
}

.secondary-button:hover {
  background-color: #4D4D4D;
}
</style> 