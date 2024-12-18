<template>
  <div class="dialog-overlay" @click="handleOverlayClick">
    <div class="find-dialog" @click.stop>
      <div class="dialog-header">
        <h3>查找</h3>
        <button class="close-button" @click="$emit('close')">&times;</button>
      </div>

      <div class="dialog-content">
        <!-- 搜索类型选择 -->
        <div class="search-type">
          <label class="radio-label">
            <input type="radio" v-model="searchType" value="text">
            <span>文本</span>
          </label>
          <label class="radio-label">
            <input type="radio" v-model="searchType" value="hex">
            <span>十六进制</span>
          </label>
        </div>

        <!-- 搜索输入 -->
        <div class="search-input">
          <label>查找内容:</label>
          <input 
            v-if="searchType === 'text'"
            type="text"
            v-model="searchText"
            placeholder="输入要查找的文本"
            @keyup.enter="handleFind"
          >
          <input 
            v-else
            type="text"
            v-model="searchHex"
            placeholder="输入十六进制值 (例如: FF 00 BE)"
            @keyup.enter="handleFind"
            @input="validateHex"
          >
        </div>

        <!-- 搜索选项 -->
        <div class="search-options">
          <label class="checkbox-label">
            <input type="checkbox" v-model="matchCase">
            <span>区分大小写</span>
          </label>
          <label class="checkbox-label">
            <input type="checkbox" v-model="searchBackward">
            <span>向上搜索</span>
          </label>
          <label class="checkbox-label">
            <input type="checkbox" v-model="wrapAround">
            <span>循环搜索</span>
          </label>
        </div>

        <!-- 范围选项 -->
        <div class="search-range">
          <label>搜索范围:</label>
          <select v-model="searchRange">
            <option value="all">整个文件</option>
            <option value="selection">选定范围</option>
            <option value="cursor">从光标位置</option>
          </select>
        </div>

        <!-- 结果显示 -->
        <div class="search-result" v-if="lastResult">
          <span :class="{ 'found': lastResult.found }">
            {{ lastResult.message }}
          </span>
        </div>
      </div>

      <div class="dialog-footer">
        <button class="primary-button" @click="handleFind">查找</button>
        <button class="primary-button" @click="handleFindNext" :disabled="!canFindNext">
          查找下一个
        </button>
        <button class="secondary-button" @click="$emit('close')">关闭</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

// 状态
const searchType = ref('text')
const searchText = ref('')
const searchHex = ref('')
const matchCase = ref(false)
const searchBackward = ref(false)
const wrapAround = ref(true)
const searchRange = ref('cursor')
const lastResult = ref(null)

// 计算属性
const canFindNext = computed(() => {
  return lastResult.value?.found && 
         (searchText.value || searchHex.value)
})

// 验证十六进制输入
function validateHex(event) {
  const input = event.target.value
  // 允许空格、数字和A-F字符
  const validHex = input.replace(/[^0-9A-Fa-f\s]/g, '')
  if (validHex !== input) {
    searchHex.value = validHex
  }
}

// 将搜索文本转换为字节数组
function getSearchBytes() {
  if (searchType.value === 'text') {
    return Array.from(searchText.value).map(char => char.charCodeAt(0))
  } else {
    return searchHex.value
      .trim()
      .split(/\s+/)
      .map(hex => parseInt(hex, 16))
      .filter(num => !isNaN(num))
  }
}

// 处理查找
async function handleFind() {
  const searchBytes = getSearchBytes()
  if (!searchBytes.length) {
    lastResult.value = {
      found: false,
      message: '请输入要查找的内容'
    }
    return
  }

  try {
    const options = {
      matchCase: matchCase.value,
      backward: searchBackward.value,
      wrapAround: wrapAround.value,
      range: searchRange.value
    }

    const result = await window.go.main.HexEditor.Search(searchBytes, options)
    
    if (result >= 0) {
      lastResult.value = {
        found: true,
        message: `找到位置: ${result.toString(16).toUpperCase()}`
      }
      emit('found', result)
    } else {
      lastResult.value = {
        found: false,
        message: '未找到匹配内容'
      }
    }
  } catch (err) {
    lastResult.value = {
      found: false,
      message: `搜索错误: ${err.message}`
    }
  }
}

// 查找下一个
function handleFindNext() {
  handleFind()
}

// 处理遮罩层点击
function handleOverlayClick(event) {
  if (event.target === event.currentTarget) {
    emit('close')
  }
}

// 定义事件
const emit = defineEmits(['close', 'found'])
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

.find-dialog {
  background-color: #252526;
  border: 1px solid #3D3D3D;
  border-radius: 4px;
  width: 400px;
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

.search-type,
.search-input,
.search-options,
.search-range {
  margin-bottom: 16px;
}

.radio-label,
.checkbox-label {
  display: inline-flex;
  align-items: center;
  margin-right: 16px;
  color: #D4D4D4;
  cursor: pointer;
}

input[type="radio"],
input[type="checkbox"] {
  margin-right: 6px;
}

.search-input label {
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

.search-result {
  margin-top: 16px;
  padding: 8px;
  background-color: #2D2D2D;
  border-radius: 3px;
}

.search-result .found {
  color: #4EC9B0;
}

.search-result :not(.found) {
  color: #F48771;
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

.secondary-button {
  background-color: #3D3D3D;
  color: #D4D4D4;
}

.secondary-button:hover {
  background-color: #4D4D4D;
}

.primary-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style> 