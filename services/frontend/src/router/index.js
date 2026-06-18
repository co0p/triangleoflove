import { createRouter, createWebHistory } from 'vue-router';
import LoginView from '../views/LoginView.vue';
import DashboardView from '../views/DashboardView.vue';
import CheckinView from '../views/CheckinView.vue';
import PairingView from '../views/PairingView.vue';
import ProfileView from '../views/ProfileView.vue';
import InsightsView from '../views/InsightsView.vue';
import InsightsWeeklyView from '../views/InsightsWeeklyView.vue';
import AdminView from '../views/AdminView.vue';
import RegisterView from '../views/RegisterView.vue';

function getTokenRole(token) {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.role || null;
  } catch {
    return null;
  }
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/login' },
    { path: '/login', component: LoginView },
    { path: '/register', component: RegisterView },
    { path: '/dashboard', component: DashboardView },
    { path: '/session', component: CheckinView },
    { path: '/pairing', component: PairingView },
    { path: '/profile', component: ProfileView },
    { path: '/insights', component: InsightsWeeklyView },
    { path: '/insights/:date', component: InsightsView },
    { path: '/admin', component: AdminView }
  ]
});

router.beforeEach((to) => {
  const token = localStorage.getItem('token');
  if (to.path === '/dashboard' && !token) {
    return '/login';
  }
  if (to.path === '/session' && !token) {
    return '/login';
  }
  if (to.path === '/pairing' && !token) {
    return '/login';
  }
  if (to.path === '/profile' && !token) {
    return '/login';
  }
  if (to.path === '/insights' && !token) {
    return '/login';
  }
  if (to.path.startsWith('/insights/') && !token) {
    return '/login';
  }
  if (to.path === '/admin') {
    if (!token) return '/login';
    if (getTokenRole(token) !== 'admin') return '/dashboard';
  }
  if (to.path === '/register' && token) {
    return '/dashboard';
  }
  if (to.path === '/login' && token) {
    return '/dashboard';
  }
});

export default router;
