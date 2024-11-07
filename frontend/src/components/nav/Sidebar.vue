<template>
  <div class="sidebar-container">
    <div class="sidebar acrylic-effect">
      <div class="menu-container">
        <ul>
          <li>
            <router-link to="/" class="nav-link">
              <span>端口扫描器</span>
            </router-link>
          </li>
          <li>
            <router-link to="/urlScan" class="nav-link">
              <span>url扫描器</span>
            </router-link>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
</script>
<style scoped>
.sidebar-container {
  position: fixed;
  top: 20px;  /* 增加顶部间距 */
  left: 20px;  /* 增加左侧间距 */
  bottom: 20px;  /* 增加底部间距 */
  width: 200px;
  z-index: 100;
}

.sidebar {
  height: 100%;
  border-radius: 16px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.menu-container {
  padding-top: 20px;
  flex-grow: 1;
}

.acrylic-effect {
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.15),
    rgba(255, 255, 255, 0.08)
  );
  backdrop-filter: blur(20px) saturate(125%);
  -webkit-backdrop-filter: blur(20px) saturate(125%);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.15);
}

.sidebar ul {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.sidebar li {
  margin: 8px 0;
}

.nav-link {
  color: #2c3e50; /* 更深的蓝灰色 */
  text-decoration: none;
  padding: 12px 24px;
  display: block;
  position: relative;
  transition: all 0.3s ease;
  margin: 0 8px;
  border-radius: 8px;
  font-weight: 600; /* 加粗 */
  letter-spacing: 0.5px;
  text-shadow: 0 1px 1px rgba(255, 255, 255, 0.5); /* 白色文字阴影 */
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.3),
    rgba(255, 255, 255, 0.2)
  );
}

.nav-link:hover {
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.4),
    rgba(255, 255, 255, 0.3)
  );
  color: #1a365d; /* 悬停时更深的颜色 */
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.7);
}

.nav-link:active {
  background: rgba(255, 255, 255, 0.35);
}

.router-link-active {
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.5),
    rgba(255, 255, 255, 0.4)
  );
  color: #1a365d;
  font-weight: 700;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.7);
}

.nav-link::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.nav-link:hover::before,
.router-link-active::before {
  opacity: 1;
}

.nav-link span {
  position: relative;
  z-index: 1;
}

@media (prefers-color-scheme: dark) {
  .acrylic-effect {
    background: linear-gradient(
      135deg,
      rgba(255, 255, 255, 0.1),
      rgba(255, 255, 255, 0.05)
    );
  }
  
  .nav-link {
    color: rgba(255, 255, 255, 0.9); /* 深色模式下使用白色 */
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
    background: linear-gradient(
      90deg,
      rgba(255, 255, 255, 0.15),
      rgba(255, 255, 255, 0.1)
    );
  }

  .nav-link:hover,
  .router-link-active {
    color: #ffffff;
  }
}
</style>