import { describe, expect, test } from 'vitest'
import { getRecordTabs, getWorkbenchMode } from '../src/utils/tab-config.js'

describe('tab config', () => {
  test('工人只看到扫码记录', () => {
    expect(getRecordTabs('worker')).toEqual([{ key: 'scan', label: '扫码记录' }])
  })

  test('管理员看到三个子 tab', () => {
    expect(getRecordTabs('admin')).toEqual([
      { key: 'scan', label: '扫码记录' },
      { key: 'status', label: '流程状态' },
      { key: 'order', label: '订单记录' },
    ])
  })

  test('工作台模式随角色变化', () => {
    expect(getWorkbenchMode('worker')).toBe('scan')
    expect(getWorkbenchMode('admin')).toBe('admin')
  })
})
