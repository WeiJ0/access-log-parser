// TypeScript 型別定義，對應後端 Go 結構
// 文件路徑: frontend/src/types/log.ts

/**
 * LogEntry 表示單筆日誌記錄
 * 對應 Go: internal/models/log_entry.go
 */
export interface LogEntry {
  lineNumber: number
  ip: string
  user: string
  timestamp: string // ISO 8601 格式
  method: string
  url: string
  protocol: string
  statusCode: number
  responseBytes: number
  referer: string
  userAgent: string
  rawLine: string
}

/**
 * ParseError 表示解析錯誤
 * 對應 Go: internal/parser/parser.go
 */
export interface ParseError {
  lineNumber: number
  line: string
  error: string
}

/**
 * ParseResult 包含解析結果和統計資訊
 * 對應 Go: internal/parser/parser.go
 */
export interface ParseResult {
  entries: LogEntry[]
  totalLines: number
  parsedLines: number
  errorLines: number
  errorSamples: ParseError[]
  parseTime: number // 毫秒
  memoryUsed: number // 位元組
  throughputMB: number // MB/秒
}

/**
 * Statistics 表示統計資料
 * 對應 Go: internal/models/statistics.go
 */
export interface Statistics {
  // 基本統計
  totalRequests: number
  uniqueIPs: number
  totalBytes: number
  
  // 時間範圍
  startTime: string
  endTime: string
  
  // 狀態碼分布
  statusCodeDist: Record<number, number>
  statusCatDist: Record<string, number>
  
  // URL 統計
  topURLs: URLStat[]
  
  // 時間分布
  hourlyDist: Record<number, number>
  dailyDist: Record<string, number>
  
  // IP 統計
  topIPs: IPStat[]
  
  // HTTP 方法分布
  methodDist: Record<string, number>
  
  // 回應大小統計
  avgResponseSize: number
  minResponseSize: number
  maxResponseSize: number
  
  // 錯誤統計
  errorCount: number
  clientErrorCount: number
  serverErrorCount: number
  errorRate: number
  topErrorURLs: URLStat[]
  topErrorIPs: IPStat[]
  
  // User Agent 統計
  topUserAgents: UserAgentStat[]
  browserDist: Record<string, number>
  osDist: Record<string, number>
  
  // Referer 統計
  topReferers: RefererStat[]
}

/**
 * URLStat 表示 URL 統計資料
 */
export interface URLStat {
  url: string
  count: number
  totalBytes: number
  errorCount: number
  errorRate: number
  avgBytes: number
}

/**
 * IPStat 表示 IP 統計資料
 */
export interface IPStat {
  ip: string
  count: number
  totalBytes: number
  errorCount: number
  errorRate: number
  uniqueUrls: number
}

/**
 * UserAgentStat 表示 User Agent 統計資料
 */
export interface UserAgentStat {
  userAgent: string
  count: number
  percentage: number
}

/**
 * RefererStat 表示 Referer 統計資料
 */
export interface RefererStat {
  referer: string
  count: number
  percentage: number
}

/**
 * LogFile 表示開啟的日誌檔案
 * 對應 Go: internal/models/log_file.go
 */
export interface LogFile {
  id: string
  path: string
  name: string
  size: number
  parseResult?: ParseResult
  statistics?: Statistics
  isLoading: boolean
  error?: string
}

/**
 * SelectFileResponse 檔案選擇 API 回應
 */
export interface SelectFileResponse {
  success: boolean
  filePath?: string
  fileName?: string
  fileSize?: number
  error?: string
}

/**
 * ParseFileRequest 解析檔案 API 請求
 */
export interface ParseFileRequest {
  filePath: string
  fileSize: number
}

/**
 * ParseFileResponse 解析檔案 API 回應
 */
export interface ParseFileResponse {
  success: boolean
  result?: ParseResult
  error?: string
}

/**
 * ValidateFormatRequest 驗證格式 API 請求
 */
export interface ValidateFormatRequest {
  filePath: string
  sampleLines: number
}

/**
 * ValidateFormatResponse 驗證格式 API 回應
 */
export interface ValidateFormatResponse {
  success: boolean
  isValid: boolean
  format?: string
  error?: string
}
