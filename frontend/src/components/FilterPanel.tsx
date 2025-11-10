/**
 * 篩選器面板組件
 * 文件路徑: frontend/src/components/FilterPanel.tsx
 * 目的: 提供進階篩選功能，支援狀態碼、時間範圍、HTTP 方法等篩選
 */

import React, { useState, useCallback } from 'react'
import {
  Box,
  Paper,
  Typography,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  FormGroup,
  FormControlLabel,
  Checkbox,
  TextField,
  Button,
  Stack,
  Chip,
  Divider,
  Grid
} from '@mui/material'
import ExpandMoreIcon from '@mui/icons-material/ExpandMore'
import FilterAltOffIcon from '@mui/icons-material/FilterAltOff'
import CheckIcon from '@mui/icons-material/Check'
import { FilterCriteria, FilterService } from '../services/filterService'

/**
 * FilterPanel 組件屬性介面
 */
export interface FilterPanelProps {
  /** 套用篩選回調函式 */
  onApplyFilter: (criteria: FilterCriteria) => void
  
  /** 清除篩選回調函式 */
  onClearFilter: () => void
  
  /** 篩選結果統計 */
  filterStats?: {
    total: number
    filtered: number
    percentage: number
  }
  
  /** 是否禁用 */
  disabled?: boolean
}

/**
 * FilterPanel 篩選器面板組件
 * 
 * 功能:
 * - 狀態碼篩選（單選或範圍）
 * - 時間範圍篩選
 * - HTTP 方法篩選
 * - 回應大小篩選
 * - 預定義的快捷篩選
 * - 清除所有篩選
 */
