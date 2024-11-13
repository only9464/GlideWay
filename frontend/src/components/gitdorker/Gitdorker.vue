<template>
  <div class="scanner-component">
    <!-- 参数配置区域 -->
    <div class="input-group">
      <div class="input-item acrylic-input-box">
        <span class="input-label">主关键词</span>
        <el-input
          v-model="mainKeyword"
          placeholder="输入主关键词"
          clearable
        />
      </div>

      <div class="input-item acrylic-input-box">
        <span class="input-label">次关键词</span>
        <div class="subkeyword-wrapper">
          <el-input
            type="textarea"
            v-model="subKeyword"
            placeholder="输入多个次关键词（支持空格、逗号、分号、换行符分隔）"
            :rows="3"
          />
          <el-upload
            class="upload-btn"
            action=""
            :auto-upload="false"
            :show-file-list="false"
            accept=".txt"
            @change="handleFileUpload"
          >
            <el-button type="primary" size="small">
              从文件导入
            </el-button>
          </el-upload>
        </div>
      </div>

      <div class="input-item acrylic-input-box">
        <span class="input-label">GitHub Token</span>
        <el-input
          v-model="token"
          placeholder="输入 GitHub Token"
          clearable
          show-password
        />
      </div>
    </div>

    <!-- 预览关键词区域 -->
    <div v-if="splitKeywords.length > 0" class="keywords-preview acrylic-mini">
      <div class="preview-header">
        <span class="status-text">已识别的关键词</span>
        <span class="keyword-count">({{ splitKeywords.length }})</span>
      </div>
      <div class="keywords-list">
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

    <!-- 搜索按钮区域 -->
    <div class="action-container">
      <el-button
        type="primary"
        :loading="isSearching"
        @click="searchGithub"
        class="search-button"
      >
        {{ isSearching ? '搜索中...' : '开始搜索' }}
      </el-button>
    </div>

    <!-- 搜索结果表格 -->
    <el-table
      v-if="searchResults.length > 0"
      :data="searchResults"
      style="width: 100%"
      size="small"
      class="acrylic-effect"
    >
      <el-table-column
        type="index"
        label="序号"
        width="60"
        align="center"
        header-align="center"
      />

      <el-table-column
        prop="Total"
        label="总数"
        width="100"
        align="center"
        header-align="center"
      />

      <el-table-column
        label="搜索链接"
        min-width="400"
        align="left"
        header-align="center"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          <el-link
            type="primary"
            @click="openUrl(row.Link)"
            style="cursor: pointer"
          >
            {{ row.Link }}
          </el-link>
        </template>
      </el-table-column>

      <el-table-column
        label="Github仓库"
        width="120"
        align="center"
        header-align="center"
      >
        <template #default="{ row }">
          <el-button
            type="primary"
            size="small"
            @click="showItemsDialog(row.Items)"
          >
            查看详情 ({{ row.Items.length }})
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 详情抽屉 -->
    <el-drawer
      v-model="dialogVisible"
      title="搜索结果详情"
      size="75%"
      :destroy-on-close="true"
      direction="rtl"
    >
      <el-table
        :data="dialogData"
        style="width: 100%"
        size="small"
      >
        <el-table-column
          type="index"
          label="序号"
          width="60"
          align="center"
          header-align="center"
        />
        <el-table-column
          label="URL"
          min-width="180"
          align="left"
          header-align="center"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            <el-link
              type="primary"
              @click="openUrl(row)"
              style="cursor: pointer"
            >
              {{ row }}
            </el-link>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>
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
// 添加 openUrl 函数
const openUrl = (url) => {
  window.runtime.BrowserOpenURL(url)
}
// 添加文件上传处理函数
const handleFileUpload = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    subKeyword.value = e.target.result
  }
  reader.readAsText(file.raw)
}
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
.scanner-component {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 0 20px;
  gap: 20px;
  overflow: hidden;
}

.input-group {
  display: flex;
  gap: 16px;
  align-items: stretch;
}

.input-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
}

.input-label {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}

.subkeyword-wrapper {
  position: relative;
  width: 100%;
}

.upload-btn {
  position: absolute;
  right: 0;
  top: -30px;
}

.keywords-preview {
  padding: 12px;
}

.preview-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.keyword-count {
  color: #909399;
  font-size: 12px;
}

.keywords-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.keyword-tag {
  margin: 0;
}

.action-container {
  display: flex;
  justify-content: flex-end;
}

.search-button {
  min-width: 120px;
  transition: all 0.3s ease;
}

.search-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 毛玻璃效果类 */
.acrylic-input-box {
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.acrylic-mini {
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
  border-radius: 6px;
  transition: all 0.3s ease;
}

.acrylic-effect {
  background: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  transition: all 0.3s ease;
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
  .acrylic-input-box,
  .acrylic-mini,
  .acrylic-effect {
    background: rgba(255, 255, 255, 0.15);
  }

  .input-label,
  .status-text {
    color: rgba(255, 255, 255, 0.9);
  }

  :deep(.el-table) {
    background-color: transparent;
    --el-table-border-color: rgba(255, 255, 255, 0.1);
    --el-table-header-bg-color: rgba(255, 255, 255, 0.05);
    --el-table-row-hover-bg-color: rgba(255, 255, 255, 0.08);
    --el-table-text-color: #e6e6e6;
    --el-table-header-text-color: #ffffff;
  }
}
</style>