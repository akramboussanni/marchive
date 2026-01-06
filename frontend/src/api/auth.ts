import apiClient from './client'
import type { User, LoginRequest, RegisterRequest } from '@/types/user'

export const authApi = {
  async login(credentials: LoginRequest) {
    const response = await apiClient.post('/auth/login', credentials)
    return response.data
  },

  async logout() {
    const response = await apiClient.post('/auth/logout')
    return response.data
  },

  async getMe(): Promise<User> {
    const response = await apiClient.get('/auth/me')
    return response.data
  },

  async register(data: RegisterRequest) {
    const response = await apiClient.post('/invites/use', data)
    return response.data
  },

  async changePassword(currentPassword: string, newPassword: string) {
    const response = await apiClient.post('/auth/change-password', {
      current_password: currentPassword,
      new_password: newPassword
    })
    return response.data
  }
}
