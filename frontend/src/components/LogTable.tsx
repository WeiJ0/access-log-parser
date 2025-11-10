// LogTable 虛擬化表格組件
// 使用 ag-Grid 顯示日誌記錄，支援虛擬化滾動
// 文件路徑: frontend/src/components/LogTable.tsx

import { useMemo } from 'react'
import { AgGridReact } from 'ag-grid-react'
import { ColDef } from 'ag-grid-community'
import { Box } from '@mui/material'
import type { models } from '../../wailsjs/wailsjs/go/models'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-material.css'

// 類型別名便於使用
type LogEntry = models.LogEntry

interface LogTableProps {
  entries: models.LogEntry[]
  height?: string | number
}

/**
 * LogTable 組件
 * 使用 ag-Grid 顯示日誌記錄，支援虛擬化滾動以處理大量數據
 * 
 * @param entries - 日誌記錄陣列
 * @param height - 表格高度（預設：'calc(100vh - 250px)'）
 */
export default function LogTable({ entries, height = 'calc(100vh - 250px)' }: LogTableProps) {
  // 定義表格欄位
  const columnDefs = useMemo<ColDef<LogEntry>[]>(() => [
    {
      field: 'lineNumber',
      headerName: '行號',
      width: 100,
      sortable: true,
      filter: 'agNumberColumnFilter',
      pinned: 'left',
    },
    {
      field: 'ip',
      headerName: 'IP 位址',
      width: 150,
      sortable: true,
      filter: 'agTextColumnFilter',
    },
    {
      field: 'timestamp',
      headerName: '時間戳',
      width: 200,
      sortable: true,
      filter: 'agDateColumnFilter',
      valueFormatter: (params) => {
        if (!params.value) return ''
        const date = new Date(params.value)
        return date.toLocaleString('zh-TW', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
        })
      },
    },
    {
      field: 'method',
      headerName: 'HTTP 方法',
      width: 120,
      sortable: true,
      filter: 'agTextColumnFilter',
      cellStyle: (params) => {
        // 根據 HTTP 方法著色
        const colors: { [key: string]: string } = {
          GET: '#4caf50',
          POST: '#2196f3',
          PUT: '#ff9800',
          DELETE: '#f44336',
          PATCH: '#9c27b0',
        }
        const color = colors[params.value as string] || '#757575'
        return { color, fontWeight: 'bold' }
      },
    },
    {
      field: 'url',
      headerName: 'URL 路徑',
      width: 300,
      sortable: true,
      filter: 'agTextColumnFilter',
      flex: 1, // 自動調整寬度
    },
    {
      field: 'protocol',
      headerName: '協定',
      width: 120,
      sortable: true,
      filter: 'agTextColumnFilter',
    },
    {
      field: 'statusCode',
      headerName: '狀態碼',
      width: 120,
      sortable: true,
      filter: 'agNumberColumnFilter',
      cellStyle: (params) => {
        // 根據狀態碼著色
        const code = params.value as number
        const style: Record<string, string> = {}
        
        if (code >= 200 && code < 300) {
          style.color = '#4caf50'
          style.fontWeight = 'bold'
        } else if (code >= 300 && code < 400) {
          style.color = '#ff9800'
          style.fontWeight = 'bold'
        } else if (code >= 400 && code < 500) {
          style.color = '#ff5722'
          style.fontWeight = 'bold'
        } else if (code >= 500) {
          style.color = '#f44336'
          style.fontWeight = 'bold'
        }
        
        return style
      },
    },
    {
      field: 'responseBytes',
      headerName: '回應大小',
      width: 130,
      sortable: true,
      filter: 'agNumberColumnFilter',
      valueFormatter: (params) => {
        if (params.value === undefined || params.value === null) return ''
        const bytes = params.value as number
        
        // 格式化為 KB 或 MB
        if (bytes < 1024) {
          return `${bytes} B`
        } else if (bytes < 1024 * 1024) {
          return `${(bytes / 1024).toFixed(2)} KB`
        } else {
          return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
        }
      },
    },
    {
      field: 'referer',
      headerName: 'Referer',
      width: 200,
      sortable: true,
      filter: 'agTextColumnFilter',
    },
    {
      field: 'userAgent',
      headerName: 'User Agent',
      width: 300,
      sortable: true,
      filter: 'agTextColumnFilter',
    },
  ], [])

  // 預設欄位設定
  const defaultColDef = useMemo<ColDef>(() => ({
    resizable: true,
    sortable: true,
    filter: true,
  }), [])

  return (
    <Box
      className="ag-theme-material"
      sx={{
        width: '100%',
        height: height,
      }}
    >
      <AgGridReact<LogEntry>
        rowData={entries}
        columnDefs={columnDefs}
        defaultColDef={defaultColDef}
        rowSelection="multiple"
        animateRows={true}
        enableCellTextSelection={true}
        ensureDomOrder={true}
        // 啟用虛擬化以處理大量數據
        rowBuffer={10}
        suppressColumnVirtualisation={false}
        // 分頁設定（可選）
        pagination={false}
        // 效能優化
        suppressRowHoverHighlight={false}
        suppressCellFocus={false}
      />
    </Box>
  )
}
