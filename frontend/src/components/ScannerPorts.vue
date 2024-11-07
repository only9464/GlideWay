<template>
  <div class="scanner-component">
    <div class="input-group">
      <el-input
        v-model="target"
        placeholder="输入IP地址"
        class="acrylic-effect"
        clearable
        @clear="handleClear"
      ></el-input>
      <el-button
        @click="handleScan"
        :loading="scanning"
        type="primary"
        class="scan-button"
      >
        {{ scanning ? '扫描中...' : '扫描' }}
      </el-button>
    </div>

    <!-- 修改进度条显示条件 -->
    <div class="progress-container" v-if="showProgress">
      <div class="progress-info">
        {{ openPorts.length }} 个开放端口 / {{ scannedPorts }}/65535 已扫描
      </div>
      <el-progress 
        :percentage="scanProgress" 
        :format="percentageFormat"
        :stroke-width="15"
        class="scan-progress"
      />
    </div>

    <el-table 
      :data="openPorts" 
      class="acrylic-effect"
      :default-sort="{ prop: 'port', order: 'ascending' }"
    >
      <el-table-column 
        label="端口" 
        prop="port" 
        sortable 
        width="100"
      ></el-table-column>
      <el-table-column 
        prop="service" 
        label="服务"
      ></el-table-column>
      <el-table-column 
        prop="banner" 
        label="服务信息"
      >
        <template #default="scope">
          <el-tooltip 
            v-if="scope.row.banner" 
            :content="scope.row.banner" 
            placement="top"
          >
            <span class="banner-text">{{ scope.row.banner.substring(0, 50) }}</span>
          </el-tooltip>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column 
        prop="protocol" 
        label="协议"
      ></el-table-column>
      
      <!-- 添加空数据提示 -->
      <template #empty>
        <div class="empty-text">
          暂无扫描数据
        </div>
      </template>
    </el-table>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useScannerStore } from '../stores/scannerStore'

const store = useScannerStore()
const scanning = ref(false)
const totalPorts = 65535

// 使用计算属性来双向绑定 target
const target = computed({
  get: () => store.target,
  set: (value) => store.setTarget(value)
})

// 使用 store 中的数据
const openPorts = computed(() => store.openPorts)
const scannedPorts = computed(() => store.scannedPorts)
const showProgress = computed(() => store.showProgress)
const scanComplete = computed(() => store.scanComplete)

const scanProgress = computed(() => {
  if (store.scanComplete) return 100
  if (store.openPorts.length === 0) return 0
  const maxPort = Math.max(...store.openPorts.map(port => port.port))
  return Math.round((maxPort / totalPorts) * 100)
})

const percentageFormat = (percentage) => {
  return `${percentage}%`
}

function handleClear() {
  store.setTarget('')
}

function validateIP(ip) {
  const pattern = /^(\d{1,3}\.){3}\d{1,3}$/
  if (!pattern.test(ip)) return false
  return ip.split('.').every(num => parseInt(num) >= 0 && parseInt(num) <= 255)
}

async function handleScan() {
  if (!target.value) {
    ElMessage.error('请输入IP地址')
    return
  }

  if (!validateIP(target.value)) {
    ElMessage.error('请输入有效的IP地址')
    return
  }

  try {
    // 清理之前的事件监听器
    window.runtime.EventsOff("port-found")
    window.runtime.EventsOff("scan-complete")

    scanning.value = true
    store.resetScan()
    store.setShowProgress(true)

    // 监听端口发现事件
    window.runtime.EventsOn("port-found", (portInfo) => {
      store.addPort(portInfo)
    })

    // 监听扫描完成事件
    window.runtime.EventsOn("scan-complete", () => {
      scanning.value = false
      store.setScanComplete(true)
      ElMessage.success('扫描完成')
      // 清理事件监听
      window.runtime.EventsOff("port-found")
      window.runtime.EventsOff("scan-complete")
    })

    // 开始扫描
    await window.go.main.App.ScanPorts(target.value)

  } catch (err) {
    ElMessage.error('扫描出错: ' + err.message)
    scanning.value = false
    store.setScanComplete(false)
    // 清理事件监听
    window.runtime.EventsOff("port-found")
    window.runtime.EventsOff("scan-complete")
  }
}
</script>

<style scoped>
.scanner-component {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: 0 20px 0 20px;
}

.progress-container {
  margin: 10px 0;
}

.progress-info {
  margin-bottom: 8px;
  color: #606266;
  font-size: 14px;
}

.scan-progress {
  margin-bottom: 20px;
}

.input-group {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.input-group .el-input {
  flex: 1;
}

.banner-text {
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 300px;
  display: inline-block;
}

.acrylic-effect {
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.empty-text {
  padding: 40px 0;
  color: #909399;
  font-size: 14px;
}
</style>