<template>
  <div class="toolbar">
    <!-- 文件操作组 -->
    <div class="toolbar-group">
      <button class="toolbar-button" @click="$emit('file-open')" title="打开文件 (Ctrl+O)">
        <i class="icon-file-open"></i>
        <span>打开</span>
      </button>
      <button class="toolbar-button" 
              @click="$emit('file-save')" 
              :disabled="!canSave"
              title="保存文件 (Ctrl+S)">
        <i class="icon-save"></i>
        <span>保存</span>
      </button>
    </div>

    <!-- 编辑操作组 -->
    <div class="toolbar-group">
      <button class="toolbar-button" 
              @click="$emit('undo')"
              :disabled="!canUndo"
              title="撤销 (Ctrl+Z)">
        <i class="icon-undo"></i>
      </button>
      <button class="toolbar-button" 
              @click="$emit('redo')"
              :disabled="!canRedo"
              title="重做 (Ctrl+Y)">
        <i class="icon-redo"></i>
      </button>
    </div>

    <!-- 查找操作组 -->
    <div class="toolbar-group">
      <button class="toolbar-button" @click="$emit('find')" title="查找 (Ctrl+F)">
        <i class="icon-search"></i>
        <span>查找</span>
      </button>
      <button class="toolbar-button" @click="$emit('goto')" title="跳转到 (Ctrl+G)">
        <i class="icon-goto"></i>
        <span>跳转</span>
      </button>
    </div>

    <!-- 视图操作组 -->
    <div class="toolbar-group">
      <button class="toolbar-button" 
              @click="$emit('toggle-file-panel')"
              :class="{ active: showFilePanel }"
              title="文件面板">
        <i class="icon-file-panel"></i>
      </button>
      <button class="toolbar-button" 
              @click="$emit('toggle-inspector')"
              :class="{ active: showInspector }"
              title="数据检查器">
        <i class="icon-inspector"></i>
      </button>
    </div>

    <!-- 数据操作组 -->
    <div class="toolbar-group">
      <select class="toolbar-select" v-model="selectedEncoding">
        <option value="hex">Hex</option>
        <option value="decimal">Decimal</option>
        <option value="binary">Binary</option>
      </select>
      <select class="toolbar-select" v-model="selectedEndian">
        <option value="little">Little Endian</option>
        <option value="big">Big Endian</option>
      </select>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

// Props
const props = defineProps({
  canSave: {
    type: Boolean,
    default: false
  },
  canUndo: {
    type: Boolean,
    default: false
  },
  canRedo: {
    type: Boolean,
    default: false
  },
  showFilePanel: {
    type: Boolean,
    default: true
  },
  showInspector: {
    type: Boolean,
    default: true
  }
})

// 状态
const selectedEncoding = ref('UTF-8')
const selectedEndian = ref('little')

// 事件
const emit = defineEmits([
  'file-open',
  'file-save',
  'undo',
  'redo',
  'find',
  'goto',
  'toggle-file-panel',
  'toggle-inspector',
  'encoding-change',
  'endian-change'
])

// 监听编码和字节序变化
watch(selectedEncoding, (newValue) => {
  emit('encoding-change', newValue)
})

watch(selectedEndian, (newValue) => {
  emit('endian-change', newValue)
})
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 8px;
  padding: 4px 8px;
  background-color: #2D2D2D;
  border-bottom: 1px solid #3D3D3D;
  user-select: none;
}

.toolbar-group {
  display: flex;
  gap: 4px;
  align-items: center;
  padding: 0 4px;
  border-right: 1px solid #3D3D3D;
}

.toolbar-group:last-child {
  border-right: none;
}

.toolbar-button {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background-color: transparent;
  border: 1px solid transparent;
  border-radius: 3px;
  color: #D4D4D4;
  cursor: pointer;
  font-size: 12px;
}

.toolbar-button:hover:not(:disabled) {
  background-color: #3D3D3D;
  border-color: #4D4D4D;
}

.toolbar-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.toolbar-button.active {
  background-color: #3D3D3D;
  border-color: #4D4D4D;
}

.toolbar-button i {
  width: 16px;
  height: 16px;
  background-size: contain;
  background-repeat: no-repeat;
}

.toolbar-select {
  padding: 2px 4px;
  background-color: #3D3D3D;
  border: 1px solid #4D4D4D;
  border-radius: 3px;
  color: #D4D4D4;
  font-size: 12px;
  outline: none;
}

.toolbar-select:focus {
  border-color: #569CD6;
}

/* 图标样式 */
.icon-file-open {
  background-image: url('../assets/icons/file-open.svg');
}

.icon-save {
  background-image: url('../assets/icons/save.svg');
}

.icon-undo {
  background-image: url('../assets/icons/undo.svg');
}

.icon-redo {
  background-image: url('../assets/icons/redo.svg');
}

.icon-search {
  background-image: url('../assets/icons/search.svg');
}

.icon-goto {
  background-image: url('../assets/icons/goto.svg');
}

.icon-file-panel {
  background-image: url('../assets/icons/file-panel.svg');
}

.icon-inspector {
  background-image: url('../assets/icons/inspector.svg');
}
</style> 