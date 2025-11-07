/**
 * 搜尋框組件
 * 文件路徑: frontend/src/components/SearchBar.tsx
 * 目的: 提供即時搜尋功能，支援 IP、URL、User-Agent 等欄位搜尋
 */

import React, { useState, useCallback, useEffect } from 'react'
import {
  TextField,
  InputAdornment,
  IconButton,
  Box,
  Chip,
  Stack,
  Tooltip,
  Paper,
  FormControlLabel,
  Switch
} from '@mui/material'
import SearchIcon from '@mui/icons-material/Search'
import ClearIcon from '@mui/icons-material/Clear'
import FilterListIcon from '@mui/icons-material/FilterList'
import { SearchCriteria } from '../services/searchService'

/**
 * SearchBar 組件屬性介面
 */
export interface SearchBarProps {
  /** 搜尋回調函式 */
  onSearch: (criteria: SearchCriteria) => void
  
  /** 清除搜尋回調函式 */
  onClear: () => void
  
  /** 顯示進階篩選面板回調函式 */
  onToggleFilter?: () => void
  
  /** 搜尋結果統計 */
  searchStats?: {
    total: number
    results: number
    percentage: number
  }
  
  /** 是否顯示進階篩選按鈕 */
  showFilterButton?: boolean
  
  /** 防抖延遲時間（毫秒，預設 300ms） */
  debounceDelay?: number
  
  /** 是否禁用 */
  disabled?: boolean
}

/**
 * SearchBar 搜尋框組件
 * 
 * 功能:
 * - 即時搜尋（帶防抖）
 * - 支援多種搜尋欄位（IP、URL、User-Agent、關鍵字）
 * - 顯示搜尋結果統計
 * - 大小寫敏感切換
 * - 快速清除功能
 */
