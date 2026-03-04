import { API_PATH } from '../config/config'
import type { LoginResponse } from '../../types'

const TOKEN_KEYS = {
    ACCESS_TOKEN: 'access_token',
    REFRESH_TOKEN: 'refresh_token',
    ACCESS_EXPIRES: 'access_expires_at',
    REFRESH_EXPIRES: 'refresh_expires_at',
} as const

export class TokenManager {
    static isServer() {
        return typeof window === 'undefined'
    }

    static getTokens() {
        if (this.isServer()) return null

        const accessToken = localStorage.getItem(TOKEN_KEYS.ACCESS_TOKEN)
        const refreshToken = localStorage.getItem(TOKEN_KEYS.REFRESH_TOKEN)

        if (!accessToken || !refreshToken) return null

        return {
            access_token: accessToken,
            refresh_token: refreshToken,
            access_expires_at: localStorage.getItem(TOKEN_KEYS.ACCESS_EXPIRES) || '',
            refresh_expires_at: localStorage.getItem(TOKEN_KEYS.REFRESH_EXPIRES) || '',
        }
    }

    static setTokens(tokens: LoginResponse) {
        if (this.isServer()) return

        localStorage.setItem(TOKEN_KEYS.ACCESS_TOKEN, tokens.access_token)
        localStorage.setItem(TOKEN_KEYS.REFRESH_TOKEN, tokens.refresh_token)
        localStorage.setItem(TOKEN_KEYS.ACCESS_EXPIRES, tokens.access_expires_at)
        localStorage.setItem(TOKEN_KEYS.REFRESH_EXPIRES, tokens.refresh_expires_at)
    }

    static clearTokens() {
        if (this.isServer()) return
        Object.values(TOKEN_KEYS).forEach(key => localStorage.removeItem(key))
    }

    static isExpired(expiresAt: string): boolean {
        if (!expiresAt) return true
        return new Date(expiresAt) <= new Date()
    }

    static async getValidToken(): Promise<string | null> {
        try {
            const tokens = this.getTokens()
            if (!tokens) return null

            // 检查 refresh token 是否过期
            if (this.isExpired(tokens.refresh_expires_at)) {
                this.clearTokens()
                return null
            }

            // 如果 access token 未过期，直接返回
            if (!this.isExpired(tokens.access_expires_at)) {
                return tokens.access_token
            }

            // access token 过期，直接调用refresh API避免循环依赖
            const response = await fetch(`${API_PATH.base}${API_PATH.auth.refresh}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh_token: tokens.refresh_token })
            })

            if (!response.ok) {
                this.clearTokens()
                return null
            }

            const result = await response.json()
            this.setTokens(result.data)
            return result.data.access_token
        } catch (error) {
            console.error('Token refresh failed:', error)
            this.clearTokens()
            return null
        }
    }
}

// 导出实例方法给外部使用
export const tokenManager = {
    getTokens: () => TokenManager.getTokens(),
    setTokens: (tokens: LoginResponse) => TokenManager.setTokens(tokens),
    clearTokens: () => TokenManager.clearTokens(),
    isExpired: (expiresAt: string) => TokenManager.isExpired(expiresAt),
    getValidToken: () => TokenManager.getValidToken(),
}
