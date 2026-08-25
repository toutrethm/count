import { describe, expect, test } from 'vitest'
import { formatCapabilityLabel, formatCapabilityListLabel } from '../src/utils/capability-labels.js'

describe('capability labels', () => {
  test('formats a single capability code', () => {
    expect(formatCapabilityLabel('turn_sleeve_auto')).toBe('车套（自动）')
  })

  test('formats multiple capability codes', () => {
    expect(formatCapabilityListLabel('turn_outer,turn_sleeve_manual')).toBe('车外圆、车套（非自动）')
  })
})
