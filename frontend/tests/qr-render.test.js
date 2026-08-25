import { describe, expect, test } from 'vitest'
import { buildQrRectangles } from '../src/utils/qr-render.js'

describe('qr render helper', () => {
  test('builds rectangles for black modules only', () => {
    const rects = buildQrRectangles([
      [{ isBlack: false }, { isBlack: true }],
      [{ isBlack: true }, { isBlack: false }],
    ])

    expect(rects).toEqual([
      { x: 1, y: 0 },
      { x: 0, y: 1 },
    ])
  })
})
