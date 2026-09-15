import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '../layouts/AppLayout.vue'
import AuthLayout from '../layouts/AuthLayout.vue'
import LoginView from '../views/auth/LoginView.vue'
import DashboardOverview from '../views/dashboard/DashboardOverview.vue'
import OrderCreateWizard from '../views/orders/OrderCreateWizard.vue'
import OrderDetailView from '../views/orders/OrderDetailView.vue'
import OrderListView from '../views/orders/OrderListView.vue'

const routes = [
  {
    path: '/login',
    component: AuthLayout,
    meta: { guest: true },
    children: [{ path: '', name: 'login', component: LoginView }],
  },
  {
    path: '/',
    component: AppLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: { name: 'dashboard' } },
      { path: 'dashboard', name: 'dashboard', component: DashboardOverview },
      { path: 'orders', name: 'orders', component: OrderListView },
      { path: 'orders/create', name: 'order-create', component: OrderCreateWizard },
      { path: 'orders/:id', name: 'order-detail', component: OrderDetailView, props: true },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const token = localStorage.getItem('enotary_token')
  if (to.meta.requiresAuth && !token) return { name: 'login' }
  if (to.meta.guest && token) return { name: 'dashboard' }
  return true
})

export default router
