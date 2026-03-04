/**
 * 更新相关类型定义
 */

export interface LatestInfo {
  /** 版本标签 */
  tag_name: string
  /** 发布时间 */
  published_at: string
  /** 更新内容 */
  body: string
}

export interface UpdateResponse {
  bestsub: LatestInfo
  webui: LatestInfo
  subconverter: LatestInfo
}

export type UpdateComponent = 'subconverter' | 'webui'

export interface SystemVersion {
  /** 版本号 */
  version: string
  /** 构建时间 */
  buildTime: string
  /** Git 提交哈希 */
  commit: string
  /** 作者 */
  author: string
  /** 仓库地址 */
  repo: string
  /** Subconverter 版本 */
  subconverter_version: string
}