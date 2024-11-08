import { createRouter, createWebHistory } from 'vue-router'
import ScannerPortsView from '../views/ScannerPortsView.vue'
import DirsearchView from '../views/DirsearchView.vue'

const routes = [
  { path: '/', component: ScannerPortsView },
  { path: '/dirsearch', component: DirsearchView }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router