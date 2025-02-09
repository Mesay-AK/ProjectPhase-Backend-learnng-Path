import { createRouter, createWebHistory } from 'vue-router';
import Home from '../pages/Home.vue';
import Dashboard from '../pages/Dashboard.vue';
import Login from '../components/Login.vue';

const routes = [
  { path: '/', component: Home },
  { path: '/dashboard', component: Dashboard },
//   { path: '/login', component: Login },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
