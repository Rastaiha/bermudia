import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/Home.vue'
import About from '../pages/About.vue'
import ApiTest from '../pages/ApiTest.vue'
import Territory from '../pages/Territory.vue'
import Island from '../pages/TerritoryIsland.vue'
import Login from '../pages/Login.vue'
import UserPage from '../pages/UserPage.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/about', component: About },
  { path: '/api', component: ApiTest },
  { path: '/territory/:id', component: Territory, props: true },
  { path: '/territory/:id/:islandId', component: Island, props: true },
  { path: '/login', component: Login},
  { path: '/user_page', component: UserPage},
]

export default createRouter({
  history: createWebHistory(),
  routes
})