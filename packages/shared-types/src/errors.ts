/**
 * Centralized error codes. Keep in sync with backend/internal/shared/errors/errors.go.
 *
 * 10xxx  Generic
 * 20xxx  Auth
 * 30xxx  User / 31xxx Catalog / 32xxx Order / 33xxx Payment
 * 34xxx  Chat / 35xxx Review, Aftersales
 * 40xxx  Admin
 */
export const ErrorCode = {
  OK: 0,
  PARAM_INVALID: 10001,
  SYSTEM_ERROR: 10002,

  UNAUTHORIZED: 20001,
  TOKEN_EXPIRED: 20002,
  FORBIDDEN: 20003,

  USER_NOT_FOUND: 30001,
  CAR_NOT_FOUND: 31001,
  CAR_ALREADY_SOLD: 31002,
  ORDER_STATE_INVALID: 32001,
  ORDER_TIMEOUT: 32002,
  PAYMENT_FAILED: 33001,
  REFUND_FAILED: 33002,
} as const

export type ErrorCodeType = (typeof ErrorCode)[keyof typeof ErrorCode]
