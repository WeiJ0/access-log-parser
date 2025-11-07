// @ts-nocheck
/**
 * 搜尋服務測試
 * 文件路徑: frontend/src/services/searchService.test.ts
 * 目的: 測試各種搜尋條件（IP、路徑、User-Agent）
 */

import { describe, it, expect, beforeEach } from 'vitest'
import { SearchService, SearchCriteria } from './searchService'
import { models } from '../../wailsjs/wailsjs/go/models'

type LogEntry = models.LogEntry

describe('SearchService', () => {
  let service: SearchService
  let testData: LogEntry[]

  beforeEach(() => {
    service = new SearchService()
    
    // 準備測試資料
    testData = [
      new models.LogEntry({
        lineNumber: 1,
        ip: '192.168.1.100',
        user: '-',
        timestamp: '2024-01-01T10:00:00Z',
        method: 'GET',
        url: '/index.html',
        protocol: 'HTTP/1.1',
        statusCode: 200,
        responseBytes: 1024,
        referer: '-',
        userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',
        rawLine: '192.168.1.100 - - [01/Jan/2024:10:00:00 +0000] "GET /index.html HTTP/1.1" 200 1024 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"'
      }),
      new models.LogEntry({
        lineNumber: 2,
        ip: '10.0.0.50',
        user: '-',
        timestamp: '2024-01-01T10:01:00Z',
        method: 'POST',
        url: '/api/users',
        protocol: 'HTTP/1.1',
        statusCode: 201,
        responseBytes: 512,
        referer: 'https://example.com',
        userAgent: 'curl/7.68.0',
        rawLine: '10.0.0.50 - - [01/Jan/2024:10:01:00 +0000] "POST /api/users HTTP/1.1" 201 512 "https://example.com" "curl/7.68.0"'
      }),
      new models.LogEntry({
        lineNumber: 3,
        ip: '192.168.1.100',
        user: 'admin',
        timestamp: '2024-01-01T10:02:00Z',
        method: 'GET',
        url: '/admin/dashboard',
        protocol: 'HTTP/1.1',
        statusCode: 403,
        responseBytes: 256,
        referer: '-',
        userAgent: 'Googlebot/2.1 (+http://www.google.com/bot.html)',
        rawLine: '192.168.1.100 admin - [01/Jan/2024:10:02:00 +0000] "GET /admin/dashboard HTTP/1.1" 403 256 "-" "Googlebot/2.1"'
      }),
      new models.LogEntry({
        lineNumber: 4,
        ip: '172.16.0.1',
        user: '-',
        timestamp: '2024-01-01T10:03:00Z',
        method: 'GET',
        url: '/api/data',
        protocol: 'HTTP/1.1',
        statusCode: 500,
        responseBytes: 128,
        referer: '-',
        userAgent: 'python-requests/2.28.0',
        rawLine: '172.16.0.1 - - [01/Jan/2024:10:03:00 +0000] "GET /api/data HTTP/1.1" 500 128 "-" "python-requests/2.28.0"'
      })
    ]
  })

  describe('基本搜尋功能', () => {
    it('應該搜尋 IP 地址（完全符合）', () => {
      const criteria: SearchCriteria = {
        ip: '192.168.1.100'
      }
      
      const results = service.search(testData, criteria)
      
      expect(results).toHaveLength(2)
      expect(results[0].ip).toBe('192.168.1.100')
      expect(results[1].ip).toBe('192.168.1.100')
    })

    it('應該搜尋 IP 地址（部分符合）', () => {
      const criteria: SearchCriteria = {
        ip: '192.168'
      }
      
      const results = service.search(testData, criteria)
      
      expect(results).toHaveLength(2)
      expect(results[0].ip).toContain('192.168')
      expect(results[1].ip).toContain('192.168')
    })

    it('應該搜尋 URL 路徑（完全符合）', () => {
      const criteria: SearchCriteria = {
        url: '/api/users'
      }
      
      const results = service.search(testData, criteria)
      
      expect(results).toHaveLength(1)
      expect(results[0].url).toBe('/api/users')
    })

    it('應該搜尋 URL 路徑（部分符合）', () => {
      const criteria: SearchCriteria = {
        url: '/api/'
      }
      
      const results = service.search(testData, criteria)
      
      expect(results).toHaveLength(2)
      expect(results[0].url).toContain('/api/')
      expect(results[1].url).toContain('/api/')
    })

    it('應該搜尋 User-Agent（完全符合）', () => {
      const criteria: SearchCriteria = {
        userAgent: 'curl/7.68.0'
      }
      
      const results = service.search(testData, criteria)
      
      expect(results).toHaveLength(1)
      expect(results[0].userAgent).toBe('curl/7.68.0')
    })

    it('應該搜尋 User-Agent（部分符合）', () => {
      const criteria: SearchCriteria = {
        userAgent: 'bot'
      }
      
      const results = service.search(testData, criteria)
      
      expect(results).toHaveLength(1)
      expect(results[0].userAgent).toContain('bot')
    })
  })

  describe('不區分大小寫搜尋', () => {
    it('應該對 URL 搜尋不區分大小寫', () => {
      const criteria: SearchCriteria = {
        url: '/API/USERS',
        caseSensitive: false
      }
      
      const results = service.search(testData, criteria)
      
      expect(results).toHaveLength(1)
      expect(results[0].url.toLowerCase()).toContain('api/users')
    })

    it('應該對 User-Agent 搜尋不區分大小寫', () => {
      const criteria: SearchCriteria = {
        userAgent: 'CURL',
        caseSensitive: false
      }
      
      const results = service.search(testData, criteria)
      
      expect(results).toHaveLength(1)
    })

    it('應該支援區分大小寫搜尋', () => {
      const criteria: SearchCriteria = {
        url: '/API/',
        caseSensitive: true
      }
      
      const results = service.search(testData, criteria)
      
      expect(results).toHaveLength(0)
    })
  })

  describe('複合搜尋條件', () => {
    it('應該組合 IP 和 URL 搜尋（AND 邏輯）', () => {
      const criteria: SearchCriteria = {
        ip: '192.168.1.100',
        url: '/admin/'
      }
      
      const results = service.search(testData, criteria)
      
      expect(results).toHaveLength(1)
      expect(results[0].ip).toBe('192.168.1.100')
      expect(results[0].url).toContain('/admin/')
    })

    it('應該組合多個搜尋條件', () => {
      const criteria: SearchCriteria = {
        ip: '192.168.1.100',
        url: '/admin/',
        userAgent: 'bot'
      }
      
      const results = service.search(testData, criteria)
      
      expect(results).toHaveLength(1)
      expect(results[0].lineNumber).toBe(3)
    })

    it('應該在沒有符合的記錄時返回空陣列', () => {
      const criteria: SearchCriteria = {
        ip: '192.168.1.100',
        url: '/api/users' // 這個 IP 沒有訪問過這個 URL
      }
      
      const results = service.search(testData, criteria)
      
      expect(results).toHaveLength(0)
    })
  })

  describe('通用關鍵字搜尋', () => {
    it('應該在所有文字欄位中搜尋關鍵字', () => {
      const criteria: SearchCriteria = {
        keyword: 'admin'
      }
      
      const results = service.search(testData, criteria)
      
      // 應該找到 user='admin' 和 url='/admin/dashboard' 的記錄
      expect(results).toHaveLength(1)
      expect(results[0].lineNumber).toBe(3)
    })

    it('應該搜尋多個欄位的關鍵字', () => {
      const criteria: SearchCriteria = {
        keyword: 'api'
      }
      
      const results = service.search(testData, criteria)
      
      // 應該找到所有包含 'api' 的記錄
      expect(results.length).toBeGreaterThan(0)
      results.forEach((entry: LogEntry) => {
        const matched = entry.url.toLowerCase().includes('api') ||
                       entry.userAgent.toLowerCase().includes('api')
        expect(matched).toBe(true)
      })
    })
  })

  describe('邊界情況處理', () => {
    it('應該處理空的搜尋條件', () => {
      const criteria: SearchCriteria = {}
      
      const results = service.search(testData, criteria)
      
      // 沒有條件應該返回所有資料
      expect(results).toHaveLength(testData.length)
    })

    it('應該處理空的資料集', () => {
      const criteria: SearchCriteria = {
        ip: '192.168.1.100'
      }
      
      const results = service.search([], criteria)
      
      expect(results).toHaveLength(0)
    })

    it('應該處理特殊字元搜尋', () => {
      const criteria: SearchCriteria = {
        url: '/api/users?id=123&name=test'
      }
      
      // 這個測試確保特殊字元不會被當作正規表達式處理
      expect(() => service.search(testData, criteria)).not.toThrow()
    })

    it('應該處理空字串搜尋', () => {
      const criteria: SearchCriteria = {
        ip: '',
        url: '',
        userAgent: ''
      }
      
      const results = service.search(testData, criteria)
      
      // 空字串應該被忽略，返回所有資料
      expect(results).toHaveLength(testData.length)
    })
  })

  describe('效能要求', () => {
    it('應該快速搜尋大量資料（效能測試）', () => {
      // 生成 10,000 筆測試資料
      const largeDataset: LogEntry[] = []
      for (let i = 0; i < 10000; i++) {
        largeDataset.push(new models.LogEntry({
          ...testData[0],
          lineNumber: i + 1,
          ip: `192.168.${Math.floor(i / 256)}.${i % 256}`,
          url: `/page${i % 100}`
        }))
      }

      const startTime = performance.now()
      
      const criteria: SearchCriteria = {
        ip: '192.168.1.'
      }
      
      const results = service.search(largeDataset, criteria)
      
      const endTime = performance.now()
      const searchTime = endTime - startTime

      // 搜尋應該在 100ms 內完成
      expect(searchTime).toBeLessThan(100)
      expect(results.length).toBeGreaterThan(0)
    })
  })
})
