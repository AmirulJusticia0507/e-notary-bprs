import { reactive } from 'vue'

// useFormValidation: state error sederhana per-field tanpa dependency tambahan.
export function useFormValidation() {
  const errors = reactive({})

  function setErrors(next) {
    Object.keys(errors).forEach((k) => delete errors[k])
    Object.assign(errors, next)
  }

  function clear() {
    Object.keys(errors).forEach((k) => delete errors[k])
  }

  function required(value, field, label) {
    if (value === undefined || value === null || String(value).trim() === '') {
      errors[field] = `${label} wajib diisi`
      return false
    }
    return true
  }

  function isValidEmail(value, field = 'email') {
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value ?? '')) {
      errors[field] = 'Format email tidak valid'
      return false
    }
    return true
  }

  return { errors, setErrors, clear, required, isValidEmail }
}
