import { createRouter, createWebHistory } from 'vue-router';
import LoginView from '../views/LoginView.vue';
import DashboardView from '../views/DashboardView.vue';
import CheckinView from '../views/CheckinView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/login' },
    { path: '/login', component: LoginView },
    { path: '/dashboard', component: DashboardView },
    { path: '/checkin', component: CheckinView }
  ]
});

router.beforeEach((to) => {
  const token = localStorage.getItem('token');
  if (to.path === '/dashboard' && !token) {
    return '/login';
  }
  if (to.path === '/checkin' && !token) {
    return '/login';
  }
  if (to.path === '/login' && token) {
    return '/dashboard';
  }
});

export default router;
