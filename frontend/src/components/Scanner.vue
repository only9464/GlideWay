<script setup>
import { ref } from 'vue'
// 引入动态调用的接口
import { ExecuteFunction } from '../../wailsjs/go/main/App'

const target = ref('')
const scanning = ref(false)
const openPorts = ref([])
const error = ref('')

function validateIP(ip) {
  const ipPattern = /^(\d{1,3}\.){3}\d{1,3}$/
  if (!ipPattern.test(ip)) {
    return false
  }
  const parts = ip.split('.')
  return parts.every(part => {
    const num = parseInt(part)
    return num >= 0 && num <= 255
  })
}

async function handleScan() {
  error.value = ''
  openPorts.value = []

  if (!target.value) {
    error.value = '请输入你要干爆的IP地址'
    return
  }

  if (!validateIP(target.value)) {
    error.value = '请输入有效的IIIIIIIIIIIIP地址'
    return
  }

  try {
    scanning.value = true
    console.log(target.value)
    const ports = await ExecuteFunction('ScanPorts', [target.value])
    openPorts.value = ports.sort((a, b) => a - b)
  } catch (err) {
    error.value = '扫描出错: ' + err.message
  } finally {
    scanning.value = false
  }
}

</script>

<template>
  <div class="scanner-component">
    <div class="input-group mb-4">
      <input
        v-model="target"
        type="text"
        class="border p-2 rounded mr-2"
        placeholder="输入NIDE IP地址"
      />
      <button
        @click="handleScan"
        :disabled="scanning"
        class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 disabled:bg-gray-400"
      >
        {{ scanning ? '扫描中...' : '扫描' }}
      </button>
    </div>

    <div v-if="error" class="text-red-500 mb-4">
      {{ error }}
    </div>

    <div v-if="openPorts.length > 0" class="results">
      <h3 class="text-lg font-bold mb-2">开放的端口:</h3>
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2">
        <div
          v-for="port in openPorts"
          :key="port"
          class="bg-green-100 p-2 rounded text-center"
        >
          {{ port }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 可根据需要添加样式 */
</style>
