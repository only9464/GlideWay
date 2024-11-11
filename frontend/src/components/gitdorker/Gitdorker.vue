<template>
  <div class="Gitdorker-container">
    <h1>Gitdorker</h1>
    
    <!-- GitHub 搜索部分 -->
    <div class="search-section">
      <input type="text" v-model="mainKeyword" placeholder="主关键词" />
      <el-input
        type="textarea"
        v-model="subKeyword"
        placeholder="输入多个次关键词（支持空格、逗号、分号、换行符分隔）"
        :rows="3"
      />
      <input type="text" v-model="token" placeholder="GitHub Token" />
      <button @click="searchGithub" :disabled="isSearching">
        {{ isSearching ? '搜索中...' : '搜索 GitHub' }}
      </button>
      
      <!-- 预览分割后的关键词 -->
      <div v-if="splitKeywords.length > 0" class="preview-keywords">
        <p>已识别的关键词：</p>
        <el-tag
          v-for="(keyword, index) in splitKeywords"
          :key="index"
          class="keyword-tag"
          size="small"
        >
          {{ keyword }}
        </el-tag>
      </div>
    </div>

    <!-- 搜索结果表格 -->
    <el-table v-if="searchResults.length > 0" :data="searchResults" style="width: 100%; margin-top: 20px">
      <el-table-column 
        type="index" 
        label="序号" 
        width="80"
        align="center"
        header-align="center" />

      <el-table-column 
        prop="Total" 
        label="总数" 
        width="120"
        align="center"
        header-align="center" />

      <el-table-column 
        label="搜索链接" 
        min-width="500"
        align="center"
        header-align="center">
        <template #default="{ row }">
          <el-link type="primary" :href="row.Link" target="_blank">
            {{ row.Link }}
          </el-link>
        </template>
      </el-table-column>

      <el-table-column 
        label="Github仓库" 
        width="120"
        align="center"
        header-align="center">
        <template #default="{ row }">
          <el-button type="primary" @click="showItemsDialog(row.Items)">
            查看详情 ({{ row.Items.length }})
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 数据项详情对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="搜索结果详情"
      width="70%"
      :destroy-on-close="true"
    >
      <el-table :data="dialogData" style="width: 100%">
        <el-table-column label="URL" min-width="180">
          <template #default="{ row }">
            <el-link type="primary" :href="row" target="_blank">
              {{ row }}
            </el-link>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useGitdorkerStore } from '../../stores/gitdorkerStore'
import { ElMessage } from 'element-plus'

const store = useGitdorkerStore()

const mainKeyword = ref(store.mainKeyword)
const subKeyword = ref(store.subKeyword)
const token = ref(store.token)
const searchResults = ref([])
const isSearching = ref(false)

// 使用计算属性来处理关键词分割
const splitKeywords = computed(() => {
  if (!subKeyword.value) return [];
  
  // 使用多个分隔符组合的正则表达式
  const regex = /[,，;；\s\n\r]+/;
  
  // 分割、过滤空值、去重、去除首尾空格
  const keywords = subKeyword.value
    .split(regex)
    .map(k => k.trim())
    .filter(k => k !== '')
    .filter((value, index, self) => self.indexOf(value) === index);
    
  return keywords;
});

// 对话框相关
const dialogVisible = ref(false)
const dialogData = ref([])

const showItemsDialog = (items) => {
  dialogData.value = items
  dialogVisible.value = true
}

// 单个关键词搜索
async function searchSingleKeyword(subKey) {
  try {
    const result = await window.go.gitdorker.App.Gitdorker(
      mainKeyword.value,
      subKey,
      token.value
    )
    return result
  } catch (error) {
    console.error(`搜索关键词 "${subKey}" 失败:`, error)
    ElMessage.warning(`关键词 "${subKey}" 搜索失败: ${error.message}`)
    return null
  }
}

async function searchGithub() {
  try {
    if (!mainKeyword.value) {
      ElMessage.warning('请输入主关键词');
      return;
    }
    
    if (splitKeywords.value.length === 0) {
      ElMessage.warning('请输入至少一个次关键词');
      return;
    }

    if (!token.value) {
      ElMessage.warning('请输入 GitHub Token');
      return;
    }
    
    isSearching.value = true
    store.setIsSearching(true)
    store.setKeywords(mainKeyword.value, subKeyword.value)
    store.setToken(token.value)
    
    // 清空之前的搜索结果
    searchResults.value = []
    
    // 逐个搜索每个关键词
    for (const keyword of splitKeywords.value) {
      const result = await searchSingleKeyword(keyword)
      if (result) {
        searchResults.value.push(result)
      }
      // 可以添加适当的延时，避免触发 GitHub API 限制
      await new Promise(resolve => setTimeout(resolve, 1000))
    }
    
    store.setSearchStatus('completed')
    console.log("GitHub search results:", searchResults.value)
  } catch (error) {
    console.error("搜索 GitHub 失败:", error)
    store.setSearchStatus('error')
    ElMessage.error('搜索失败：' + error.message)
  } finally {
    isSearching.value = false
    store.setIsSearching(false)
  }
}
</script>

<style scoped>
.Gitdorker-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.search-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin: 20px 0;
}

.preview-keywords {
  margin-top: 10px;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.keyword-tag {
  margin: 4px;
}

/* 添加表格内链接样式 */
:deep(.el-table .cell) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

input {
  padding: 8px;
  width: 100%;
}

button {
  padding: 10px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover:not(:disabled) {
  background-color: #45a049;
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}
</style>