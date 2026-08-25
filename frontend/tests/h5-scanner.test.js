import { describe, expect, test } from 'vitest'
import { extractScanText, isH5CameraScanSupported } from '../src/utils/h5-scanner.js'

describe('h5 scanner', () => {
  test('detects camera support', () => {
    expect(
      isH5CameraScanSupported({
        navigator: { mediaDevices: { getUserMedia: () => Promise.resolve() } },
        window: {},
      })
    ).toBe(true)
    expect(isH5CameraScanSupported({ navigator: {}, window: {} })).toBe(false)
  })

  test('extracts qr text from common result shapes', () => {
    expect(extractScanText('12345')).toBe('12345')
    expect(extractScanText({ text: '  abc  ' })).toBe('abc')
    expect(extractScanText({ result: ' xyz ' })).toBe('xyz')
  })
})
