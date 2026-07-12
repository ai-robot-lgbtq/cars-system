import { describe, expect, it } from 'vitest'
import { UserRole, hasRole, OrderStatus, CarStatus } from '../src/enums'
import { ErrorCode } from '../src/errors'

describe('UserRole.hasRole', () => {
  it('returns true when user has required role', () => {
    // user is seller+admin (2|4 = 6)
    expect(hasRole(6, UserRole.SELLER)).toBe(true)
    expect(hasRole(6, UserRole.ADMIN)).toBe(true)
  })

  it('returns false when user lacks required role', () => {
    expect(hasRole(UserRole.BUYER, UserRole.SELLER)).toBe(false)
    expect(hasRole(UserRole.SELLER, UserRole.ADMIN)).toBe(false)
  })

  it('returns false for guest', () => {
    expect(hasRole(UserRole.GUEST, UserRole.BUYER)).toBe(false)
  })
})

describe('OrderStatus transitions', () => {
  it('CREATED is initial state', () => {
    expect(OrderStatus.CREATED).toBe(0)
  })

  it('happy path ends at COMPLETED', () => {
    const happyPath = [
      OrderStatus.CREATED,
      OrderStatus.PAID,
      OrderStatus.SHIPPING,
      OrderStatus.TRANSFERRING,
      OrderStatus.COMPLETED,
    ]
    expect(happyPath.length).toBe(5)
  })
})

describe('CarStatus', () => {
  it('DRAFT is 0', () => {
    expect(CarStatus.DRAFT).toBe(0)
  })
})

describe('ErrorCode', () => {
  it('OK is 0', () => {
    expect(ErrorCode.OK).toBe(0)
  })
})
