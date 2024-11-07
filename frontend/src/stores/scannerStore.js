import { defineStore } from 'pinia'

export const useScannerStore = defineStore('scanner', {
  state: () => ({
    openPorts: [],
    scannedPorts: 0,
    showProgress: false,
    scanComplete: false,
    target: '127.0.0.1',
    startPort: 1,
    endPort: 65535,
    maxThreads: 500,
    isScanning: false
  }),
  
  getters: {
    scanProgress: (state) => {
      if (state.scanComplete) return 100
      if (!state.showProgress) return 0
      
      const total = state.endPort - state.startPort + 1
      const scanned = state.scannedPorts - state.startPort + 1
      return Math.min(Math.round((scanned / total) * 100), 99)
    }
  },
  
  actions: {
    resetScan() {
      this.openPorts = []
      this.scannedPorts = 0
      this.showProgress = false
      this.scanComplete = false
      this.isScanning = false
    },
    
    setTarget(value) {
      this.target = value || '127.0.0.1'
    },
    
    setStartPort(value) {
      const port = parseInt(value)
      if (port >= 1 && port <= 65535) {
        this.startPort = port
      }
    },
    
    setEndPort(value) {
      const port = parseInt(value)
      if (port >= 1 && port <= 65535) {
        this.endPort = port
      }
    },
    
    setMaxThreads(value) {
      const threads = parseInt(value)
      if (threads >= 1 && threads <= 1000) {
        this.maxThreads = threads
      }
    },
    
    setIsScanning(value) {
      this.isScanning = value
    },
    
    addPort(portInfo) {
      // 扩展端口信息，包含所有指纹识别结果
      this.openPorts.push({
        port: portInfo.port,
        protocol: portInfo.protocol,
        service: portInfo.service,
        product_name: portInfo.product_name,
        version: portInfo.version,
        info: portInfo.info,
        hostname: portInfo.hostname,
        operating_system: portInfo.operating_system,
        device_type: portInfo.device_type,
        probe_name: portInfo.probe_name,
        tls: portInfo.tls
      })
      // 按端口号排序
      this.openPorts.sort((a, b) => a.port - b.port)
    },
    
    setScanComplete(value) {
      this.scanComplete = value
      if (value) {
        this.isScanning = false
        this.scannedPorts = this.endPort
      }
    },
    
    setShowProgress(value) {
      this.showProgress = value
    },
    
    setScannedPorts(value) {
      this.scannedPorts = value
    },
    
    clearAll() {
      this.resetScan()
      this.target = '127.0.0.1'
      this.startPort = 1
      this.endPort = 65535
      this.maxThreads = 500
    },

    // 新增：检查端口是否有额外信息
    hasAdditionalInfo(port) {
      return !!(port.hostname || port.operating_system || port.device_type || port.probe_name)
    },

    // 新增：获取端口的服务描述
    getServiceDescription(port) {
      const parts = []
      if (port.service) parts.push(port.service)
      if (port.product_name) parts.push(port.product_name)
      if (port.version) parts.push(port.version)
      return parts.join(' - ') || '未知服务'
    },

    // 新增：导出扫描结果
    exportResults() {
      return this.openPorts.map(port => ({
        port: port.port,
        service: this.getServiceDescription(port),
        details: {
          protocol: port.protocol,
          tls: port.tls,
          info: port.info,
          hostname: port.hostname,
          operating_system: port.operating_system,
          device_type: port.device_type,
          probe_name: port.probe_name
        }
      }))
    }
  }
})