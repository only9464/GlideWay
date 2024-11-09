import { createRouter, createWebHistory } from 'vue-router'
import ScannerPortsView from '../views/ScannerPortsView.vue'
import DirsearchView from '../views/DirsearchView.vue'
import JsfinderView from '../views/JsfinderView.vue'
import GitdorkerView from '../views/GitdorkerView.vue'
const routes = [
  { path: '/', component: ScannerPortsView },
  { path: '/dirsearch', component: DirsearchView },
  { path: '/jsfinder', component: JsfinderView },
  { path: '/gitdorker', component: GitdorkerView }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router