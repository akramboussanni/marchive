export interface User {
  id: string
  username: string
  role: string
  created_at: string
  request_credits: number
  invite_tokens: number
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  token: string
  username: string
  password: string
}
