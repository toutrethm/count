import { request } from './request'

export function previewScanOrder(qrToken) {
  return request({
    url: `/api/scans/preview/${encodeURIComponent(qrToken)}`,
    method: 'GET',
  })
}

export function recordScan(payload) {
  return request({
    url: '/api/scans/record',
    method: 'POST',
    data: payload,
  })
}

export function listMyScanRecords() {
  return request({
    url: '/api/scan-records/mine',
    method: 'GET',
  })
}

export function listAllScanRecords() {
  return request({
    url: '/api/admin/scan-records',
    method: 'GET',
  })
}
