import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { authService } from '../services/authService'

const TOKEN_KEY = 'enotary_token'
const USER_KEY = 'enotary_user'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem(TOKEN_KEY) ?? '')
  const user = ref(JSON.parse(localStorage.getItem(USER_KEY) ?? 'null'))
  const loading = ref(false)
  const error = ref('')

  const isAuthenticated = computed(() => Boolean(token.value))
  const role = computed(() => user.value?.role ?? '')

  function persist(nextToken, nextUser) {
    token.value = nextToken
    user.value = nextUser
    localStorage.setItem(TOKEN_KEY, nextToken)
    localStorage.setItem(USER_KEY, JSON.stringify(nextUser))
  }

  async function login(email, password) {
    loading.value = true
    error.value = ''
    try {
      const { token: t, user: u } = await authService.login(email, password)
      persist(t, u)
      return u
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Email atau password salah'
      throw e
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  }

  return { token, user, loading, error, isAuthenticated, role, login, logout }
})
