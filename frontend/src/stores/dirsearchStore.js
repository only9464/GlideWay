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
    // 计算扫描进度百分比
    scanProgress: (state) => {
      if (state.totalPaths <= 0) return 0
      const progress = Math.round((state.scannedPaths / state.totalPaths) * 100)
      return Math.min(progress, 99)
    },
  },
  
  actions: {
    // 重置扫描状态 - 仅在开始新扫描时调用
    resetScan() {
      this.foundPaths = []
      this.scannedPaths = 0
      this.totalPaths = 0
      this.isScanning = false
      this.scanStatus = 'idle'
      // 注意：不重置 showProgress，因为我们想保持进度条显示
    },
    
    // 设置扫描状态
    setIsScanning(value) {
      this.isScanning = value
      if (!value) {
        if (this.scanStatus === 'scanning') {
          this.scanStatus = 'cancelled'
        }
        // 移除这里的 setTimeout，不再自动隐藏进度条
      } else {
        this.scanStatus = 'scanning'
        this.showProgress = true
      }
    },

    // 设置扫描状态
    setScanStatus(status) {
      this.scanStatus = status
      if (status === 'completed' || status === 'cancelled' || status === 'error') {
        this.isScanning = false
        // 不再隐藏进度条和清除结果
      }
    },

    // 设置是否显示进度 - 只在开始新扫描时设置为 true
    setShowProgress(value) {
      if (value) {
        this.showProgress = true
      }
      // 忽略设置为 false 的情况，保持进度条显示
    },
    
// 设置已扫描路径数
setScannedPaths(value) {
    if (typeof value === 'number') {
        this.scannedPaths = value
    }
},

// 设置总路径数
setTotalPaths(value) {
    if (typeof value === 'number' && value > 0) {
        this.totalPaths = value
    }
},

    // 清除所有发现的路径 - 仅在用户明确要求时调用
    clearFoundPaths() {
      this.foundPaths = []
    },

    // 添加发现的路径
    addPath(pathInfo) {
      this.foundPaths.push({
        path: pathInfo.path,
        fullUrl: pathInfo.fullUrl,
        statusCode: pathInfo.statusCode,
        contentType: pathInfo.contentType,
        contentLength: pathInfo.contentLength,
      })
    },

    // 设置排序配置
    setSortConfig({ prop, order }) {
      if (!prop || !order) {
        this.sortConfig.prop = null
        this.sortConfig.order = null
        return
      }
      this.sortConfig.prop = prop
      this.sortConfig.order = order
    },

    // 设置完成状态
    setComplete() {
      this.isScanning = false
      this.scanStatus = 'completed'
      // 不再添加隐藏进度条的逻辑
    },

    // 导出扫描结果
    exportResults() {
      return {
        timestamp: new Date().toISOString(),
        totalScanned: this.scannedPaths,
        foundPaths: this.foundPaths
      }
    }
  }
})