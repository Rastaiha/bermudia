import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/Home.vue'
import About from '../pages/About.vue'
import ApiTest from '../pages/ApiTest.vue'
import Territory from '../pages/Territory.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/about', component: About },
  { path: '/api', component: ApiTest },
  { path: '/territory', component: Territory },
]

export default createRouter({
  history: createWebHistory(),
  routes
})