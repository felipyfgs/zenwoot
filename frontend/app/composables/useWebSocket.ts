type WsHandler = (data: unknown) => void

export function useWebSocket() {
  const config = useRuntimeConfig()
  const wsBase = (config.public.apiBase as string).replace(/^http/, 'ws')

  let socket: WebSocket | null = null
  const handlers: WsHandler[] = []
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null

  function connect(accountId = 'default') {
    if (socket?.readyState === WebSocket.OPEN) return

    const url = `${wsBase}/ws?accountId=${accountId}`
    socket = new WebSocket(url)

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        handlers.forEach(h => h(data))
      } catch { /* ignore malformed */ }
    }

    socket.onclose = () => {
      reconnectTimer = setTimeout(() => connect(accountId), 3000)
    }

    socket.onerror = () => {
      socket?.close()
    }
  }

  function disconnect() {
    if (reconnectTimer) clearTimeout(reconnectTimer)
    socket?.close()
    socket = null
  }

  function onMessage(handler: WsHandler) {
    handlers.push(handler)
  }

  onUnmounted(disconnect)

  return { connect, disconnect, onMessage }
}
