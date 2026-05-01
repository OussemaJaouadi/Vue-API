export type GlobalRole = 'manager' | 'user'

export interface AuthUser {
  userId: string
  email: string
  username: string
  globalRole: GlobalRole
}

export interface AuthResponse extends AuthUser {
  accessToken: string
}

export interface AuthCredentials {
  login: string
  password: string
}

export interface RegisterPayload {
  email: string
  username: string
  password: string
}

export interface EventTicketResponse {
  ticket: string
}

export interface ApiEvent<T = unknown> {
  type: string
  data?: T
}
