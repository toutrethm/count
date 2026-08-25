const COMPONENT_TYPES = [
  {
    key: 'guide_pillar',
    label: '导柱',
    fields: [
      { key: 'innerDiameter', label: '内径' },
      { key: 'length', label: '长度' },
    ],
  },
  {
    key: 'top_pin',
    label: '顶针',
    fields: [
      { key: 'innerDiameter', label: '内径' },
      { key: 'length', label: '长度' },
    ],
  },
  {
    key: 'b_pillar',
    label: 'B柱',
    fields: [
      { key: 'smallHeadDiameter', label: '小头直径' },
      { key: 'totalLength', label: '总长' },
      { key: 'largeDiameterAndHeadLength', label: '大径和头长度的和' },
    ],
  },
  {
    key: 'pull_rod',
    label: '拉杆',
    fields: [
      { key: 'innerDiameter', label: '内径' },
      { key: 'length', label: '长度' },
    ],
  },
  {
    key: 'middle_guide_pillar',
    label: '中导柱',
    fields: [
      { key: 'innerDiameter', label: '内径' },
      { key: 'length', label: '长度' },
    ],
  },
  {
    key: 'guide_bush',
    label: '导套',
    fields: [
      { key: 'innerDiameter', label: '内径' },
      { key: 'length', label: '长度' },
    ],
  },
  {
    key: 'straight_sleeve',
    label: '直套',
    fields: [
      { key: 'innerDiameter', label: '内径' },
      { key: 'length', label: '长度' },
    ],
  },
  {
    key: 'a_sleeve',
    label: 'A套',
    fields: [
      { key: 'innerDiameter', label: '内径' },
      { key: 'length', label: '长度' },
    ],
  },
  {
    key: 'b_sleeve',
    label: 'B套',
    fields: [
      { key: 'innerDiameter', label: '内径' },
      { key: 'length', label: '长度' },
    ],
  },
  {
    key: 'middle_guide_sleeve',
    label: '中导套',
    fields: [
      { key: 'innerHoleDiameter', label: '内孔直径' },
      { key: 'totalLength', label: '总长' },
      { key: 'tightFitAndHeadThickness', label: '紧位长度和头部厚度的和' },
    ],
  },
]

export function getComponentTypes() {
  return COMPONENT_TYPES.map(({ key, label }) => ({ key, label }))
}

export function getComponentDefinition(typeKey) {
  return COMPONENT_TYPES.find((item) => item.key === typeKey) || COMPONENT_TYPES[0]
}

export function createOrderItemDraft(typeKey = COMPONENT_TYPES[0].key) {
  const definition = getComponentDefinition(typeKey)
  return {
    typeKey: definition.key,
    quantity: 1,
    dimensions: definition.fields.reduce((acc, field) => {
      acc[field.key] = ''
      return acc
    }, {}),
  }
}

export function buildSizeSpec(typeKey, dimensions) {
  const definition = getComponentDefinition(typeKey)
  const values = definition.fields.map((field) => String(dimensions?.[field.key] || '').trim())
  return values.join('*')
}

export function buildOrderItemPayload(draft, index) {
  const definition = getComponentDefinition(draft.typeKey)
  return {
    item_no: String(index + 1),
    component_type: definition.key,
    part_name: definition.label,
    spec: buildSizeSpec(draft.typeKey, draft.dimensions),
    dimensions: { ...draft.dimensions },
    quantity: Number(draft.quantity || 0),
    unit: '件',
    remark: '',
  }
}

export function isSupportedComponentType(typeKey) {
  return COMPONENT_TYPES.some((item) => item.key === typeKey)
}
