<template>
  <div class="scanner-component">
    <!-- 参数配置区域 -->
    <div class="input-group">
      <div class="input-item acrylic-input-box">
        <span class="input-label">目标URL</span>
        <el-input
          v-model="target"
          placeholder="URL地址 (例如: http://example.com)"
          clearable
          @clear="handleClear"
          :status="target && !validateUrl(target) ? 'error' : ''"
        >
          <template #append>
            <el-tooltip content="URL将自动添加末尾的/">
              <el-icon><InfoFilled /></el-icon>
            </el-tooltip>
          </template>
        </el-input>
      </div>
      
      <div class="input-item acrylic-input-box">
        <span class="input-label">字典文件</span>
        <el-button type="primary" @click="handleSelectFile">选择字典文件</el-button>
        <span v-if="selectedFile" class="selected-file">
          已选择: {{ selectedFile.name }}
        </span>
      </div>

      <div class="input-item acrylic-input-box">
        <span class="input-label">最大线程</span>
        <el-input
          v-model="maxThreads"
          placeholder="最大线程"
          type="number"
          :min="1"
          :max="10000"
        ></el-input>
      </div>
    </div>

    <!-- 进度信息和控制按钮区域 -->
    <div class="progress-container" v-if="store.showProgress">
      <div class="progress-info">
        <div class="status-group left">
          <div class="info-box acrylic-mini">
            <span class="status-text">{{ store.foundPaths.length }} 个有效路径</span>
          </div>
        </div>
        <div class="status-group right">
          <div class="info-box acrylic-mini">
            <span class="status-text">当前扫描速度：{{ formatSpeed(store.scanSpeed) }}个/s</span>
          </div>
          <div class="info-box acrylic-mini">
            <span class="status-text">{{ store.scannedPaths }}/{{ store.totalPaths }} 已扫描</span>
          </div>
          <el-button
            v-if="!store.isScanning"
            @click="handleScan"
            type="primary"
            class="scan-button"
          >
            扫描
          </el-button>
          <el-button
            v-else
            @click="handleStop"
            type="danger"
            class="scan-button"
          >
            停止
          </el-button>
        </div>
      </div>
      <div class="progress-wrapper acrylic-mini">
        <el-progress 
          :percentage="store.scanProgress" 
          :format="percentageFormat"
          :stroke-width="15"
          class="scan-progress"
        />
      </div>
    </div>

    <!-- 如果没有显示进度条，显示初始扫描按钮 -->
    <div v-else class="initial-scan-container">
      <el-button
        @click="handleScan"
        type="primary"
        class="scan-button"
      >
        开始扫描
      </el-button>
    </div>

    <!-- 分页控制器 -->
    <div class="pagination-container">
      <el-select v-model="pageSize" class="page-size-select" size="small">
        <el-option :value="10" label="10条/页" />
        <el-option :value="20" label="20条/页" />
        <el-option :value="50" label="50条/页" />
        <el-option :value="100" label="100条/页" />
      </el-select>
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="store.sortedPaths.length"
        layout="prev, pager, next, jumper"
        background
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 扫描结果表格 -->
    <el-table 
      :data="paginatedPaths" 
      style="width: 100%"
      :max-height="tableHeight"
      class="acrylic-effect"
      @sort-change="handleSortChange"
    >
      <el-table-column type="index" label="序号" width="60">
        <template #default="scope">
          {{ (currentPage - 1) * pageSize + scope.$index + 1 }}
        </template>
      </el-table-column>
      <!-- 路径列 -->
      <el-table-column 
        prop="fullUrl" 
        label="完整路径" 
        min-width="300"
        sortable="custom"
      >
        <template #default="scope">
          <el-link 
            type="primary" 
            @click="openInBrowser(scope.row.fullUrl)"
            style="cursor: pointer"
          >
            {{ scope.row.fullUrl }}
          </el-link>
        </template>
      </el-table-column>
      <!-- 状态码列 -->
      <el-table-column 
        prop="statusCode" 
        label="状态码" 
        width="100"
        sortable="custom"
      >
        <template #default="scope">
          <el-tag 
            :type="getStatusCodeType(scope.row.statusCode)"
            size="small"
          >
            {{ scope.row.statusCode }}
          </el-tag>
        </template>
      </el-table-column>
      <!-- 内容类型列 -->
      <el-table-column 
        prop="contentType" 
        label="内容类型" 
        width="200"
        sortable="custom"
      >
        <template #default="scope">
          <span>{{ formatContentType(scope.row.contentType) }}</span>
        </template>
      </el-table-column>
      <!-- 大小列 -->
      <el-table-column 
        prop="contentLength" 
        label="大小" 
        width="120"
        sortable="custom"
      >
        <template #default="scope">
          <span>{{ formatSize(scope.row.contentLength) }}</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { InfoFilled } from '@element-plus/icons-vue'
