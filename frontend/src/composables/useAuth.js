import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/userStore'

// useAuth: akses state otentikasi + guard logout.
export function useAuth() {
  const store = useUserStore()
  const router = useRouter()

  async function logout() {
    store.logout()
    await router.push({ name: 'login' })
  }

  return { store, logout }
}
