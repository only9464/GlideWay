import { createRouter, createWebHistory } from 'vue-router'
import ScannerView from '../views/ScannerView.vue'
// import About from './components/About.vue'

const routes = [
  { path: '/', component: ScannerView },
  // { path: '/about', component: About }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router