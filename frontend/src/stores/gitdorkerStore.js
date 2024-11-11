import { defineStore } from 'pinia'

export const useGitdorkerStore = defineStore('gitdorker', {
  state: () => ({
    mainKeyword: '',
    subKeyword: '',
    token: '',
    searchResults: null,
    isSearching: false,
    searchStatus: 'idle', // idle, searching, completed, error
  }),

  actions: {
    resetSearch() {
      this.searchResults = null
      this.isSearching = false
      this.searchStatus = 'idle'
    },

    setIsSearching(value) {
      this.isSearching = value
      if (!value) {
        if (this.searchStatus === 'searching') {
          this.searchStatus = 'cancelled'
        }
      } else {
        this.searchStatus = 'searching'
      }
    },

    setSearchStatus(status) {
      this.searchStatus = status
      if (status === 'completed' || status === 'cancelled' || status === 'error') {
        this.isSearching = false
      }
    },

    setSearchResults(results) {
      this.searchResults = results
    },

    setKeywords(main, sub) {
      this.mainKeyword = main
      this.subKeyword = sub
    },

    setToken(token) {
      this.token = token
    },

    exportResults() {
      return {
        timestamp: new Date().toISOString(),
        mainKeyword: this.mainKeyword,
        subKeyword: this.subKeyword,
        results: this.searchResults
      }
    }
  }
}) 