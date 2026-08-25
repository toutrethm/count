export function resolveApiBaseURL(explicitBaseURL) {
  if (explicitBaseURL) {
    return explicitBaseURL
  }

  if (typeof window !== 'undefined' && window.location && window.location.hostname) {
    const protocol = window.location.protocol || 'http:'
    return `${protocol}//${window.location.hostname}:8080`
  }

  return 'http://localhost:8080'
}