export const SearchBar: React.FC<SearchBarProps> = ({
  onSearch,
  onClear,
  onToggleFilter,
  searchStats,
  showFilterButton = true,
  debounceDelay = 300,
  disabled = false
}) => {
  // 搜尋欄位狀態
  const [searchText, setSearchText] = useState('')
  const [searchField, setSearchField] = useState<'keyword' | 'ip' | 'url' | 'userAgent'>('keyword')
  const [caseSensitive, setCaseSensitive] = useState(false)

  // 防抖計時器
  const [debounceTimer, setDebounceTimer] = useState<number | null>(null)

  /**
   * 執行搜尋（帶防抖）
   */
  const performSearch = useCallback((text: string, field: string, sensitive: boolean) => {
    if (!text || text.trim() === '') {
      onClear()
      return
    }

    // 建立搜尋條件
    const criteria: SearchCriteria = {
      caseSensitive: sensitive
    }

    // 根據選擇的欄位設定搜尋條件
    switch (field) {
      case 'ip':
        criteria.ip = text
        break
      case 'url':
        criteria.url = text
        break
      case 'userAgent':
        criteria.userAgent = text
        break
      case 'keyword':
      default:
        criteria.keyword = text
        break
    }

    onSearch(criteria)
  }, [onSearch, onClear])

  /**
   * 處理搜尋文字變更（帶防抖）
   */
  const handleSearchChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const text = event.target.value
    setSearchText(text)

    // 清除舊的計時器
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }

    // 設定新的計時器（防抖）
    const timer = setTimeout(() => {
      performSearch(text, searchField, caseSensitive)
    }, debounceDelay)

    setDebounceTimer(timer)
  }, [debounceTimer, searchField, caseSensitive, performSearch, debounceDelay])

  /**
   * 處理搜尋欄位變更
   */
  const handleFieldChange = useCallback((field: 'keyword' | 'ip' | 'url' | 'userAgent') => {
    setSearchField(field)
    
    // 如果有搜尋文字，立即執行搜尋
    if (searchText.trim() !== '') {
      performSearch(searchText, field, caseSensitive)
    }
  }, [searchText, caseSensitive, performSearch])

  /**
   * 處理大小寫敏感切換
   */
  const handleCaseSensitiveChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const sensitive = event.target.checked
    setCaseSensitive(sensitive)
    
    // 如果有搜尋文字，立即執行搜尋
    if (searchText.trim() !== '') {
      performSearch(searchText, searchField, sensitive)
    }
  }, [searchText, searchField, performSearch])

  /**
   * 處理清除搜尋
   */
  const handleClear = useCallback(() => {
    setSearchText('')
    onClear()
    
    // 清除防抖計時器
    if (debounceTimer) {
      clearTimeout(debounceTimer)
      setDebounceTimer(null)
    }
  }, [onClear, debounceTimer])

  /**
   * 處理 Enter 鍵立即搜尋
   */
  const handleKeyDown = useCallback((event: React.KeyboardEvent) => {
    if (event.key === 'Enter') {
      // 清除防抖計時器，立即執行搜尋
      if (debounceTimer) {
        clearTimeout(debounceTimer)
        setDebounceTimer(null)
      }
      
      performSearch(searchText, searchField, caseSensitive)
    }
  }, [debounceTimer, searchText, searchField, caseSensitive, performSearch])

  /**
   * 組件卸載時清理計時器
   */
  useEffect(() => {
    return () => {
      if (debounceTimer) {
        clearTimeout(debounceTimer)
      }
    }
  }, [debounceTimer])

  /**
   * 取得欄位標籤
   */
  const getFieldLabel = () => {
    switch (searchField) {
      case 'ip':
        return '搜尋 IP 地址'
      case 'url':
        return '搜尋 URL 路徑'
      case 'userAgent':
        return '搜尋 User-Agent'
      case 'keyword':
      default:
        return '搜尋所有欄位'
    }
  }

  return (
    <Paper elevation={0} sx={{ p: 2, mb: 2 }}>
      <Box>
        {/* 搜尋輸入框 */}
        <TextField
          fullWidth
          variant="outlined"
          size="small"
          value={searchText}
          onChange={handleSearchChange}
          onKeyDown={handleKeyDown}
          placeholder={getFieldLabel()}
          disabled={disabled}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon />
              </InputAdornment>
            ),
            endAdornment: (
              <InputAdornment position="end">
                {searchText && (
                  <IconButton
                    size="small"
                    onClick={handleClear}
                    edge="end"
                    disabled={disabled}
                  >
                    <ClearIcon />
                  </IconButton>
                )}
                {showFilterButton && onToggleFilter && (
                  <Tooltip title="進階篩選">
                    <IconButton
                      size="small"
                      onClick={onToggleFilter}
                      edge="end"
                      disabled={disabled}
                      sx={{ ml: 1 }}
                    >
                      <FilterListIcon />
                    </IconButton>
                  </Tooltip>
                )}
              </InputAdornment>
            )
          }}
        />

        {/* 搜尋欄位選擇和選項 */}
        <Stack
          direction="row"
          spacing={1}
          alignItems="center"
          sx={{ mt: 2 }}
          flexWrap="wrap"
        >
          {/* 搜尋欄位選擇 */}
          <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
            <Chip
              label="所有欄位"
              onClick={() => handleFieldChange('keyword')}
              color={searchField === 'keyword' ? 'primary' : 'default'}
              variant={searchField === 'keyword' ? 'filled' : 'outlined'}
              size="small"
              disabled={disabled}
            />
            <Chip
              label="IP 地址"
              onClick={() => handleFieldChange('ip')}
              color={searchField === 'ip' ? 'primary' : 'default'}
              variant={searchField === 'ip' ? 'filled' : 'outlined'}
              size="small"
              disabled={disabled}
            />
            <Chip
              label="URL 路徑"
              onClick={() => handleFieldChange('url')}
              color={searchField === 'url' ? 'primary' : 'default'}
              variant={searchField === 'url' ? 'filled' : 'outlined'}
              size="small"
              disabled={disabled}
            />
            <Chip
              label="User-Agent"
              onClick={() => handleFieldChange('userAgent')}
              color={searchField === 'userAgent' ? 'primary' : 'default'}
              variant={searchField === 'userAgent' ? 'filled' : 'outlined'}
              size="small"
              disabled={disabled}
            />
          </Box>

          {/* 大小寫敏感開關 */}
          <FormControlLabel
            control={
              <Switch
                checked={caseSensitive}
                onChange={handleCaseSensitiveChange}
                size="small"
                disabled={disabled}
              />
            }
            label="區分大小寫"
            sx={{ ml: 'auto' }}
          />
        </Stack>

        {/* 搜尋結果統計 */}
        {searchStats && searchStats.results < searchStats.total && (
          <Box sx={{ mt: 2, color: 'text.secondary', fontSize: '0.875rem' }}>
            顯示 {searchStats.results.toLocaleString()} 筆，共 {searchStats.total.toLocaleString()} 筆
            （{searchStats.percentage.toFixed(1)}%）
          </Box>
        )}
      </Box>
    </Paper>
  )
}

export default SearchBar
