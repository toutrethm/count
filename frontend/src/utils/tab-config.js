export function getWorkbenchMode(role) {
  return role === 'admin' ? 'admin' : 'scan'
}

export function getRecordTabs(role) {
  if (role === 'admin') {
    return [
      { key: 'scan', label: '扫码记录' },
      { key: 'status', label: '流程状态' },
      { key: 'order', label: '订单记录' },
    ]
  }
  return [{ key: 'scan', label: '扫码记录' }]
}
