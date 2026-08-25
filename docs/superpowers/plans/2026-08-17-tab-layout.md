# 工厂计件前端 Tab 结构实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把现有 uni-app 前端改成底部 3 个 tab，并让记录页按角色显示不同子 tab 与内容。

**Architecture:** 以 `pages.json` 的 `tabBar` 作为主导航，`index` 作为工作台，新增 `records` 和 `profile` 页面。登录成功后统一 `switchTab` 到工作台；记录页内部用本地状态切换“扫码记录”和“流程状态”，管理员可见后者，普通用户只看自己的扫码记录。

**Tech Stack:** uni-app, Vue 3, Vite, Go 后端现有接口

---

### Task 1: 角色与 tab 配置

**Files:**
- Create: `frontend/src/utils/tab-config.js`
- Create: `frontend/tests/tab-config.test.js`

- [ ] **Step 1: Write the failing test**

```js
import { describe, expect, test } from 'vitest'
import { getRecordTabs, getWorkbenchMode } from '../src/utils/tab-config.js'

describe('tab config', () => {
  test('worker only sees scan record tab', () => {
    expect(getRecordTabs('worker')).toEqual([{ key: 'scan', label: '扫码记录' }])
  })

  test('admin sees scan record and order status tabs', () => {
    expect(getRecordTabs('admin')).toEqual([
      { key: 'scan', label: '扫码记录' },
      { key: 'status', label: '流程状态' },
    ])
  })

  test('workbench mode follows role', () => {
    expect(getWorkbenchMode('worker')).toBe('scan')
    expect(getWorkbenchMode('admin')).toBe('admin')
  })
})
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm exec vitest run tests/tab-config.test.js`
Expected: fail because `tab-config.js` does not exist yet

- [ ] **Step 3: Write minimal implementation**

```js
export function getWorkbenchMode(role) {
  return role === 'admin' ? 'admin' : 'scan'
}

export function getRecordTabs(role) {
  if (role === 'admin') {
    return [
      { key: 'scan', label: '扫码记录' },
      { key: 'status', label: '流程状态' },
    ]
  }
  return [{ key: 'scan', label: '扫码记录' }]
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `npm exec vitest run tests/tab-config.test.js`
Expected: PASS

### Task 2: 页面与导航结构

**Files:**
- Modify: `frontend/src/pages.json`
- Modify: `frontend/src/pages/login/login.vue`
- Modify: `frontend/src/utils/navigation.js`
- Modify: `frontend/src/pages/index/index.vue`
- Create: `frontend/src/pages/records/records.vue`
- Create: `frontend/src/pages/profile/profile.vue`

- [ ] Add tabBar pages and switchTab navigation
- [ ] Split workbench into admin/worker states
- [ ] Make records page show role-specific sub tabs
- [ ] Add profile page with logout and basic identity info

### Task 3: Verification

**Files:**
- None

- [ ] Run frontend build

Run: `npm run build:h5`
Expected: build succeeds with the new tab layout
