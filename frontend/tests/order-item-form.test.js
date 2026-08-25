import { describe, expect, test } from 'vitest'
import { buildOrderItemPayload, buildSizeSpec, getComponentDefinition } from '../src/utils/order-item-form.js'

describe('order item form helpers', () => {
  test('B柱使用大径和头长度的和', () => {
    expect(getComponentDefinition('b_pillar').fields).toEqual([
      { key: 'smallHeadDiameter', label: '小头直径' },
      { key: 'totalLength', label: '总长' },
      { key: 'largeDiameterAndHeadLength', label: '大径和头长度的和' },
    ])
    expect(
      buildSizeSpec('b_pillar', {
        smallHeadDiameter: '8',
        totalLength: '120',
        largeDiameterAndHeadLength: '30',
      })
    ).toBe('8*120*30')
  })

  test('中导套使用紧位长度和头部厚度的和', () => {
    expect(getComponentDefinition('middle_guide_sleeve').fields).toEqual([
      { key: 'innerHoleDiameter', label: '内孔直径' },
      { key: 'totalLength', label: '总长' },
      { key: 'tightFitAndHeadThickness', label: '紧位长度和头部厚度的和' },
    ])
    expect(
      buildSizeSpec('middle_guide_sleeve', {
        innerHoleDiameter: '26',
        totalLength: '90',
        tightFitAndHeadThickness: '24',
      })
    ).toBe('26*90*24')
  })

  test('提交 payload 包含元件类型和尺寸明细', () => {
    expect(
      buildOrderItemPayload(
        {
          typeKey: 'guide_bush',
          quantity: '20',
          dimensions: {
            innerDiameter: '30',
            length: '40',
          },
        },
        1
      )
    ).toEqual({
      item_no: '2',
      component_type: 'guide_bush',
      part_name: '导套',
      spec: '30*40',
      dimensions: {
        innerDiameter: '30',
        length: '40',
      },
      quantity: 20,
      unit: '件',
      remark: '',
    })
  })
})
