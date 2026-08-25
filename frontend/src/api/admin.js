import { request } from './request'

export function listWorkers() {
  return request({
    url: '/api/admin/workers',
    method: 'GET',
  })
}

export function setWorkerProcesses(workerId, processIds) {
  return request({
    url: `/api/admin/workers/${workerId}/processes`,
    method: 'PUT',
    data: {
      process_ids: processIds,
    },
  })
}
