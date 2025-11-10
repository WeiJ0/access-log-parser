/**
 * 搜尋服務
 * 文件路徑: frontend/src/services/searchService.ts
 * 目的: 提供高效的日誌記錄搜尋功能，支援 IP、URL、User-Agent 等欄位搜尋
 */

import { models } from '../../wailsjs/wailsjs/go/models'

// 使用 Wails 生成的類型
type LogEntry = models.LogEntry

/**
 * 搜尋條件介面
 */
export interface SearchCriteria {
  ip?: string           // IP 地址搜尋（支援部分符合）
  url?: string          // URL 路徑搜尋（支援部分符合）
  userAgent?: string    // User-Agent 搜尋（支援部分符合）
  method?: string       // HTTP 方法搜尋
  user?: string         // 使用者名稱搜尋
  keyword?: string      // 通用關鍵字搜尋（搜尋所有文字欄位）
  caseSensitive?: boolean // 是否區分大小寫（預設：false）
}

/**
 * 搜尋服務類別
 * 提供高效的日誌記錄搜尋功能
 */
export class SearchService {
  /**
   * 搜尋日誌記錄
   * @param entries 日誌記錄陣列
   * @param criteria 搜尋條件
   * @returns 符合條件的日誌記錄陣列
   * 
   * 效能要求: 100 萬筆記錄搜尋 ≤100ms
   * 複雜度: O(n) - 線性掃描，使用字串包含匹配
   */
  search(entries: LogEntry[], criteria: SearchCriteria): LogEntry[] {
    // 如果沒有搜尋條件，返回所有資料
    if (this.isEmptyCriteria(criteria)) {
      return entries
    }

    const caseSensitive = criteria.caseSensitive ?? false

    return entries.filter(entry => {
      // 檢查 IP 搜尋
      if (criteria.ip && !this.matchField(entry.ip, criteria.ip, caseSensitive)) {
        return false
      }

      // 檢查 URL 搜尋
      if (criteria.url && !this.matchField(entry.url, criteria.url, caseSensitive)) {
        return false
      }

      // 檢查 User-Agent 搜尋
      if (criteria.userAgent && !this.matchField(entry.userAgent, criteria.userAgent, caseSensitive)) {
        return false
      }

      // 檢查 HTTP 方法搜尋
      if (criteria.method && !this.matchField(entry.method, criteria.method, caseSensitive)) {
        return false
      }

      // 檢查使用者名稱搜尋
      if (criteria.user && entry.user && !this.matchField(entry.user, criteria.user, caseSensitive)) {
        return false
      }

      // 檢查通用關鍵字搜尋（搜尋所有文字欄位）
      if (criteria.keyword) {
        const keyword = caseSensitive ? criteria.keyword : criteria.keyword.toLowerCase()
        
        const searchableFields = [
          entry.ip,
          entry.user || '',
          entry.method,
          entry.url,
          entry.protocol,
          entry.referer,
          entry.userAgent
        ]

        const matched = searchableFields.some(field => {
          if (!field) return false
          const fieldValue = caseSensitive ? field : field.toLowerCase()
          return fieldValue.includes(keyword)
        })

        if (!matched) {
          return false
        }
      }

      // 所有條件都符合
      return true
    })
  }

  /**
   * 檢查搜尋條件是否為空
   * @param criteria 搜尋條件
   * @returns 如果所有條件都未設定或為空字串，返回 true
   */
  private isEmptyCriteria(criteria: SearchCriteria): boolean {
    const { ip, url, userAgent, method, user, keyword } = criteria
    
    return !ip && !url && !userAgent && !method && !user && !keyword
  }

  /**
   * 欄位匹配檢查（支援部分符合和大小寫設定）
   * @param fieldValue 欄位值
   * @param searchValue 搜尋值
   * @param caseSensitive 是否區分大小寫
   * @returns 如果符合返回 true
   * 
   * 實作細節:
   * - 使用字串的 includes() 方法進行部分匹配
   * - 支援大小寫敏感和不敏感搜尋
   * - 自動過濾空字串搜尋
   */
  private matchField(fieldValue: string, searchValue: string, caseSensitive: boolean): boolean {
    // 空字串搜尋視為無條件，返回 true
    if (!searchValue || searchValue.trim() === '') {
      return true
    }

    // 如果不區分大小寫，將兩者都轉為小寫
    const field = caseSensitive ? fieldValue : fieldValue.toLowerCase()
    const search = caseSensitive ? searchValue : searchValue.toLowerCase()

    return field.includes(search)
  }

  /**
   * 高亮搜尋結果中的匹配文字
   * @param text 原始文字
   * @param searchText 搜尋文字
   * @param caseSensitive 是否區分大小寫
   * @returns 包含高亮標記的 HTML 字串
   * 
   * 用途: 在 UI 中顯示搜尋結果時高亮匹配的文字
   */
  highlightMatch(text: string, searchText: string, caseSensitive: boolean = false): string {
    if (!searchText || searchText.trim() === '') {
      return text
    }

    // 轉義特殊字元，避免被當作正規表達式處理
    const escapedSearch = searchText.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
    
    // 建立正規表達式（支援大小寫設定）
    const flags = caseSensitive ? 'g' : 'gi'
    const regex = new RegExp(escapedSearch, flags)

    // 替換匹配的文字為帶有高亮標記的版本
    return text.replace(regex, (match) => `<mark>${match}</mark>`)
  }

  /**
   * 取得搜尋結果統計
   * @param totalCount 總記錄數
   * @param resultCount 搜尋結果數
   * @returns 統計資訊
   */
  getSearchStats(totalCount: number, resultCount: number): {
    total: number
    results: number
    percentage: number
  } {
    const percentage = totalCount > 0 ? (resultCount / totalCount) * 100 : 0

    return {
      total: totalCount,
      results: resultCount,
      percentage: Math.round(percentage * 100) / 100 // 保留兩位小數
    }
  }
}

/**
 * 匯出單例實例（方便直接使用）
 */
export const searchService = new SearchService()