export const FilterPanel: React.FC<FilterPanelProps> = ({
  onApplyFilter,
  onClearFilter,
  filterStats,
  disabled = false
}) => {
  // 狀態碼篩選
  const [selectedStatusCodes, setSelectedStatusCodes] = useState<number[]>([])
  const [statusCodeRangeMode, setStatusCodeRangeMode] = useState(false)
  const [statusCodeMin, setStatusCodeMin] = useState('')
  const [statusCodeMax, setStatusCodeMax] = useState('')

  // 時間範圍篩選
  const [timeRangeStart, setTimeRangeStart] = useState('')
  const [timeRangeEnd, setTimeRangeEnd] = useState('')

  // HTTP 方法篩選
  const [selectedMethods, setSelectedMethods] = useState<string[]>([])

  // 回應大小篩選
  const [responseSizeMin, setResponseSizeMin] = useState('')
  const [responseSizeMax, setResponseSizeMax] = useState('')

  /**
   * 常用狀態碼列表
   */
  const commonStatusCodes = [
    { code: 200, label: '200 成功' },
    { code: 201, label: '201 已建立' },
    { code: 204, label: '204 無內容' },
    { code: 301, label: '301 永久重新導向' },
    { code: 302, label: '302 暫時重新導向' },
    { code: 304, label: '304 未修改' },
    { code: 400, label: '400 錯誤請求' },
    { code: 401, label: '401 未授權' },
    { code: 403, label: '403 禁止' },
    { code: 404, label: '404 找不到' },
    { code: 500, label: '500 伺服器錯誤' },
    { code: 502, label: '502 閘道錯誤' },
    { code: 503, label: '503 服務無法使用' }
  ]

  /**
   * 常用 HTTP 方法列表
   */
  const httpMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS']

  /**
   * 處理狀態碼選擇
   */
  const handleStatusCodeToggle = useCallback((code: number) => {
    setSelectedStatusCodes(prev => {
      if (prev.includes(code)) {
        return prev.filter(c => c !== code)
      } else {
        return [...prev, code]
      }
    })
  }, [])

  /**
   * 處理 HTTP 方法選擇
   */
  const handleMethodToggle = useCallback((method: string) => {
    setSelectedMethods(prev => {
      if (prev.includes(method)) {
        return prev.filter(m => m !== method)
      } else {
        return [...prev, method]
      }
    })
  }, [])

  /**
   * 套用快捷狀態碼範圍
   */
  const applyStatusCodeRange = useCallback((min: number, max: number) => {
    setStatusCodeRangeMode(true)
    setStatusCodeMin(min.toString())
    setStatusCodeMax(max.toString())
    setSelectedStatusCodes([])
  }, [])

  /**
   * 套用快捷時間範圍
   */
  const applyTimeRange = useCallback((start: string, end: string) => {
    setTimeRangeStart(start)
    setTimeRangeEnd(end)
  }, [])

  /**
   * 建立並套用篩選條件
   */
  const handleApplyFilter = useCallback(() => {
    const criteria: FilterCriteria = {}

    // 狀態碼篩選
    if (statusCodeRangeMode) {
      const min = parseInt(statusCodeMin)
      const max = parseInt(statusCodeMax)
      if (!isNaN(min) && !isNaN(max)) {
        criteria.statusCodeRange = { min, max }
      }
    } else if (selectedStatusCodes.length > 0) {
      criteria.statusCodes = selectedStatusCodes
    }

    // 時間範圍篩選
    if (timeRangeStart || timeRangeEnd) {
      criteria.timeRange = {}
      if (timeRangeStart) {
        criteria.timeRange.start = new Date(timeRangeStart).toISOString()
      }
      if (timeRangeEnd) {
        criteria.timeRange.end = new Date(timeRangeEnd).toISOString()
      }
    }

    // HTTP 方法篩選
    if (selectedMethods.length > 0) {
      criteria.methods = selectedMethods
    }

    // 回應大小篩選
    if (responseSizeMin || responseSizeMax) {
      criteria.responseSizeRange = {}
      if (responseSizeMin) {
        criteria.responseSizeRange.min = parseInt(responseSizeMin)
      }
      if (responseSizeMax) {
        criteria.responseSizeRange.max = parseInt(responseSizeMax)
      }
    }

    // 驗證並套用篩選
    const filterService = new FilterService()
    const validation = filterService.validateCriteria(criteria)
    
    if (!validation.valid) {
      alert('篩選條件無效:\n' + validation.errors.join('\n'))
      return
    }

    onApplyFilter(criteria)
  }, [
    statusCodeRangeMode,
    statusCodeMin,
    statusCodeMax,
    selectedStatusCodes,
    timeRangeStart,
    timeRangeEnd,
    selectedMethods,
    responseSizeMin,
    responseSizeMax,
    onApplyFilter
  ])

  /**
   * 清除所有篩選
   */
  const handleClearFilter = useCallback(() => {
    setSelectedStatusCodes([])
    setStatusCodeRangeMode(false)
    setStatusCodeMin('')
    setStatusCodeMax('')
    setTimeRangeStart('')
    setTimeRangeEnd('')
    setSelectedMethods([])
    setResponseSizeMin('')
    setResponseSizeMax('')
    onClearFilter()
  }, [onClearFilter])

  /**
   * 快捷時間範圍按鈕
   */
  const timeRanges = FilterService.createTimeRanges()

  return (
    <Paper elevation={1} sx={{ p: 2, mb: 2 }}>
      <Box sx={{ mb: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h6">進階篩選</Typography>
        <Button
          startIcon={<FilterAltOffIcon />}
          onClick={handleClearFilter}
          size="small"
          disabled={disabled}
        >
          清除所有篩選
        </Button>
      </Box>

      {/* 狀態碼篩選 */}
      <Accordion defaultExpanded>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography>HTTP 狀態碼</Typography>
          {selectedStatusCodes.length > 0 && (
            <Chip
              label={`${selectedStatusCodes.length} 個`}
              size="small"
              color="primary"
              sx={{ ml: 2 }}
            />
          )}
        </AccordionSummary>
        <AccordionDetails>
          {/* 狀態碼範圍快捷按鈕 */}
          <Stack direction="row" spacing={1} flexWrap="wrap" sx={{ mb: 2 }}>
            <Chip
              label="2xx 成功"
              onClick={() => applyStatusCodeRange(200, 299)}
              size="small"
              color={statusCodeRangeMode && statusCodeMin === '200' && statusCodeMax === '299' ? 'success' : 'default'}
              disabled={disabled}
            />
            <Chip
              label="3xx 重新導向"
              onClick={() => applyStatusCodeRange(300, 399)}
              size="small"
              color={statusCodeRangeMode && statusCodeMin === '300' && statusCodeMax === '399' ? 'info' : 'default'}
              disabled={disabled}
            />
            <Chip
              label="4xx 客戶端錯誤"
              onClick={() => applyStatusCodeRange(400, 499)}
              size="small"
              color={statusCodeRangeMode && statusCodeMin === '400' && statusCodeMax === '499' ? 'warning' : 'default'}
              disabled={disabled}
            />
            <Chip
              label="5xx 伺服器錯誤"
              onClick={() => applyStatusCodeRange(500, 599)}
              size="small"
              color={statusCodeRangeMode && statusCodeMin === '500' && statusCodeMax === '599' ? 'error' : 'default'}
              disabled={disabled}
            />
          </Stack>

          <Divider sx={{ my: 2 }} />

          {/* 自訂範圍 */}
          <FormControlLabel
            control={
              <Checkbox
                checked={statusCodeRangeMode}
                onChange={(e) => {
                  setStatusCodeRangeMode(e.target.checked)
                  if (e.target.checked) {
                    setSelectedStatusCodes([])
                  }
                }}
                disabled={disabled}
              />
            }
            label="自訂範圍"
          />
          
          {statusCodeRangeMode && (
            <Box sx={{ display: 'flex', gap: 2, mt: 1 }}>
              <TextField
                label="最小值"
                type="number"
                size="small"
                value={statusCodeMin}
                onChange={(e) => setStatusCodeMin(e.target.value)}
                disabled={disabled}
                sx={{ flex: 1 }}
              />
              <TextField
                label="最大值"
                type="number"
                size="small"
                value={statusCodeMax}
                onChange={(e) => setStatusCodeMax(e.target.value)}
                disabled={disabled}
                sx={{ flex: 1 }}
              />
            </Box>
          )}

          {!statusCodeRangeMode && (
            <>
              <Divider sx={{ my: 2 }} />
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                選擇特定狀態碼:
              </Typography>
              <Grid container spacing={1}>
                {commonStatusCodes.map(({ code, label }) => (
                  <Grid item xs={6} key={code}>
                    <FormControlLabel
                      control={
                        <Checkbox
                          checked={selectedStatusCodes.includes(code)}
                          onChange={() => handleStatusCodeToggle(code)}
                          disabled={disabled}
                        />
                      }
                      label={label}
                    />
                  </Grid>
                ))}
              </Grid>
            </>
          )}
        </AccordionDetails>
      </Accordion>

      {/* 時間範圍篩選 */}
      <Accordion>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography>時間範圍</Typography>
        </AccordionSummary>
        <AccordionDetails>
          {/* 快捷時間範圍 */}
          <Stack direction="row" spacing={1} flexWrap="wrap" sx={{ mb: 2 }}>
            <Chip
              label="今天"
              onClick={() => applyTimeRange(timeRanges.TODAY.start!, timeRanges.TODAY.end!)}
              size="small"
              disabled={disabled}
            />
            <Chip
              label="昨天"
              onClick={() => applyTimeRange(timeRanges.YESTERDAY.start!, timeRanges.YESTERDAY.end!)}
              size="small"
              disabled={disabled}
            />
            <Chip
              label="過去 7 天"
              onClick={() => applyTimeRange(timeRanges.LAST_7_DAYS.start!, timeRanges.LAST_7_DAYS.end!)}
              size="small"
              disabled={disabled}
            />
            <Chip
              label="過去 30 天"
              onClick={() => applyTimeRange(timeRanges.LAST_30_DAYS.start!, timeRanges.LAST_30_DAYS.end!)}
              size="small"
              disabled={disabled}
            />
          </Stack>

          <Divider sx={{ my: 2 }} />

          {/* 自訂時間範圍 */}
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="開始時間"
              type="datetime-local"
              size="small"
              value={timeRangeStart}
              onChange={(e) => setTimeRangeStart(e.target.value)}
              disabled={disabled}
              InputLabelProps={{ shrink: true }}
            />
            <TextField
              label="結束時間"
              type="datetime-local"
              size="small"
              value={timeRangeEnd}
              onChange={(e) => setTimeRangeEnd(e.target.value)}
              disabled={disabled}
              InputLabelProps={{ shrink: true }}
            />
          </Box>
        </AccordionDetails>
      </Accordion>

      {/* HTTP 方法篩選 */}
      <Accordion>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography>HTTP 方法</Typography>
          {selectedMethods.length > 0 && (
            <Chip
              label={`${selectedMethods.length} 個`}
              size="small"
              color="primary"
              sx={{ ml: 2 }}
            />
          )}
        </AccordionSummary>
        <AccordionDetails>
          <FormGroup>
            <Grid container spacing={1}>
              {httpMethods.map(method => (
                <Grid item xs={6} key={method}>
                  <FormControlLabel
                    control={
                      <Checkbox
                        checked={selectedMethods.includes(method)}
                        onChange={() => handleMethodToggle(method)}
                        disabled={disabled}
                      />
                    }
                    label={method}
                  />
                </Grid>
              ))}
            </Grid>
          </FormGroup>
        </AccordionDetails>
      </Accordion>

      {/* 回應大小篩選 */}
      <Accordion>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography>回應大小</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <Box sx={{ display: 'flex', gap: 2 }}>
            <TextField
              label="最小值 (bytes)"
              type="number"
              size="small"
              value={responseSizeMin}
              onChange={(e) => setResponseSizeMin(e.target.value)}
              disabled={disabled}
              sx={{ flex: 1 }}
            />
            <TextField
              label="最大值 (bytes)"
              type="number"
              size="small"
              value={responseSizeMax}
              onChange={(e) => setResponseSizeMax(e.target.value)}
              disabled={disabled}
              sx={{ flex: 1 }}
            />
          </Box>
        </AccordionDetails>
      </Accordion>

      {/* 套用按鈕和統計 */}
      <Box sx={{ mt: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Button
          variant="contained"
          startIcon={<CheckIcon />}
          onClick={handleApplyFilter}
          disabled={disabled}
        >
          套用篩選
        </Button>

        {filterStats && filterStats.filtered < filterStats.total && (
          <Typography variant="body2" color="text.secondary">
            顯示 {filterStats.filtered.toLocaleString()} / {filterStats.total.toLocaleString()} 筆
            （{filterStats.percentage.toFixed(1)}%）
          </Typography>
        )}
      </Box>
    </Paper>
  )
}

export default FilterPanel
