
const API_PORT = process.env.NEXT_PUBLIC_API_PORT

export function getApiPrefix() {
    return API_PORT && window ? `${window.location.protocol}//${window.location.hostname}:${API_PORT}` : ''
}

export function getWsPrefix() {
   return `${window?.location.protocol.startsWith('https') ? 'wss' : 'ws'}://${window?.location.hostname}:${API_PORT || window?.location.port}`
}