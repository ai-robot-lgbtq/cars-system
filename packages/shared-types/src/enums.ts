/**
 * User roles. Stored as bit flags in DB (1=buyer, 2=seller, 4=admin).
 */
export enum UserRole {
  GUEST = 0,
  BUYER = 1,
  SELLER = 2,
  ADMIN = 4,
}

export function hasRole(userRole: number, required: UserRole): boolean {
  return (userRole & required) === required
}

/**
 * Car status.
 */
export enum CarStatus {
  DRAFT = 0,
  PENDING = 1,
  ONLINE = 2,
  SOLD = 3,
  OFFLINE = 4,
}

/**
 * Order status.
 */
export enum OrderStatus {
  CREATED = 0,
  PAID = 1,
  SHIPPING = 2,
  TRANSFERRING = 3,
  COMPLETED = 4,
  CANCELLED = 5,
  REFUNDED = 6,
}

/**
 * Payment status.
 */
export enum PaymentStatus {
  PENDING = 0,
  SUCCESS = 1,
  FAILED = 2,
  REFUNDED = 3,
}
