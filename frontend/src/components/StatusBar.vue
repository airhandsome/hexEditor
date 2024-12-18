<template>
  <div class="status-bar">
    <!-- 左侧信息区 -->
    <div class="status-left">
      <!-- 当前位置 -->
      <div class="status-item">
        <span class="label">位置:</span>
        <span class="value">{{ formatAddress(position) }}</span>
      </div>

      <!-- 选择范围 -->
      <div class="status-item" v-if="hasSelection">
        <span class="label">选择:</span>
        <span class="value">
          {{ formatAddress(selection.start) }} - {{ formatAddress(selection.end) }}
          ({{ selectionSize }} 字节)
        </span>
      </div>

      <!-- 文件大小 -->
      <div class="status-item">
        <span class="label">大小:</span>
        <span class="value">{{ formatFileSize(fileSize) }}</span>
      </div>
    </div>

    <!-- 右侧信息区 -->
    <div class="status-right">
      <!-- 修改状态 -->
      <div class="status-item" v-if="modified">
        <span class="modified-indicator">已修改</span>
      </div>

      <!-- 插入/覆盖模式 -->
      <div class="status-item">
        <span class="mode-indicator" :class="{ active: insertMode }">
          {{ insertMode ? '插入' : '覆盖' }}
        </span>
      </div>

      <!-- 字节序 -->
      <div class="status-item">
        <span class="endian-indicator">
          {{ littleEndian ? 'LE' : 'BE' }}
        </span>
      </div>

      <!-- 编码方式 -->
      <div class="status-item">
        <span class="encoding-indicator">
          {{ encoding }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// Props
const props = defineProps({
  position: {
    type: Number,
    default: 0
  },
  selection: {
    type: Object,
    default: () => ({ start: -1, end: -1 })
  },
  fileSize: {
    type: Number,
    default: 0
  },
  modified: {
    type: Boolean,
    default: false
  },
  insertMode: {
    type: Boolean,
    default: false
  },
  littleEndian: {
    type: Boolean,
    default: true
  },
  encoding: {
    type: String,
    default: 'UTF-8'
  }
})

// 计算属性
const hasSelection = computed(() => {
  return props.selection.start !== -1 && props.selection.end !== -1
})

const selectionSize = computed(() => {
  if (!hasSelection.value) return 0
  return props.selection.end - props.selection.start + 1
})

// 格式化函数
function formatAddress(addr) {
  return '0x' + addr.toString(16).padStart(8, '0').toUpperCase()
}

function formatFileSize(size) {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let value = size
  let unitIndex = 0

  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024
    unitIndex++
  }

  if (unitIndex === 0) {
    return `${value} ${units[unitIndex]}`
  } else {
    return `${value.toFixed(2)} ${units[unitIndex]} (${size} bytes)`
  }
}
</script>

<style scoped>
.status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 24px;
  background-color: #007ACC;
  color: #FFFFFF;
  font-size: 12px;
  font-family: 'Consolas', monospace;
  user-select: none;
}

.status-left, .status-right {
  display: flex;
  align-items: center;
  height: 100%;
}

.status-item {
  display: flex;
  align-items: center;
  padding: 0 8px;
  height: 100%;
  border-right: 1px solid rgba(255, 255, 255, 0.2);
}

.status-item:last-child {
  border-right: none;
}

.label {
  color: rgba(255, 255, 255, 0.8);
  margin-right: 4px;
}

.value {
  font-weight: 500;
}

.modified-indicator {
  color: #FFA500;
  font-weight: bold;
}

.mode-indicator {
  padding: 2px 6px;
  border-radius: 2px;
  background-color: rgba(255, 255, 255, 0.1);
}

.mode-indicator.active {
  background-color: rgba(255, 255, 255, 0.2);
}

.endian-indicator, .encoding-indicator {
  padding: 2px 6px;
  border-radius: 2px;
  background-color: rgba(255, 255, 255, 0.1);
}

/* 鼠标悬停效果 */
.status-item:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

/* 可点击项样式 */
.mode-indicator,
.endian-indicator,
.encoding-indicator {
  cursor: pointer;
}

.mode-indicator:hover,
.endian-indicator:hover,
.encoding-indicator:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

/* 动画效果 */
.modified-indicator {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { opacity: 1; }
  50% { opacity: 0.5; }
  100% { opacity: 1; }
}
</style> 