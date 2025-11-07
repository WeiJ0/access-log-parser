/**
 * 篩選器服務
 * 文件路徑: frontend/src/services/filterService.ts
 * 目的: 提供高效的日誌記錄篩選功能，支援狀態碼、時間範圍、HTTP 方法等篩選
 */

import { models } from '../../wailsjs/wailsjs/go/models'

// 使用 Wails 生成的類型
type LogEntry = models.LogEntry

/**
 * 狀態碼範圍介面
 */
export interface StatusCodeRange {
  min: number
  max: number
}

/**
 * 時間範圍介面
 */
export interface TimeRange {
  start?: string // ISO 8601 格式
  end?: string   // ISO 8601 格式
}

/**
 * 回應大小範圍介面
 */
export interface ResponseSizeRange {
  min?: number // 位元組
  max?: number // 位元組
}

/**
 * 篩選條件介面
 */
export interface FilterCriteria {
  statusCodes?: number[]              // 特定狀態碼列表（精確匹配）
  statusCodeRange?: StatusCodeRange   // 狀態碼範圍
  timeRange?: TimeRange               // 時間範圍
  methods?: string[]                  // HTTP 方法列表
  responseSizeRange?: ResponseSizeRange // 回應大小範圍
}

/**
 * 篩選統計資訊介面
 */
export interface FilterStats {
  total: number       // 總記錄數
  filtered: number    // 篩選後的記錄數
  percentage: number  // 篩選結果百分比
}

/**
 * 篩選器服務類別
 * 提供高效的日誌記錄篩選功能
 */
export class FilterService {
  /**
   * 篩選日誌記錄
   * @param entries 日誌記錄陣列
   * @param criteria 篩選條件
   * @returns 符合條件的日誌記錄陣列
   * 
   * 效能要求: 100 萬筆記錄篩選 ≤100ms
   * 複雜度: O(n) - 線性掃描，所有條件在單次遍歷中完成
   */
  filter(entries: LogEntry[], criteria: FilterCriteria): LogEntry[] {
    // 如果沒有篩選條件，返回所有資料
    if (this.isEmptyCriteria(criteria)) {
      return entries
    }

    return entries.filter(entry => {
      // 檢查特定狀態碼列表
      if (criteria.statusCodes && criteria.statusCodes.length > 0) {
        if (!criteria.statusCodes.includes(entry.statusCode)) {
          return false
        }
      }

      // 檢查狀態碼範圍
      if (criteria.statusCodeRange) {
        const { min, max } = criteria.statusCodeRange
        
        // 驗證範圍有效性
        if (min > max) {
          return false
        }
        
        if (entry.statusCode < min || entry.statusCode > max) {
          return false
        }
      }

      // 檢查時間範圍
      if (criteria.timeRange) {
        const { start, end } = criteria.timeRange
        
        // 驗證時間範圍有效性
        if (start && end && start > end) {
          return false
        }
        
        if (start && entry.timestamp < start) {
          return false
        }
        
        if (end && entry.timestamp > end) {
          return false
        }
      }

      // 檢查 HTTP 方法列表
      if (criteria.methods && criteria.methods.length > 0) {
        if (!criteria.methods.includes(entry.method)) {
          return false
        }
      }

      // 檢查回應大小範圍
      if (criteria.responseSizeRange) {
        const { min, max } = criteria.responseSizeRange
        
        // 驗證範圍有效性
        if (min !== undefined && max !== undefined && min > max) {
          return false
        }
        
        if (min !== undefined && entry.responseBytes < min) {
          return false
        }
        
        if (max !== undefined && entry.responseBytes > max) {
          return false
        }
      }

      // 所有條件都符合
      return true
    })
  }

  /**
   * 檢查篩選條件是否為空
   * @param criteria 篩選條件
   * @returns 如果所有條件都未設定，返回 true
   */
  private isEmptyCriteria(criteria: FilterCriteria): boolean {
    const {
      statusCodes,
      statusCodeRange,
      timeRange,
      methods,
      responseSizeRange
    } = criteria

    return (
      (!statusCodes || statusCodes.length === 0) &&
      !statusCodeRange &&
      !timeRange &&
      (!methods || methods.length === 0) &&
      !responseSizeRange
    )
  }

  /**
   * 取得篩選結果統計
   * @param originalEntries 原始日誌記錄陣列
   * @param filteredEntries 篩選後的日誌記錄陣列
   * @returns 篩選統計資訊
   */
  getFilterStats(originalEntries: LogEntry[], filteredEntries: LogEntry[]): FilterStats {
    const total = originalEntries.length
    const filtered = filteredEntries.length
    const percentage = total > 0 ? (filtered / total) * 100 : 0

    return {
      total,
      filtered,
      percentage: Math.round(percentage * 100) / 100 // 保留兩位小數
    }
  }

