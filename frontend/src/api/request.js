import { clearSession, getToken } from '../utils/session'
import { redirectTo } from '../utils/navigation'
import { resolveApiBaseURL } from '../utils/api-base-url'

const baseURL = resolveApiBaseURL(import.meta.env.VITE_API_BASE_URL)

export function request({ url, method = 'GET', data, header = {} }) {
  const token = getToken()
  const mergedHeader = { ...header }
  if (token) {
    mergedHeader.Authorization = `Bearer ${token}`
  }

  return new Promise((resolve, reject) => {
    uni.request({
      url: `${baseURL}${url}`,
      method,
      data,
      header: mergedHeader,
      success: (response) => {
        if (response.statusCode === 401) {
          clearSession()
          redirectTo('/pages/login/login')
          reject(response)
          return
        }
        if (response.statusCode >= 200 && response.statusCode < 300) {
          resolve(response.data)
          return
        }
        reject(response)
      },
      fail: reject,
    })
  })
}
