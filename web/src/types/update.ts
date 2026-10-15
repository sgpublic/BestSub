/**
 * 更新相关类型定义
 */

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
