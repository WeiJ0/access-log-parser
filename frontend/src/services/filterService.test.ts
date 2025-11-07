// @ts-nocheck
/**
 * 篩選器服務測試
 * 文件路徑: frontend/src/services/filterService.test.ts
 * 目的: 測試各種篩選條件（狀態碼、時間範圍、HTTP 方法）
 */

import { describe, it, expect, beforeEach } from 'vitest'
import { FilterService, FilterCriteria } from './filterService'
import { models } from '../../wailsjs/wailsjs/go/models'

type LogEntry = models.LogEntry

describe('FilterService', () => {
  let service: FilterService
  let testData: LogEntry[]

  beforeEach(() => {
    service = new FilterService()
    
    // 準備測試資料（涵蓋不同的狀態碼和時間）
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
        userAgent: 'Mozilla/5.0',
        rawLine: ''
      }),
      new models.LogEntry({
        lineNumber: 2,
        ip: '10.0.0.50',
        user: '-',
        timestamp: '2024-01-01T11:00:00Z',
        method: 'POST',
        url: '/api/users',
        protocol: 'HTTP/1.1',
        statusCode: 201,
        responseBytes: 512,
        referer: '-',
        userAgent: 'curl/7.68.0',
        rawLine: ''
      }),
      new models.LogEntry({
        lineNumber: 3,
        ip: '192.168.1.100',
        user: '-',
        timestamp: '2024-01-01T12:00:00Z',
        method: 'GET',
        url: '/admin',
        protocol: 'HTTP/1.1',
        statusCode: 403,
        responseBytes: 256,
        referer: '-',
        userAgent: 'Mozilla/5.0',
        rawLine: ''
      }),
      new models.LogEntry({
        lineNumber: 4,
        ip: '172.16.0.1',
        user: '-',
        timestamp: '2024-01-01T13:00:00Z',
        method: 'GET',
        url: '/missing',
        protocol: 'HTTP/1.1',
        statusCode: 404,
        responseBytes: 128,
        referer: '-',
        userAgent: 'Googlebot',
        rawLine: ''
      }),
      new models.LogEntry({
        lineNumber: 5,
        ip: '10.0.0.1',
        user: '-',
        timestamp: '2024-01-01T14:00:00Z',
        method: 'POST',
        url: '/api/error',
        protocol: 'HTTP/1.1',
        statusCode: 500,
        responseBytes: 64,
        referer: '-',
        userAgent: 'python-requests',
        rawLine: ''
      }),
      new models.LogEntry({
        lineNumber: 6,
        ip: '192.168.1.1',
        user: '-',
        timestamp: '2024-01-02T10:00:00Z',
        method: 'DELETE',
        url: '/api/resource',
        protocol: 'HTTP/1.1',
        statusCode: 204,
        responseBytes: 0,
        referer: '-',
        userAgent: 'curl/7.68.0',
        rawLine: ''
      })
    ]
  })

  describe('狀態碼篩選', () => {
    it('應該篩選單一狀態碼', () => {
      const criteria: FilterCriteria = {
        statusCodes: [200]
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(1)
      expect(results[0].statusCode).toBe(200)
    })

    it('應該篩選多個狀態碼', () => {
      const criteria: FilterCriteria = {
        statusCodes: [200, 201, 204]
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(3)
      expect(results.every((r: LogEntry) => [200, 201, 204].includes(r.statusCode))).toBe(true)
    })

    it('應該篩選狀態碼範圍（2xx 成功）', () => {
      const criteria: FilterCriteria = {
        statusCodeRange: { min: 200, max: 299 }
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(3)
      expect(results.every((r: LogEntry) => r.statusCode >= 200 && r.statusCode < 300)).toBe(true)
    })

    it('應該篩選狀態碼範圍（4xx 客戶端錯誤）', () => {
      const criteria: FilterCriteria = {
        statusCodeRange: { min: 400, max: 499 }
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(2)
      expect(results.every((r: LogEntry) => r.statusCode >= 400 && r.statusCode < 500)).toBe(true)
    })

    it('應該篩選狀態碼範圍（5xx 伺服器錯誤）', () => {
      const criteria: FilterCriteria = {
        statusCodeRange: { min: 500, max: 599 }
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(1)
      expect(results[0].statusCode).toBe(500)
    })
  })

  describe('時間範圍篩選', () => {
    it('應該篩選開始時間之後的記錄', () => {
      const criteria: FilterCriteria = {
        timeRange: {
          start: '2024-01-01T12:00:00Z'
        }
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(4)
      expect(results.every((r: LogEntry) => r.timestamp >= '2024-01-01T12:00:00Z')).toBe(true)
    })

    it('應該篩選結束時間之前的記錄', () => {
      const criteria: FilterCriteria = {
        timeRange: {
          end: '2024-01-01T12:00:00Z'
        }
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(3)
      expect(results.every((r: LogEntry) => r.timestamp <= '2024-01-01T12:00:00Z')).toBe(true)
    })

    it('應該篩選時間範圍（含開始和結束）', () => {
      const criteria: FilterCriteria = {
        timeRange: {
          start: '2024-01-01T11:00:00Z',
          end: '2024-01-01T13:00:00Z'
        }
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(3)
      expect(results.every((r: LogEntry) => 
        r.timestamp >= '2024-01-01T11:00:00Z' && 
        r.timestamp <= '2024-01-01T13:00:00Z'
      )).toBe(true)
    })
  })

  describe('HTTP 方法篩選', () => {
    it('應該篩選單一 HTTP 方法', () => {
      const criteria: FilterCriteria = {
        methods: ['GET']
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(3)
      expect(results.every((r: LogEntry) => r.method === 'GET')).toBe(true)
    })

    it('應該篩選多個 HTTP 方法', () => {
      const criteria: FilterCriteria = {
        methods: ['POST', 'DELETE']
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(3)
      expect(results.every((r: LogEntry) => ['POST', 'DELETE'].includes(r.method))).toBe(true)
    })
  })

  describe('回應大小篩選', () => {
    it('應該篩選最小回應大小', () => {
      const criteria: FilterCriteria = {
        responseSizeRange: {
          min: 256
        }
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(3)
      expect(results.every((r: LogEntry) => r.responseBytes >= 256)).toBe(true)
    })

    it('應該篩選最大回應大小', () => {
      const criteria: FilterCriteria = {
        responseSizeRange: {
          max: 256
        }
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(4)
      expect(results.every((r: LogEntry) => r.responseBytes <= 256)).toBe(true)
    })

    it('應該篩選回應大小範圍', () => {
      const criteria: FilterCriteria = {
        responseSizeRange: {
          min: 100,
          max: 600
        }
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(3)
      expect(results.every((r: LogEntry) => 
        r.responseBytes >= 100 && r.responseBytes <= 600
      )).toBe(true)
    })
  })

  describe('複合篩選', () => {
    it('應該組合狀態碼和時間篩選', () => {
      const criteria: FilterCriteria = {
        statusCodeRange: { min: 200, max: 299 },
        timeRange: {
          start: '2024-01-01T11:00:00Z',
          end: '2024-01-01T13:00:00Z'
        }
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(1)
      expect(results[0].statusCode).toBe(201)
      expect(results[0].timestamp).toBe('2024-01-01T11:00:00Z')
    })

    it('應該組合多個篩選條件（AND 邏輯）', () => {
      const criteria: FilterCriteria = {
        statusCodeRange: { min: 200, max: 299 },
        methods: ['POST'],
        responseSizeRange: { min: 400 }
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(1)
      expect(results.every((r: LogEntry) => 
        r.statusCode >= 200 && r.statusCode < 300 &&
        r.method === 'POST' &&
        r.responseBytes >= 400
      )).toBe(true)
    })

    it('應該返回沒有符合記錄的空陣列', () => {
      const criteria: FilterCriteria = {
        statusCodes: [999]
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(0)
    })
  })

  describe('邊界條件', () => {
    it('應該處理空的資料集', () => {
      const criteria: FilterCriteria = {
        statusCodes: [200]
      }
      
      const results = service.filter([], criteria)
      
      expect(results).toHaveLength(0)
    })

    it('應該處理空的篩選條件（返回於結果）', () => {
      const criteria: FilterCriteria = {
        timeRange: {
          start: '2024-01-01T00:00:00Z',
          end: '2024-01-03T00:00:00Z'
        }
      }
      
      const results = service.filter(testData, criteria)
      
      expect(results).toHaveLength(6)
    })

    it('應該檢查篩選條件是否為空', () => {
      const emptyCriteria: FilterCriteria = {}
      const emptyResult = service.filter(testData, emptyCriteria)
      
      // 空條件應該返回所有記錄
      expect(emptyResult).toHaveLength(testData.length)
      
      const nonEmptyCriteria: FilterCriteria = {
        statusCodes: [200]
      }
      const nonEmptyResult = service.filter(testData, nonEmptyCriteria)
      
      // 非空條件應該篩選記錄
      expect(nonEmptyResult.length).toBeLessThan(testData.length)
    })
  })

  describe('效能要求', () => {
    it('應該快速篩選大量資料（效能測試）', () => {
      // 生成 100,000 筆測試資料
      const largeDataset: LogEntry[] = []
      for (let i = 0; i < 100000; i++) {
        largeDataset.push(new models.LogEntry({
          ...testData[0],
          lineNumber: i + 1,
          statusCode: 200 + (i % 400), // 200-599
          timestamp: new Date(2024, 0, 1 + (i % 30), i % 24).toISOString()
        }))
      }

      const startTime = performance.now()
      
      const criteria: FilterCriteria = {
        statusCodeRange: { min: 400, max: 499 },
        timeRange: {
          start: '2024-01-10T00:00:00Z',
          end: '2024-01-20T23:59:59Z'
        }
      }
      
      const results = service.filter(largeDataset, criteria)
      
      const endTime = performance.now()
      const filterTime = endTime - startTime

      // 篩選應該在 100ms 內完成
      expect(filterTime).toBeLessThan(100)
      
      // 驗證結果正確性
      expect(results.length).toBeGreaterThan(0)
    })
  })

  describe('統計資訊', () => {
    it('應該返回篩選結果統計資訊', () => {
      const criteria: FilterCriteria = {
        statusCodeRange: { min: 400, max: 499 }
      }
      
      const results = service.filter(testData, criteria)
      const stats = service.getFilterStats(testData, results)
      
      expect(stats.total).toBe(6)
      expect(stats.filtered).toBe(2)
      expect(stats.percentage).toBeCloseTo(33.33, 1)
    })
  })
})
