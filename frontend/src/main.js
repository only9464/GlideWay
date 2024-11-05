import { createApp } from 'vue'
import App from './App.vue'
import router from './router/index' // 新增
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import './style.css'

createApp(App)
  .use(router) // 新增
  .use(ElementPlus)
  .mount('#app')