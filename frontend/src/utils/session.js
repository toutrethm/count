const TOKEN_KEY = 'piecework_token'
const USER_KEY = 'piecework_user'

export function getToken() {
  return uni.getStorageSync(TOKEN_KEY) || ''
}

export function setSession(token, user) {
  uni.setStorageSync(TOKEN_KEY, token)
  uni.setStorageSync(USER_KEY, user)
}

export function getUser() {
  return uni.getStorageSync(USER_KEY) || null
}

export function clearSession() {
  uni.removeStorageSync(TOKEN_KEY)
  uni.removeStorageSync(USER_KEY)
}
