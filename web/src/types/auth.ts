/**
 * 认证相关类型定义
 */

/**
 * 登录请求类型
 */
export interface LoginRequest {
    username: string
    password: string
}

/**
 * 登录响应类型
 */
export interface LoginResponse {
    access_token: string
    refresh_token: string
    access_expires_at: string
    refresh_expires_at: string
}

/**
 * 用户信息类型
 */
export interface UserInfo {
    username: string
    created_at?: string
    updated_at?: string
}

/**
 * 会话信息类型
 */
export interface SessionInfo {
    id: number
    is_active: boolean
    client_ip: string
    user_agent: string
    expires_at: string
    created_at: string
    last_access_at: string
}

/**
 * 会话列表响应类型
 */
export interface SessionListResponse {
    sessions: SessionInfo[]
    total: number
}

/**
 * 修改密码请求类型
 */
export interface ChangePasswordRequest {
    username: string
    old_password: string
    new_password: string
}

/**
 * 更新用户信息请求类型
 */
export interface UpdateUserInfoRequest {
    username: string
}

/**
 * 刷新令牌请求类型
 */
export interface RefreshTokenRequest {
    refresh_token: string
} 