<template>
  <div class="jsfinder-container">
    <h1>Gitdorker</h1>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const url = ref('')
const results = ref([])

const startScan = async () => {
  try {
    // 这里调用后端的扫描方法
    const response = await window.go.main.ScanJS(url.value)
    results.value = response
  } catch (error) {
    console.error('扫描出错:', error)
  }
}
</script>

<style scoped>
.jsfinder-container {
  padding: 20px;
}

.input-section {
  margin: 20px 0;
}

.input-section input {
  padding: 8px;
  margin-right: 10px;
  width: 300px;
}

.input-section button {
  padding: 8px 16px;
}

.results-section {
  margin-top: 20px;
}
</style>