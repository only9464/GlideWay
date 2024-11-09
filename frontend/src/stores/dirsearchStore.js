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
    scanStatus: 'idle',
    scanSpeed: 0,  // 新增扫描速度状态
  }),
  
  getters: {
    scanProgress: (state) => {
      if (state.totalPaths <= 0) return 0
      const progress = (state.scannedPaths / state.totalPaths) * 100
      return Math.min(progress, 100)  // 移除 Math.round，允许小数点显示
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
      this.scanSpeed = 0  // 重置扫描速度
    },
    
    setIsScanning(value) {
      this.isScanning = value
      if (!value) {
    // 修改后：只在实际扫描状态下才设置为cancelled
    if (this.scanStatus === 'scanning') {
      this.scanStatus = 'cancelled'
    }
    this.scanSpeed = 0
  } else {
    this.scanStatus = 'scanning'
    this.showProgress = true
  }
    },

    setScanStatus(status) {
      this.scanStatus = status
      if (status === 'completed' || status === 'cancelled' || status === 'error') {
        store.setIsScanning(false)
        store.setShowProgress(false)
        this.isScanning = false
        this.scanSpeed = 0  // 扫描结束时重置速度
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

    setScanSpeed(value) {
      if (typeof value === 'number') {
        this.scanSpeed = value
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
      this.scanSpeed = 0  // 完成时重置速度
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