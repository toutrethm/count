function callUni(method, payload) {
  if (typeof uni !== 'undefined' && typeof uni[method] === 'function') {
    uni[method](payload)
    return true
  }
  return false
}

export function redirectTo(url) {
  if (callUni('redirectTo', { url })) {
    return
  }
  if (typeof window !== 'undefined' && typeof window.location !== 'undefined') {
    window.location.href = url
  }
}

export function navigateTo(url) {
  if (callUni('navigateTo', { url })) {
    return
  }
  if (typeof window !== 'undefined' && typeof window.location !== 'undefined') {
    window.location.href = url
  }
}

export function switchTabTo(url) {
  if (callUni('switchTab', { url })) {
    return
  }
  if (typeof window !== 'undefined' && typeof window.location !== 'undefined') {
    window.location.href = url
  }
}
