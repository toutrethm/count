import { afterEach, describe, expect, test, vi } from 'vitest'
import { navigateTo, switchTabTo } from '../src/utils/navigation.js'

afterEach(() => {
  vi.restoreAllMocks()
  delete global.window
  delete global.uni
})

describe('navigation', () => {
  test('navigateTo uses uni.navigateTo on H5', () => {
    global.window = { location: { hash: '#/pages/login/login' } }
    global.uni = {
      navigateTo: vi.fn(),
    }

    navigateTo('/pages/register/register')

    expect(global.uni.navigateTo).toHaveBeenCalledWith({ url: '/pages/register/register' })
    expect(global.window.location.hash).toBe('#/pages/login/login')
  })

  test('switchTabTo uses uni.switchTab on H5', () => {
    global.window = { location: { hash: '#/pages/login/login' } }
    global.uni = {
      switchTab: vi.fn(),
    }

    switchTabTo('/pages/index/index')

    expect(global.uni.switchTab).toHaveBeenCalledWith({ url: '/pages/index/index' })
    expect(global.window.location.hash).toBe('#/pages/login/login')
  })
})
