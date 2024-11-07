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
      this.openPorts.push({
        port: portInfo.port
      })
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
    }
  }
})