  /**
   * 建立預定義的狀態碼範圍篩選
   * 提供常用的 HTTP 狀態碼類別快捷方式
   */
  static readonly STATUS_CODE_RANGES = {
    /** 2xx 成功 */
    SUCCESS: { min: 200, max: 299 } as StatusCodeRange,
    
    /** 3xx 重新導向 */
    REDIRECT: { min: 300, max: 399 } as StatusCodeRange,
    
    /** 4xx 客戶端錯誤 */
    CLIENT_ERROR: { min: 400, max: 499 } as StatusCodeRange,
    
    /** 5xx 伺服器錯誤 */
    SERVER_ERROR: { min: 500, max: 599 } as StatusCodeRange,
    
    /** 所有錯誤 (4xx + 5xx) */
    ALL_ERRORS: { min: 400, max: 599 } as StatusCodeRange
  }

  /**
   * 建立預定義的時間範圍篩選
   * @param referenceTime 參考時間（預設為當前時間）
   * @returns 時間範圍快捷方式物件
   */
  static createTimeRanges(referenceTime: Date = new Date()) {
    const now = referenceTime
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
    const yesterday = new Date(today)
    yesterday.setDate(yesterday.getDate() - 1)
    const lastWeek = new Date(today)
    lastWeek.setDate(lastWeek.getDate() - 7)
    const lastMonth = new Date(today)
    lastMonth.setMonth(lastMonth.getMonth() - 1)

    return {
      /** 今天 */
      TODAY: {
        start: today.toISOString(),
        end: now.toISOString()
      } as TimeRange,
      
      /** 昨天 */
      YESTERDAY: {
        start: yesterday.toISOString(),
        end: today.toISOString()
      } as TimeRange,
      
      /** 過去 7 天 */
      LAST_7_DAYS: {
        start: lastWeek.toISOString(),
        end: now.toISOString()
      } as TimeRange,
      
      /** 過去 30 天 */
      LAST_30_DAYS: {
        start: lastMonth.toISOString(),
        end: now.toISOString()
      } as TimeRange
    }
  }

  /**
   * 組合搜尋和篩選條件
   * 允許同時使用搜尋服務和篩選器服務
   * 
   * @param entries 日誌記錄陣列
   * @param filterCriteria 篩選條件
   * @param searchCallback 搜尋回調函式（先執行篩選，再執行搜尋）
   * @returns 符合所有條件的日誌記錄陣列
   */
  filterAndSearch(
    entries: LogEntry[],
    filterCriteria: FilterCriteria,
    searchCallback: (entries: LogEntry[]) => LogEntry[]
  ): LogEntry[] {
    // 先執行篩選
    const filtered = this.filter(entries, filterCriteria)
    
    // 再執行搜尋
    return searchCallback(filtered)
  }

  /**
   * 驗證篩選條件的有效性
   * @param criteria 篩選條件
   * @returns 驗證結果和錯誤訊息
   */
  validateCriteria(criteria: FilterCriteria): {
    valid: boolean
    errors: string[]
  } {
    const errors: string[] = []

    // 驗證狀態碼範圍
    if (criteria.statusCodeRange) {
      const { min, max } = criteria.statusCodeRange
      
      if (min < 100 || min > 599) {
        errors.push('狀態碼最小值必須在 100-599 之間')
      }
      
      if (max < 100 || max > 599) {
        errors.push('狀態碼最大值必須在 100-599 之間')
      }
      
      if (min > max) {
        errors.push('狀態碼最小值不能大於最大值')
      }
    }

    // 驗證時間範圍
    if (criteria.timeRange) {
      const { start, end } = criteria.timeRange
      
      if (start && end && start > end) {
        errors.push('開始時間不能晚於結束時間')
      }
    }

    // 驗證回應大小範圍
    if (criteria.responseSizeRange) {
      const { min, max } = criteria.responseSizeRange
      
      if (min !== undefined && min < 0) {
        errors.push('回應大小最小值不能為負數')
      }
      
      if (max !== undefined && max < 0) {
        errors.push('回應大小最大值不能為負數')
      }
      
      if (min !== undefined && max !== undefined && min > max) {
        errors.push('回應大小最小值不能大於最大值')
      }
    }

    return {
      valid: errors.length === 0,
      errors
    }
  }
}

/**
 * 匯出單例實例（方便直接使用）
 */
export const filterService = new FilterService()
