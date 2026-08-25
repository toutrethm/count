import { request } from './request'

export function listProcesses() {
  return request({
    url: '/api/processes',
    method: 'GET',
  })
}

export function listMyProcesses() {
  return request({
    url: '/api/me/processes',
    method: 'GET',
  })
}

export function createProcess(payload) {
  return request({
    url: '/api/admin/processes',
    method: 'POST',
    data: payload,
  })
}
