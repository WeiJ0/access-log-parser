// TabPanel 分頁管理組件
// 支援多檔案分頁顯示
// 文件路徑: frontend/src/components/TabPanel.tsx

import React from 'react'
import { Box } from '@mui/material'

interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
  id?: string
}

/**
 * TabPanel 組件
 * 用於顯示分頁內容，僅在當前分頁激活時顯示子組件
 * 
 * @param children - 分頁內容
 * @param index - 分頁索引
 * @param value - 當前激活的分頁索引
 * @param id - 可選的 ID，用於 ARIA 標籤
 */
export default function TabPanel(props: TabPanelProps) {
  const { children, value, index, id, ...other } = props

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={id || `tabpanel-${index}`}
      aria-labelledby={`tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3, height: '100%', overflow: 'auto' }}>
          {children}
        </Box>
      )}
    </div>
  )
}
