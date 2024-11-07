import { defineStore } from 'pinia'

export const useScannerStore = defineStore('scanner', {
  state: () => ({
    openPorts: [],
    scannedPorts: 0,
    showProgress: false,
    scanComplete: false,
    target: ''  // 添加 target 字段
  }),
  actions: {
    resetScan() {
      this.openPorts = []
      this.scannedPorts = 0
      this.showProgress = false
      this.scanComplete = false
      // 不重置 target，这样可以保留上次的 IP
    },
    setTarget(value) {
      this.target = value
    },
    addPort(portInfo) {
      this.openPorts.push(portInfo)
      this.scannedPorts++
    },
    setScanComplete(value) {
      this.scanComplete = value
    },
    setShowProgress(value) {
      this.showProgress = value
    }
  }
})