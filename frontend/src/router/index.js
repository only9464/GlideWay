import { createRouter, createWebHistory } from 'vue-router'
import ScannerPortsView from '../views/ScannerPortsView.vue'
import ScannerUrlView from '../views/ScannerUrlView.vue'

const routes = [
  { path: '/', component: ScannerPortsView },
  { path: '/urlScan', component: ScannerUrlView }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router