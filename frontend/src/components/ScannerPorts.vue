<template>
  <div class="scanner-component">
    <!-- 参数配置区域 -->
    <div class="input-group">
      <div class="input-item acrylic-input-box">
        <span class="input-label">目标地址</span>
        <el-input
          v-model="target"
          placeholder="输入IP地址"
          clearable
          @clear="handleClear"
        ></el-input>
      </div>
      
      <div class="input-item acrylic-input-box">
        <span class="input-label">起始端口</span>
        <el-input
          v-model="startPort"
          placeholder="起始端口"
          type="number"
          :min="1"
          :max="65535"
        ></el-input>
      </div>

      <div class="input-item acrylic-input-box">
        <span class="input-label">结束端口</span>
        <el-input
          v-model="endPort"
          placeholder="结束端口"
          type="number"
          :min="1"
          :max="65535"
        ></el-input>
      </div>

      <div class="input-item acrylic-input-box">
        <span class="input-label">最大线程</span>
        <el-input
          v-model="maxThreads"
          placeholder="最大线程"
          type="number"
          :min="1"
          :max="1000"
        ></el-input>
      </div>
    </div>

    <!-- 进度信息和控制按钮区域 -->
    <div class="progress-container" v-if="showProgress">
      <div class="progress-info">
        <div class="status-group left">
          <div class="info-box acrylic-mini">
            <span class="status-text">{{ openPorts.length }} 个开放端口</span>
          </div>
        </div>
        <div class="status-group right">
          <div class="info-box acrylic-mini">
            <span class="status-text">{{ scannedPorts }}/{{ endPort }} 已扫描</span>
          </div>
          <el-button
            v-if="!scanning"
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
          :percentage="scanProgress" 
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

    <!-- 扫描结果表格 -->
    <el-table 
      :data="openPorts" 
      style="width: 100%"
      :max-height="500"
      class="acrylic-effect"
    >
      <el-table-column prop="port" label="端口" width="100" sortable />
      <el-table-column prop="service" label="服务" width="120">
        <template #default="scope">
          <div class="service-info">
            <span>{{ scope.row.service || '-' }}</span>
            <el-tag v-if="scope.row.tls" size="small" type="success">TLS</el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="product_name" label="产品名称" width="150">
        <template #default="scope">
          {{ scope.row.product_name || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="version" label="版本" width="120">
        <template #default="scope">
          {{ scope.row.version || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="info" label="详细信息" min-width="200">
        <template #default="scope">
          <el-tooltip 
            v-if="scope.row.info" 
            :content="scope.row.info" 
            placement="top" 
            :show-after="500"
          >
            <span class="truncate-text">{{ scope.row.info }}</span>
          </el-tooltip>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="更多信息" width="100">
        <template #default="scope">
          <el-popover
            placement="left"
            :width="300"
            trigger="click"
            v-if="hasAdditionalInfo(scope.row)"
          >
            <template #reference>
              <el-button type="primary" link>详情</el-button>
            </template>
            <div class="additional-info">
              <div v-if="scope.row.hostname" class="info-item">
                <span class="info-label">主机名:</span>
                <span class="info-value">{{ scope.row.hostname }}</span>
              </div>
              <div v-if="scope.row.operating_system" class="info-item">
                <span class="info-label">操作系统:</span>
                <span class="info-value">{{ scope.row.operating_system }}</span>
              </div>
              <div v-if="scope.row.device_type" class="info-item">
                <span class="info-label">设备类型:</span>
                <span class="info-value">{{ scope.row.device_type }}</span>
              </div>
              <div v-if="scope.row.probe_name" class="info-item">
                <span class="info-label">探针名称:</span>
                <span class="info-value">{{ scope.row.probe_name }}</span>
              </div>
            </div>
          </el-popover>
          <span v-else>-</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useScannerStore } from '../stores/scannerStore'

const store = useScannerStore()

// 计算属性
const target = computed({
  get: () => store.target,
  set: (value) => store.setTarget(value)
})

const startPort = computed({
  get: () => store.startPort,
  set: (value) => store.setStartPort(value)
})

const endPort = computed({
  get: () => store.endPort,
  set: (value) => store.setEndPort(value)
})

const maxThreads = computed({
  get: () => store.maxThreads,
  set: (value) => store.setMaxThreads(value)
})

const scanning = computed(() => store.isScanning)
const openPorts = computed(() => store.openPorts)
const showProgress = computed(() => store.showProgress)
const scanProgress = computed(() => store.scanProgress)
const scannedPorts = computed(() => store.scannedPorts)

// 方法
const percentageFormat = (percentage) => `${percentage}%`

const handleClear = () => store.setTarget('127.0.0.1')

const validateIP = (ip) => {
  const pattern = /^(\d{1,3}\.){3}\d{1,3}$/
  return pattern.test(ip) && ip.split('.').every(num => parseInt(num) >= 0 && parseInt(num) <= 255)
}

const hasAdditionalInfo = (row) => {
  return row.hostname || row.operating_system || row.device_type || row.probe_name
}

const handleStop = async () => {
  try {
    await window.go.main.App.StopScan()
    store.setIsScanning(false)
    window.runtime.EventsOff("port-found")
    window.runtime.EventsOff("scan-status")
    window.runtime.EventsOff("scan-progress")
    ElMessage.info('已停止扫描')
  } catch (err) {
    ElMessage.error('停止扫描失败: ' + err.message)
  }
}

const handleScan = async () => {
  if (!target.value || !validateIP(target.value)) {
    ElMessage.error('请输入有效的IP地址')
    return
  }

  const start = parseInt(startPort.value)
  const end = parseInt(endPort.value)
  const threads = parseInt(maxThreads.value)
  
  if (start < 1 || start > 65535 || end < 1 || end > 65535) {
    ElMessage.error('端口号必须在 1-65535 之间')
    return
  }

  if (start > end) {
    ElMessage.error('起始端口不能大于结束端口')
    return
  }

  if (threads < 1 || threads > 1000) {
    ElMessage.error('线程数必须在 1-1000 之间')
    return
  }

  try {
    // 清理之前的事件监听
    window.runtime.EventsOff("port-found")
    window.runtime.EventsOff("scan-status")
    window.runtime.EventsOff("scan-progress")

    // 重置状态
    store.resetScan()
    store.setShowProgress(true)
    store.setIsScanning(true)

    // 绑定事件监听
    window.runtime.EventsOn("port-found", (portInfo) => {
      store.addPort(portInfo)
    })

    window.runtime.EventsOn("scan-status", (status) => {
      if (status === "completed") {
        store.setScanComplete(true)
        ElMessage.success('扫描完成')
      } else if (status === "error") {
        store.setIsScanning(false)
        store.setScanComplete(false)
      }
    })

    window.runtime.EventsOn("scan-progress", (progress) => {
      store.setScannedPorts(progress.current_port)
    })

    // 启动扫描
    await window.go.main.App.ScanPorts(
      target.value,
      start,
      end,
      threads
    )
  } catch (err) {
    ElMessage.error('扫描出错: ' + err.message)
    store.setIsScanning(false)
    store.setScanComplete(false)
    
    // 清理事件监听
    window.runtime.EventsOff("port-found")
    window.runtime.EventsOff("scan-status")
    window.runtime.EventsOff("scan-progress")
  }
}
</script>

<style scoped>
.scanner-component {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 0 20px;
  gap: 20px;
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

.scan-progress :deep(.el-progress-bar__outer) {
  border-radius: 4px;
  background-color: rgba(255, 255, 255, 0.2);
}

.scan-progress :deep(.el-progress-bar__inner) {
  border-radius: 4px;
}

.initial-scan-container {
  display: flex;
  justify-content: flex-end;
  margin: 10px 0;
}

.service-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.truncate-text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 200px;
  display: inline-block;
}

.additional-info {
  padding: 8px;
}

.info-item {
  margin-bottom: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  color: #909399;
  font-size: 13px;
}

.info-value {
  color: #303133;
  font-size: 14px;
  word-break: break-all;
}

/* 表格样式 */
.el-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

:deep(.el-table__inner-wrapper) {
  border-radius: 12px;
}

:deep(.el-table__header) {
  border-top-left-radius: 12px;
  border-top-right-radius: 12px;
  overflow: hidden;
}

:deep(.el-table__body) {
  border-bottom-left-radius: 12px;
  border-bottom-right-radius: 12px;
  overflow: hidden;
}

:deep(.el-table__row) {
  background-color: transparent !important;
}

:deep(.el-table__row:hover) {
  background-color: rgba(255, 255, 255, 0.1) !important;
}

:deep(.el-table__header-wrapper) {
  background-color: rgba(255, 255, 255, 0.1);
}

:deep(.el-table__cell) {
  background-color: transparent !important;
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

  .input-item :deep(.el-input__wrapper) {
    background-color: #1a1a1a;
    box-shadow: 0 0 0 1px #4c4d4f inset !important;
  }

  .input-item :deep(.el-input__inner) {
    color: #ffffff;
  }

  .empty-text {
    color: rgba(255, 255, 255, 0.6);
  }

  .info-label {
    color: #a6a6a6;
  }

  .info-value {
    color: #e6e6e6;
  }

  .el-table {
    --el-table-border-color: rgba(255, 255, 255, 0.1);
    --el-table-header-bg-color: rgba(255, 255, 255, 0.05);
    --el-table-row-hover-bg-color: rgba(255, 255, 255, 0.08);
    --el-table-text-color: #e6e6e6;
    --el-table-header-text-color: #ffffff;
  }

  :deep(.el-table__header-wrapper) {
    background-color: rgba(255, 255, 255, 0.05);
  }

  :deep(.el-table__row:hover) {
    background-color: rgba(255, 255, 255, 0.08) !important;
  }

  :deep(.el-table) {
    background-color: transparent;
  }

  :deep(.el-table__cell) {
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  :deep(.el-table__header th.el-table__cell) {
    background-color: rgba(255, 255, 255, 0.05);
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  :deep(.el-table--enable-row-hover .el-table__body tr:hover > td.el-table__cell) {
    background-color: rgba(255, 255, 255, 0.08);
  }
}
</style>

