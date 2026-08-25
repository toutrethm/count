import { afterEach, describe, expect, test } from 'vitest'
import { resolveApiBaseURL } from '../src/utils/api-base-url.js'

afterEach(() => {
  delete global.window
})

describe('api base url', () => {
  test('uses explicit base url when provided', () => {
    global.window = { location: { protocol: 'http:', hostname: '10.26.68.46' } }

    expect(resolveApiBaseURL('http://api.example.com:8080')).toBe('http://api.example.com:8080')
  })

  test('derives backend host from the current page hostname', () => {
    global.window = { location: { protocol: 'http:', hostname: '10.26.68.46' } }

    expect(resolveApiBaseURL()).toBe('http://10.26.68.46:8080')
  })

  test('falls back to localhost when window is unavailable', () => {
    expect(resolveApiBaseURL()).toBe('http://localhost:8080')
  })
})
