export function buildQrRectangles(modules) {
  const rects = []
  for (let y = 0; y < modules.length; y += 1) {
    for (let x = 0; x < modules[y].length; x += 1) {
      if (modules[y][x]?.isBlack) {
        rects.push({ x, y })
      }
    }
  }
  return rects
}
