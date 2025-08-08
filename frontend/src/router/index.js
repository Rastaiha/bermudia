import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/Home.vue'
import About from '../pages/About.vue'
import ApiTest from '../pages/ApiTest.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/about', component: About },
  { path: '/api', component: ApiTest },
]

export default createRouter({
  history: createWebHistory(),
  routes
})
