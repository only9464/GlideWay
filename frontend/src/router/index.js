import { createRouter, createWebHistory } from 'vue-router'
import ScannerPortsView from '../views/ScannerPortsView.vue'
import DirsearchView from '../views/DirsearchView.vue'
import JsfinderView from '../views/JsfinderView.vue'
const routes = [
  { path: '/', component: ScannerPortsView },
  { path: '/dirsearch', component: DirsearchView },
  { path: '/jsfinder', component: JsfinderView }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router