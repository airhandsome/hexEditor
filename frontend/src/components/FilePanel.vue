<template>
  <div class="file-panel">
    <!-- 面板标题 -->
    <div class="panel-header">
      <span>文件信息</span>
      <button class="close-button" @click="$emit('close')">&times;</button>
    </div>

    <!-- 文件信息 -->
    <div class="file-info" v-if="fileName">
      <div class="info-section">
        <h3>基本信息</h3>
        <table>
          <tr>
            <td>文件名:</td>
            <td>{{ fileName }}</td>
          </tr>
          <tr>
            <td>大小:</td>
            <td>{{ formatFileSize(fileSize) }}</td>
          </tr>
          <tr>
            <td>路径:</td>
            <td class="file-path">{{ filePath }}</td>
          </tr>
        </table>
      </div>

      <!-- 文件统计 -->
      <div class="info-section">
        <h3>统计信息</h3>
        <table>
          <tr>
            <td>字节数:</td>
            <td>{{ fileSize }}</td>
          </tr>
          <tr>
            <td>页数:</td>
            <td>{{ Math.ceil(fileSize / pageSize) }}</td>
          </tr>
          <tr>
            <td>修改字节:</td>
            <td>{{ modifiedCount }}</td>
          </tr>
        </table>
      </div>

      <!-- 文件结构 -->
      <div class="info-section">
        <h3>文件结构</h3>
        <div class="file-structure">
          <div class="structure-item" 
               v-for="section in fileStructure" 
               :key="section.offset"
               @click="handleSectionClick(section)">
            <span class="section-name">{{ section.name }}</span>
            <span class="section-offset">{{ formatAddress(section.offset) }}</span>
            <span class="section-size">{{ formatFileSize(section.size) }}</span>
          </div>
        </div>
      </div>

      <!-- 书签列表 -->
      <div class="info-section">
        <h3>书签</h3>
        <div class="bookmarks">
          <div class="bookmark-item" 
               v-for="bookmark in bookmarks" 
               :key="bookmark.offset"
               @click="handleBookmarkClick(bookmark)">
            <span class="bookmark-name">{{ bookmark.name }}</span>
            <span class="bookmark-offset">{{ formatAddress(bookmark.offset) }}</span>
          </div>
          <div v-if="bookmarks.length === 0" class="empty-message">
            无书签
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <p>未打开文件</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// Props
const props = defineProps({
  fileName: {
    type: String,
    default: ''
  },
  filePath: {
    type: String,
    default: ''
  },
  fileSize: {
    type: Number,
    default: 0
  },
  pageSize: {
    type: Number,
    default: 4096
  },
  modifiedBytes: {
    type: Set,
    default: () => new Set()
  },
  fileStructure: {
    type: Array,
    default: () => []
  },
  bookmarks: {
    type: Array,
    default: () => []
  }
})

// Emits
const emit = defineEmits(['close', 'section-click', 'bookmark-click'])

// 计算属性
const modifiedCount = computed(() => props.modifiedBytes.size)

// 格式化函数
function formatFileSize(size) {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`
  if (size < 1024 * 1024 * 1024) return `${(size / 1024 / 1024).toFixed(2)} MB`
  return `${(size / 1024 / 1024 / 1024).toFixed(2)} GB`
}

function formatAddress(addr) {
  return '0x' + addr.toString(16).padStart(8, '0').toUpperCase()
}

// 事件处理
function handleSectionClick(section) {
  emit('section-click', section)
}

function handleBookmarkClick(bookmark) {
  emit('bookmark-click', bookmark)
}
</script>

<style scoped>
.file-panel {
  width: 250px;
  background-color: #252526;
  border-right: 1px solid #3D3D3D;
  display: flex;
  flex-direction: column;
  user-select: none;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background-color: #2D2D2D;
  border-bottom: 1px solid #3D3D3D;
  font-weight: bold;
}

.close-button {
  background: none;
  border: none;
  color: #D4D4D4;
  cursor: pointer;
  font-size: 16px;
  padding: 0 4px;
}

.close-button:hover {
  color: #FFFFFF;
}

.file-info {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.info-section {
  margin-bottom: 20px;
}

.info-section h3 {
  font-size: 12px;
  color: #D4D4D4;
  margin-bottom: 8px;
  padding-bottom: 4px;
  border-bottom: 1px solid #3D3D3D;
}

table {
  width: 100%;
  font-size: 12px;
  border-collapse: collapse;
}

td {
  padding: 4px 0;
  color: #D4D4D4;
}

td:first-child {
  color: #9CDCFE;
  width: 80px;
}

.file-path {
  word-break: break-all;
}

.file-structure, .bookmarks {
  font-size: 12px;
}

.structure-item, .bookmark-item {
  display: flex;
  justify-content: space-between;
  padding: 4px 8px;
  cursor: pointer;
  border-radius: 3px;
}

.structure-item:hover, .bookmark-item:hover {
  background-color: #2D2D2D;
}

.section-name, .bookmark-name {
  flex: 1;
  color: #CE9178;
}

.section-offset, .bookmark-offset {
  color: #9CDCFE;
  margin-left: 8px;
}

.section-size {
  color: #B5CEA8;
  margin-left: 8px;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #6D6D6D;
  font-style: italic;
}

.empty-message {
  color: #6D6D6D;
  font-style: italic;
  text-align: center;
  padding: 8px;
}

::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #1E1E1E;
}

::-webkit-scrollbar-thumb {
  background: #424242;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #4D4D4D;
}
</style> 