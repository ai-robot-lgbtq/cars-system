/**
 * Standard API response envelope.
 */
export interface APIResponse<T = unknown> {
  code: number
  message: string
  data?: T
}

/**
 * Pagination result wrapper.
 */
export interface PageResult<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

/**
 * Pagination query params accepted by list endpoints.
 */
export interface PageQuery {
  page?: number
  page_size?: number
}
