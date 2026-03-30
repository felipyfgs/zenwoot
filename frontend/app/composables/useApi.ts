import type { ApiResponse, PaginatedResponse } from '~/types'

export function useApi() {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase as string

  function headers(): Record<string, string> {
    const h: Record<string, string> = { 'Content-Type': 'application/json' }
    // Use API key for authentication instead of token
    h['ApiKey'] = 'wzap-secret-apikey-123'
    return h
  }

  async function get<T>(path: string, query?: Record<string, unknown>): Promise<T> {
    const res = await $fetch<ApiResponse<T>>(`${baseURL}${path}`, {
      method: 'GET',
      headers: headers(),
      query
    })
    if (!res.success) throw new Error(res.error || res.message || 'Request failed')
    return res.data
  }

  async function post<T>(path: string, body?: unknown): Promise<T> {
    const res = await $fetch<ApiResponse<T>>(`${baseURL}${path}`, {
      method: 'POST',
      headers: headers(),
      body: body ? JSON.stringify(body) : undefined
    })
    if (!res.success) throw new Error(res.error || res.message || 'Request failed')
    return res.data
  }

  async function put<T>(path: string, body?: unknown): Promise<T> {
    const res = await $fetch<ApiResponse<T>>(`${baseURL}${path}`, {
      method: 'PUT',
      headers: headers(),
      body: body ? JSON.stringify(body) : undefined
    })
    if (!res.success) throw new Error(res.error || res.message || 'Request failed')
    return res.data
  }

  async function del<T = void>(path: string): Promise<T> {
    const res = await $fetch<ApiResponse<T>>(`${baseURL}${path}`, {
      method: 'DELETE',
      headers: headers()
    })
    if (!res.success) throw new Error(res.error || res.message || 'Request failed')
    return res.data
  }

  async function paginated<T>(path: string, query?: Record<string, unknown>): Promise<PaginatedResponse<T>> {
    const res = await $fetch<PaginatedResponse<T>>(`${baseURL}${path}`, {
      method: 'GET',
      headers: headers(),
      query
    })
    if (!res.success) throw new Error(res.error || res.message || 'Request failed')
    return res
  }

  return { get, post, put, del, paginated }
}
