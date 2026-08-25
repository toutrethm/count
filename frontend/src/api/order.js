import { request } from './request'

export function listOrders() {
  return request({
    url: '/api/admin/orders',
    method: 'GET',
  })
}

export function getOrder(id) {
  return request({
    url: `/api/orders/${id}`,
    method: 'GET',
  })
}

export function getOrderByNo(orderNo) {
  return request({
    url: `/api/orders/by-no/${encodeURIComponent(orderNo)}`,
    method: 'GET',
  })
}

export function createOrder(payload) {
  return request({
    url: '/api/admin/orders',
    method: 'POST',
    data: payload,
  })
}
