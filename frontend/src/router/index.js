import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/Home.vue'
import About from '../pages/About.vue'
import ApiTest from '../pages/ApiTest.vue'
import Territory from '../pages/Territory.vue'
import Island from '../pages/TerritoryIsland.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/about', component: About },
  { path: '/api', component: ApiTest },
  // Route for the territory map
  { path: '/territory/:id', component: Territory, props: true },
  // THE CORRECTED NESTED ROUTE for the island view
  { path: '/territory/:id/:islandId', component: Island, props: true },
]

export default createRouter({
  history: createWebHistory(),
  routes
})