const CAPABILITY_LABELS = {
  blanking_center: '切料和打中心孔',
  turn_outer: '车外圆',
  turn_head: '车大头',
  center_hole: '打中心孔',
  turn_head_center: '车大头和打中心孔',
  drill_tap_small: '钻孔攻牙',
  drill_tap_batch: '批量钻孔攻牙',
  turn_sleeve: '车套',
  turn_sleeve_auto: '车套（自动）',
  turn_sleeve_manual: '车套（非自动）',
}

export function formatCapabilityLabel(code) {
  const value = String(code || '').trim()
  if (!value) return '-'
  return CAPABILITY_LABELS[value] || value
}

export function formatCapabilityListLabel(value) {
  return (
    String(value || '')
      .split(',')
      .map((item) => formatCapabilityLabel(item))
      .filter((item) => item && item !== '-')
      .join('、') || '-'
  )
}
