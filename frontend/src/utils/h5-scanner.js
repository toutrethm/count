import jsQR from 'jsqr'

export function isH5CameraScanSupported(env = {}) {
  const navigatorObject = env.navigator || (typeof navigator !== 'undefined' ? navigator : null)
  const windowObject = env.window || (typeof window !== 'undefined' ? window : null)
  return Boolean(
    windowObject &&
      navigatorObject &&
      navigatorObject.mediaDevices &&
      typeof navigatorObject.mediaDevices.getUserMedia === 'function'
  )
}

export function extractScanText(payload) {
  if (typeof payload === 'string') {
    return payload.trim()
  }

  if (!payload || typeof payload !== 'object') {
    return ''
  }

  const candidates = [payload.text, payload.result, payload.data]
  for (const value of candidates) {
    if (typeof value === 'string' && value.trim()) {
      return value.trim()
    }
  }

  return ''
}

export function createH5CameraScanSession({ videoEl, canvasEl, facingMode = 'environment' }) {
  if (!videoEl || !canvasEl) {
    throw new Error('missing scanner elements')
  }

  const state = {
    active: true,
    finished: false,
    rafId: 0,
    stream: null,
    reject: null,
  }

  const ctx = canvasEl.getContext('2d', { willReadFrequently: true })
  if (!ctx) {
    throw new Error('cannot get canvas context')
  }

  const cleanup = () => {
    if (state.rafId) {
      cancelAnimationFrame(state.rafId)
      state.rafId = 0
    }
    if (state.stream) {
      state.stream.getTracks().forEach((track) => track.stop())
      state.stream = null
    }
    if (videoEl.srcObject) {
      videoEl.srcObject = null
    }
  }

  const fail = (error) => {
    if (state.finished) {
      return
    }
    state.finished = true
    cleanup()
    if (state.reject) {
      state.reject(error)
    }
  }

  const promise = new Promise((resolve, reject) => {
    state.reject = reject

    const begin = async () => {
      try {
        if (!isH5CameraScanSupported()) {
          throw new Error('浏览器不支持摄像头扫码')
        }

        state.stream = await navigator.mediaDevices.getUserMedia({
          audio: false,
          video: {
            facingMode: { ideal: facingMode },
          },
        })

        videoEl.srcObject = state.stream
        if (typeof videoEl.play === 'function') {
          await videoEl.play()
        }

        const scanFrame = () => {
          if (!state.active || state.finished) {
            return
          }

          if (videoEl.readyState >= 2 && videoEl.videoWidth > 0 && videoEl.videoHeight > 0) {
            canvasEl.width = videoEl.videoWidth
            canvasEl.height = videoEl.videoHeight
            ctx.drawImage(videoEl, 0, 0, canvasEl.width, canvasEl.height)
            const imageData = ctx.getImageData(0, 0, canvasEl.width, canvasEl.height)
            const result = jsQR(imageData.data, canvasEl.width, canvasEl.height, {
              inversionAttempts: 'dontInvert',
            })
            const text = extractScanText(result?.data)
            if (text) {
              state.finished = true
              cleanup()
              resolve({ text, raw: result })
              return
            }
          }

          state.rafId = window.requestAnimationFrame(scanFrame)
        }

        scanFrame()
      } catch (error) {
        fail(error)
      }
    }

    begin()
  })

  return {
    promise,
    stop() {
      if (state.finished) {
        return
      }
      state.active = false
      fail(Object.assign(new Error('H5 scanner cancelled'), { code: 'H5_SCAN_CANCELLED' }))
    },
  }
}
