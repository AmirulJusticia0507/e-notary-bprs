import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '../layouts/AppLayout.vue'
import AuthLayout from '../layouts/AuthLayout.vue'
import LoginView from '../views/auth/LoginView.vue'
import ForbiddenView from '../views/auth/ForbiddenView.vue'
import DashboardOverview from '../views/dashboard/DashboardOverview.vue'
import OrderCreateWizard from '../views/orders/OrderCreateWizard.vue'
import OrderDetailView from '../views/orders/OrderDetailView.vue'
import OrderListView from '../views/orders/OrderListView.vue'

// Pembagian tampilan per role:
// - dashboard, orders, order-detail: admin, legal_officer, notary
// - order-create: admin, legal_officer (notaris tidak boleh buat order,
//   hanya update processing dokumen dari order yang ditugaskan ke mereka)
const ALL_STAFF = ['admin', 'legal_officer', 'notary']
const LEGAL_TEAM = ['admin', 'legal_officer']

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
      { path: 'dashboard', name: 'dashboard', component: DashboardOverview, meta: { roles: ALL_STAFF } },
      { path: 'orders', name: 'orders', component: OrderListView, meta: { roles: ALL_STAFF } },
      { path: 'orders/create', name: 'order-create', component: OrderCreateWizard, meta: { roles: LEGAL_TEAM } },
      { path: 'orders/:id', name: 'order-detail', component: OrderDetailView, props: true, meta: { roles: ALL_STAFF } },
      { path: '403', name: 'forbidden', component: ForbiddenView },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

function getStoredRole() {
  try {
    return JSON.parse(localStorage.getItem('enotary_user') ?? 'null')?.role ?? ''
  } catch {
    return ''
  }
}

router.beforeEach((to) => {
  const token = localStorage.getItem('enotary_token')
  if (to.meta.requiresAuth && !token) return { name: 'login' }
  if (to.meta.guest && token) return { name: 'dashboard' }
  const allowed = to.meta.roles
  if (token && Array.isArray(allowed) && allowed.length > 0) {
    const role = getStoredRole()
    if (!allowed.includes(role)) return { name: 'forbidden' }
  }
  return true
})

export default router
