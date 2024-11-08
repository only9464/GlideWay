import { defineStore } from 'pinia'

export const useDirsearchStore = defineStore('dirsearch', {
  state: () => ({
    foundPaths: [],      
    scannedPaths: 0,     
    totalPaths: 0,       
    showProgress: false, 
    isScanning: false,   
    sortConfig: {        
      prop: null,
      order: null
    },
    scanStatus: 'idle'   
  }),
  
  getters: {
    scanProgress: (state) => {
      if (state.totalPaths <= 0) return 0
      const progress = Math.round((state.scannedPaths / state.totalPaths) * 100)
      return Math.min(progress, 99)
    },
    
    sortedPaths: (state) => {
      if (!state.sortConfig.prop || !state.sortConfig.order) {
        return state.foundPaths
      }
      
      return [...state.foundPaths].sort((a, b) => {
        const prop = state.sortConfig.prop
        const order = state.sortConfig.order
        
        if (a[prop] === b[prop]) return 0
        
        const result = a[prop] > b[prop] ? 1 : -1
        return order === 'ascending' ? result : -result
      })
    }
  },
  
  actions: {
    resetScan() {
      this.foundPaths = []
      this.scannedPaths = 0
      this.totalPaths = 0
      this.isScanning = false
      this.scanStatus = 'idle'
    },
    
    setIsScanning(value) {
      this.isScanning = value
      if (!value) {
        if (this.scanStatus === 'scanning') {
          this.scanStatus = 'cancelled'
        }
      } else {
        this.scanStatus = 'scanning'
        this.showProgress = true
      }
    },

    setScanStatus(status) {
      this.scanStatus = status
      if (status === 'completed' || status === 'cancelled' || status === 'error') {
        this.isScanning = false
      }
    },

    setShowProgress(value) {
      if (value) {
        this.showProgress = true
      }
    },
    
    setScannedPaths(value) {
      if (typeof value === 'number') {
        // 确保扫描数量至少等于找到的路径数量
        this.scannedPaths = Math.max(value, this.foundPaths.length)
      }
    },

    setTotalPaths(value) {
      if (typeof value === 'number' && value > 0) {
        this.totalPaths = value
      }
    },

    clearFoundPaths() {
      this.foundPaths = []
    },

    addPath(pathInfo) {
      this.foundPaths.push({
        path: pathInfo.path,
        fullUrl: pathInfo.fullUrl,
        statusCode: pathInfo.statusCode,
        contentType: pathInfo.contentType,
        contentLength: pathInfo.contentLength,
      })
      // 确保扫描数量至少等于找到的路径数量
      this.scannedPaths = Math.max(this.scannedPaths, this.foundPaths.length)
    },

    setSortConfig({ prop, order }) {
      if (!prop || !order) {
        this.sortConfig.prop = null
        this.sortConfig.order = null
        return
      }
      this.sortConfig.prop = prop
      this.sortConfig.order = order
    },

    setComplete() {
      this.isScanning = false
      this.scanStatus = 'completed'
    },

    exportResults() {
      return {
        timestamp: new Date().toISOString(),
        totalScanned: this.scannedPaths,
        foundPaths: this.foundPaths
      }
    }
  }
})