// API 服務封裝層 - 封裝 Wails Go ↔ TypeScript 通訊
// 文件路徑: frontend/src/services/logService.ts

import type {
  SelectFileResponse,
  ParseFileRequest,
  ParseFileResponse,
  ValidateFormatRequest,
  ValidateFormatResponse,
  LogFile,
} from '../types/log'

// 臨時 mock - 正式版本會從 wailsjs/go/internal/app/App.js 導入
// 這些函式會在 wails dev 或 wails build 時自動生成
declare global {
  interface Window {
    go?: {
      internal?: {
        app?: {
          App?: {
            SelectFile: () => Promise<SelectFileResponse>
            ParseFile: (req: ParseFileRequest) => Promise<ParseFileResponse>
            ValidateLogFormat: (req: ValidateFormatRequest) => Promise<ValidateFormatResponse>
            GetOpenFiles: () => Promise<LogFile[]>
            SetActiveFile: (fileId: string) => Promise<void>
            GetActiveFile: () => Promise<string>
            CloseFile: (fileId: string) => Promise<void>
          }
        }
      }
    }
  }
}

/**
 * LogService - 封裝所有與後端的通訊邏輯
 * 提供類型安全的 API 呼叫介面
 */
class LogService {
  private app = window.go?.internal?.app?.App

  /**
   * 開啟檔案選擇對話框
   * @returns 選擇的檔案資訊或錯誤
   */
  async selectFile(): Promise<SelectFileResponse> {
    if (!this.app) {
      throw new Error('Wails runtime 尚未初始化')
    }
    
    try {
      return await this.app.SelectFile()
    } catch (error) {
      console.error('選擇檔案失敗:', error)
      return {
        success: false,
        error: error instanceof Error ? error.message : '未知錯誤'
      }
    }
  }

  /**
   * 解析日誌檔案
   * @param filePath 檔案路徑
   * @param fileSize 檔案大小（位元組）
   * @returns 解析結果或錯誤
   */
  async parseFile(filePath: string, fileSize: number): Promise<ParseFileResponse> {
    if (!this.app) {
      throw new Error('Wails runtime 尚未初始化')
    }
    
    try {
      return await this.app.ParseFile({ filePath, fileSize })
    } catch (error) {
      console.error('解析檔案失敗:', error)
      return {
        success: false,
        error: error instanceof Error ? error.message : '未知錯誤'
      }
    }
  }

  /**
   * 驗證日誌檔案格式
   * @param filePath 檔案路徑
   * @param sampleLines 取樣行數（預設 100）
   * @returns 驗證結果
   */
  async validateFormat(filePath: string, sampleLines: number = 100): Promise<ValidateFormatResponse> {
    if (!this.app) {
      throw new Error('Wails runtime 尚未初始化')
    }
    
    try {
      return await this.app.ValidateLogFormat({ filePath, sampleLines })
    } catch (error) {
      console.error('驗證格式失敗:', error)
      return {
        success: false,
        isValid: false,
        error: error instanceof Error ? error.message : '未知錯誤'
      }
    }
  }

  /**
   * 取得所有開啟的檔案列表
   * @returns 開啟的檔案陣列
   */
  async getOpenFiles(): Promise<LogFile[]> {
    if (!this.app) {
      throw new Error('Wails runtime 尚未初始化')
    }
    
    try {
      return await this.app.GetOpenFiles()
    } catch (error) {
      console.error('取得開啟檔案失敗:', error)
      return []
    }
  }

  /**
   * 設定當前活動的檔案
   * @param fileId 檔案 ID
   */
  async setActiveFile(fileId: string): Promise<void> {
    if (!this.app) {
      throw new Error('Wails runtime 尚未初始化')
    }
    
    try {
      await this.app.SetActiveFile(fileId)
    } catch (error) {
      console.error('設定活動檔案失敗:', error)
      throw error
    }
  }

  /**
   * 取得當前活動的檔案 ID
   * @returns 檔案 ID
   */
  async getActiveFile(): Promise<string> {
    if (!this.app) {
      throw new Error('Wails runtime 尚未初始化')
    }
    
    try {
      return await this.app.GetActiveFile()
    } catch (error) {
      console.error('取得活動檔案失敗:', error)
      return ''
    }
  }

  /**
   * 關閉指定的檔案
   * @param fileId 檔案 ID
   */
  async closeFile(fileId: string): Promise<void> {
    if (!this.app) {
      throw new Error('Wails runtime 尚未初始化')
    }
    
    try {
      await this.app.CloseFile(fileId)
    } catch (error) {
      console.error('關閉檔案失敗:', error)
      throw error
    }
  }
}

// 導出單例實例
export const logService = new LogService()
export default logService
