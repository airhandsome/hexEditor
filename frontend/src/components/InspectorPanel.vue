<template>
  <div class="inspector-panel">
    <!-- 面板标题 -->
    <div class="panel-header">
      <span>数据检查器</span>
      <button class="close-button" @click="$emit('close')">&times;</button>
    </div>

    <!-- 数据类型显示 -->
    <div class="inspector-content">
      <!-- 当前位置信息 -->
      <div class="position-info">
        <span>位置: {{ formatAddress(offset) }}</span>
        <span>字节: {{ formatByte(currentByte) }}</span>
      </div>

      <!-- 整数类型 -->
      <div class="type-section">
        <h3>整数</h3>
        <table>
          <tr>
            <td>Signed Byte:</td>
            <td>{{ formatSignedByte(currentByte) }}</td>
          </tr>
          <tr>
            <td>Unsigned Byte:</td>
            <td>{{ currentByte }}</td>
          </tr>
          <tr>
            <td>Signed Short:</td>
            <td>{{ formatSignedShort(getShort()) }}</td>
          </tr>
          <tr>
            <td>Unsigned Short:</td>
            <td>{{ getShort() >>> 0 }}</td>
          </tr>
          <tr>
            <td>Signed Int:</td>
            <td>{{ formatSignedInt(getInt()) }}</td>
          </tr>
          <tr>
            <td>Unsigned Int:</td>
            <td>{{ getInt() >>> 0 }}</td>
          </tr>
          <tr>
            <td>Signed Int64:</td>
            <td>{{ formatInt64(getInt64()) }}</td>
          </tr>
          <tr>
            <td>Unsigned Int64:</td>
            <td>{{ formatUInt64(getInt64()) }}</td>
          </tr>
        </table>
      </div>

      <!-- 浮点类型 -->
      <div class="type-section">
        <h3>浮点数</h3>
        <table>
          <tr>
            <td>Float:</td>
            <td>{{ formatFloat(getFloat()) }}</td>
          </tr>
          <tr>
            <td>Double:</td>
            <td>{{ formatDouble(getDouble()) }}</td>
          </tr>
        </table>
      </div>

      <!-- 时间类型 -->
      <div class="type-section">
        <h3>时间</h3>
        <table>
          <tr>
            <td>DOS Time:</td>
            <td>{{ formatDosTime(getInt()) }}</td>
          </tr>
          <tr>
            <td>DOS Date:</td>
            <td>{{ formatDosDate(getInt()) }}</td>
          </tr>
          <tr>
            <td>Unix Time:</td>
            <td>{{ formatUnixTime(getInt()) }}</td>
          </tr>
          <tr>
            <td>File Time:</td>
            <td>{{ formatFileTime(getInt64()) }}</td>
          </tr>
        </table>
      </div>

      <!-- 字符串类型 -->
      <div class="type-section">
        <h3>字符串</h3>
        <table>
          <tr>
            <td>ASCII:</td>
            <td>{{ formatAsciiString(8) }}</td>
          </tr>
          <tr>
            <td>UTF-16:</td>
            <td>{{ formatUtf16String(8) }}</td>
          </tr>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// Props
const props = defineProps({
  currentByte: {
    type: Number,
    default: 0
  },
  offset: {
    type: Number,
    default: 0
  },
  hexData: {
    type: Array,
    default: () => []
  },
  littleEndian: {
    type: Boolean,
    default: true
  }
})

// 格式化函数
function formatAddress(addr) {
  return '0x' + addr.toString(16).padStart(8, '0').toUpperCase()
}

function formatByte(byte) {
  return byte?.toString(16).padStart(2, '0').toUpperCase() ?? '00'
}

function formatSignedByte(byte) {
  return byte >= 128 ? byte - 256 : byte
}

function formatSignedShort(value) {
  return value >= 32768 ? value - 65536 : value
}

function formatSignedInt(value) {
  return value | 0
}

function formatInt64(value) {
  return BigInt(value).toString()
}

function formatUInt64(value) {
  return BigInt.asUintN(64, BigInt(value)).toString()
}

function formatFloat(value) {
  return value.toFixed(6)
}

function formatDouble(value) {
  return value.toFixed(15)
}

// 数据获取函数
function getBytes(count) {
  const bytes = []
  for (let i = 0; i < count; i++) {
    bytes.push(props.hexData[props.offset + i] || 0)
  }
  return props.littleEndian ? bytes : bytes.reverse()
}

function getShort() {
  const bytes = getBytes(2)
  return (bytes[1] << 8) | bytes[0]
}

function getInt() {
  const bytes = getBytes(4)
  return (bytes[3] << 24) | (bytes[2] << 16) | (bytes[1] << 8) | bytes[0]
}

function getInt64() {
  const bytes = getBytes(8)
  let value = BigInt(0)
  for (let i = 0; i < 8; i++) {
    value = (value << BigInt(8)) | BigInt(bytes[7 - i])
  }
  return value
}

function getFloat() {
  const bytes = getBytes(4)
  const buffer = new ArrayBuffer(4)
  const view = new DataView(buffer)
  bytes.forEach((byte, i) => view.setUint8(i, byte))
  return view.getFloat32(0, props.littleEndian)
}

function getDouble() {
  const bytes = getBytes(8)
  const buffer = new ArrayBuffer(8)
  const view = new DataView(buffer)
  bytes.forEach((byte, i) => view.setUint8(i, byte))
  return view.getFloat64(0, props.littleEndian)
}

// 时间格式化函数
function formatDosTime(time) {
  const hour = (time >> 11) & 0x1f
  const minute = (time >> 5) & 0x3f
  const second = (time & 0x1f) * 2
  return `${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}:${second.toString().padStart(2, '0')}`
}

function formatDosDate(date) {
  const year = ((date >> 9) & 0x7f) + 1980
  const month = (date >> 5) & 0x0f
  const day = date & 0x1f
  return `${year}-${month.toString().padStart(2, '0')}-${day.toString().padStart(2, '0')}`
}

function formatUnixTime(timestamp) {
  return new Date(timestamp * 1000).toISOString()
}

function formatFileTime(fileTime) {
  // Windows FILETIME to JavaScript Date
  const WINDOWS_EPOCH = 116444736000000000n // Difference between Windows FILETIME and Unix Epoch
  const unixTime = Number((fileTime - WINDOWS_EPOCH) / 10000n) // Convert to milliseconds
  return new Date(unixTime).toISOString()
}

// 字符串格式化函数
function formatAsciiString(length) {
  return getBytes(length)
    .map(byte => byte >= 32 && byte <= 126 ? String.fromCharCode(byte) : '.')
    .join('')
}

function formatUtf16String(length) {
  const bytes = getBytes(length * 2)
  const chars = []
  for (let i = 0; i < bytes.length; i += 2) {
    const charCode = (bytes[i + 1] << 8) | bytes[i]
    chars.push(String.fromCharCode(charCode))
  }
  return chars.join('')
}
</script>

<style scoped>
.inspector-panel {
  width: 300px;
  background-color: #252526;
  border-left: 1px solid #3D3D3D;
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

.inspector-content {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.position-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  padding: 8px;
  background-color: #2D2D2D;
  border-radius: 3px;
  font-family: 'Consolas', monospace;
}

.type-section {
  margin-bottom: 20px;
}

.type-section h3 {
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
  font-family: 'Consolas', monospace;
}

td:first-child {
  color: #9CDCFE;
  width: 120px;
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