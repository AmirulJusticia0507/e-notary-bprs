import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '../layouts/AppLayout.vue'
import AuthLayout from '../layouts/AuthLayout.vue'
import LoginView from '../views/auth/LoginView.vue'
import RegisterView from '../views/auth/RegisterView.vue'
import ForbiddenView from '../views/auth/ForbiddenView.vue'
import DashboardOverview from '../views/dashboard/DashboardOverview.vue'
import OrderCreateWizard from '../views/orders/OrderCreateWizard.vue'
import OrderDetailView from '../views/orders/OrderDetailView.vue'
import OrderListView from '../views/orders/OrderListView.vue'
import NasabahDashboard from '../views/nasabah/NasabahDashboard.vue'

// Pembagian tampilan per role:
// - dashboard, orders: admin, legal_officer, notary (staf)
// - order-create: admin, legal_officer (notaris tidak boleh buat order,
//   hanya update processing dokumen dari order yang ditugaskan ke mereka)
// - saya (dashboard + pengajuan nasabah): nasabah saja
const STAFF = ['admin', 'legal_officer', 'notary']
const LEGAL_TEAM = ['admin', 'legal_officer']
const NASABAH = ['nasabah']

export function homeFor(role) {
  return role === 'nasabah' ? 'nasabah-dashboard' : 'dashboard'
}

const routes = [
  {
    path: '/login',
    component: AuthLayout,
    meta: { guest: true },
    children: [{ path: '', name: 'login', component: LoginView }],
  },
  {
    path: '/register',
    component: AuthLayout,
    meta: { guest: true },
    children: [{ path: '', name: 'register', component: RegisterView }],
  },
  {
    path: '/',
    component: AppLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: (to) => ({ name: homeFor(getStoredRole()) }) },
      { path: 'dashboard', name: 'dashboard', component: DashboardOverview, meta: { roles: STAFF } },
      { path: 'orders', name: 'orders', component: OrderListView, meta: { roles: STAFF } },
      { path: 'orders/create', name: 'order-create', component: OrderCreateWizard, meta: { roles: LEGAL_TEAM } },
      { path: 'orders/:id', name: 'order-detail', component: OrderDetailView, props: true, meta: { roles: STAFF } },
      { path: 'saya', name: 'nasabah-dashboard', component: NasabahDashboard, meta: { roles: NASABAH } },
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
  if (to.meta.guest && token) return { name: homeFor(getStoredRole()) }
  const allowed = to.meta.roles
  if (token && Array.isArray(allowed) && allowed.length > 0) {
    const role = getStoredRole()
    if (!allowed.includes(role)) return { name: 'forbidden' }
  }
  return true
})

export default router
