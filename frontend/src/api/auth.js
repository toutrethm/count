import { request } from './request'

export function login(payload) {
  return request({
    url: '/api/auth/login',
    method: 'POST',
    data: payload,
  })
}

export function registerWorker(payload) {
  return request({
    url: '/api/auth/register',
    method: 'POST',
    data: payload,
  })
}

export function getMe() {
  return request({
    url: '/api/auth/me',
    method: 'GET',
  })
}