import { useDirsearchStore } from '../stores/dirsearchStore'

const store = useDirsearchStore()
const target = ref(localStorage.getItem('dirsearch_target') || '')
const selectedFile = ref(JSON.parse(localStorage.getItem('dirsearch_selected_file') || 'null'))
const maxThreads = ref(localStorage.getItem('dirsearch_max_threads') || '10')
const tableHeight = computed(() => window.innerHeight - 300)

// 分页相关的响应式变量
const currentPage = ref(1)
const pageSize = ref(10)

// 分页计算属性
const paginatedPaths = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return store.sortedPaths.slice(start, end)
})

// 分页处理方法
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val) => {
  currentPage.value = val
}

// 添加排序处理函数
const handleSortChange = ({ prop, order }) => {
  store.setSortConfig({ prop, order })
}

// 监听参数变化并保存
watch(target, (newVal) => {
  localStorage.setItem('dirsearch_target', newVal)
})

watch(selectedFile, (newVal) => {
  localStorage.setItem('dirsearch_selected_file', JSON.stringify(newVal))
})

watch(maxThreads, (newVal) => {
  localStorage.setItem('dirsearch_max_threads', newVal)
})

// URL 验证函数
const validateUrl = (url) => {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

// URL 规范化函数
const normalizeURL = (url) => {
  if (!url) return ''
  try {
    const parsed = new URL(url)
    return parsed.toString()
  } catch {
    return ''
  }
}

// 文件选择处理
const handleSelectFile = async () => {
  try {
    const filePath = await window.go.dirsearch.App.OpenFileDialog()
    if (filePath) {
      selectedFile.value = {
        name: filePath.split('\\').pop().split('/').pop(),
        path: filePath
      }
    }
  } catch (err) {
    console.error("文件选择错误:", err)
    ElMessage.error("文件选择失败: " + err.message)
  }
}

// 清除处理
const handleClear = () => {
  target.value = ''
  localStorage.removeItem('dirsearch_target')
}

// 格式化百分比显示
const percentageFormat = (percentage) => {
  return percentage === 100 ? '完成' : `${percentage.toFixed(2)}%`
}

// 格式化速度显示
const formatSpeed = (speed) => {
  if (speed === 0) return '0'
  return speed.toFixed(1)
}

// 获取状态码类型
const getStatusCodeType = (code) => {
  if (code >= 200 && code < 300) return 'success'
  if (code >= 300 && code < 400) return 'warning'
  if (code >= 400 && code < 500) return 'danger'
  if (code >= 500) return 'error'
  return 'info'
}

// 格式化内容类型
const formatContentType = (contentType) => {
  if (!contentType) return '未知'
  return contentType.split(';')[0]
}

// 格式化大小
const formatSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 在浏览器中打开URL
const openInBrowser = (url) => {
  window.runtime.BrowserOpenURL(url)
}

const handleScan = async () => {
  try {
    // 清理之前的事件监听
    window.runtime.EventsOff("path-found")
    window.runtime.EventsOff("dirsearch-status")
    window.runtime.EventsOff("dirsearch-progress")

    // 重置状态
    store.resetScan()
    store.setShowProgress(true)
    store.setIsScanning(true)

    // 绑定事件监听
    window.runtime.EventsOn("path-found", (pathInfo) => {
      store.addPath(pathInfo)
    })

    window.runtime.EventsOn("dirsearch-status", (status) => {
      store.setScanStatus(status)
      if (status === "completed") {
        ElMessage.success('扫描完成')
      } else if (status === "error") {
        ElMessage.error('扫描出错')
      } else if (status === "cancelled") {
        ElMessage.info('扫描已取消')
      }
    })

    window.runtime.EventsOn("dirsearch-progress", (progress) => {
      if (progress && typeof progress.current === 'number' && typeof progress.total === 'number') {
        store.setScannedPaths(progress.current)
        store.setTotalPaths(progress.total)
        store.setScanSpeed(progress.speed)
        console.log(`Progress update: ${progress.current}/${progress.total}, Speed: ${progress.speed}/s`)
      }
    })

    // 启动扫描
    await window.go.dirsearch.App.StartDirsearch(
      normalizeURL(target.value),
      selectedFile.value.path,
      parseInt(maxThreads.value)
    )
  } catch (err) {
    ElMessage.error('扫描出错: ' + err.message)
    store.setIsScanning(false)
    store.setScanStatus('error')
    
    // 清理事件监听
    window.runtime.EventsOff("path-found")
    window.runtime.EventsOff("dirsearch-status")
    window.runtime.EventsOff("dirsearch-progress")
  }
}

const handleStop = async () => {
  try {
    // 1. 立即移除所有事件监听，确保不再接收任何后端事件
    window.runtime.EventsOff("path-found")
    window.runtime.EventsOff("dirsearch-status")
    window.runtime.EventsOff("dirsearch-progress")
    window.runtime.EventsOff("dirsearch-error")  // 如果有错误事件监听也要移除

    // 2. 更新前端状态为停止中
    store.setIsScanning(false)
    store.setScanStatus('stopping')

    // 3. 通知后端停止扫描
    await window.go.dirsearch.App.StopDirsearch()
    
    // 4. 显示停止消息
    ElMessage.info('已停止扫描')

  } catch (err) {
    // 5. 如果出错，也要确保前端状态正确
    store.setIsScanning(false)
    store.setScanStatus('error')
    ElMessage.error('停止扫描失败: ' + err.message)
  } finally {
    // 6. 确保在任何情况下都重置进度显示
    store.setShowProgress(false)
  }
}

// 处理窗口大小变化
const handleResize = () => {
  tableHeight.value = window.innerHeight - 300
}

// 组件挂载时初始化
onMounted(() => {
  window.addEventListener('resize', handleResize)
  handleResize()
})

// 组件卸载时清理
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>
<style scoped>
.scanner-component {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 0 20px;
  gap: 20px;
  overflow: hidden;
}

.input-group {
  display: flex;
  gap: 16px;
  align-items: stretch;
}

.input-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.input-label {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
  text-align: center;
}

.selected-file {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.acrylic-input-box {
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  width: 100%;
}

.progress-container {
  margin: 10px 0;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  margin-bottom: 8px;
}

.status-group {
  display: flex;
  align-items: center;
  gap: 16px;
}

.status-group.left {
  margin-right: auto;
}

.status-group.right {
  margin-left: auto;
}

.info-box {
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 13px;
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
}

.status-text {
  color: #606266;
  font-weight: 500;
}

.progress-wrapper {
  padding: 10px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
}

.scan-progress {
  margin-bottom: 0;
}

.initial-scan-container {
  display: flex;
  justify-content: flex-end;
  margin: 10px 0;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 16px;
}

.page-size-select {
  width: 110px;
}

/* 表格基础样式 */
.el-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  flex: 1;
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
  .acrylic-input-box,
  .info-box,
  .progress-wrapper,
  .acrylic-effect {
    background: rgba(255, 255, 255, 0.15);
  }

  .input-label,
  .status-text {
    color: rgba(255, 255, 255, 0.9);
  }

  .selected-file {
    color: rgba(255, 255, 255, 0.6);
  }

  :deep(.el-table) {
    background-color: transparent;
    --el-table-border-color: rgba(255, 255, 255, 0.1);
    --el-table-header-bg-color: rgba(255, 255, 255, 0.05);
    --el-table-row-hover-bg-color: rgba(255, 255, 255, 0.08);
    --el-table-text-color: #e6e6e6;
    --el-table-header-text-color: #ffffff;
  }

  :deep(.el-pagination) {
    --el-pagination-button-bg-color: rgba(255, 255, 255, 0.15);
    --el-pagination-hover-color: var(--el-color-primary);
  }

  :deep(.el-select-dropdown__item) {
    color: #e6e6e6;
  }

  :deep(.el-select-dropdown__item.selected) {
    color: var(--el-color-primary);
  }
}

/* 按钮过渡效果 */
.scan-button {
  transition: all 0.3s ease;
}

.scan-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 毛玻璃效果类 */
.acrylic-mini {
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
  border-radius: 6px;
  transition: all 0.3s ease;
}

.acrylic-effect {
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  transition: all 0.3s ease;
}
</style>