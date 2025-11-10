// 最近開啟的檔案列表組件
// 文件路徑: frontend/src/components/RecentFiles.tsx

import { useState } from 'react'
import {
  Menu,
  MenuItem,
  ListItemText,
  Typography,
  Divider,
  Box,
  IconButton,
  Tooltip,
} from '@mui/material'
import HistoryIcon from '@mui/icons-material/History'
import ClearIcon from '@mui/icons-material/Clear'
import * as AppAPI from '../../wailsjs/wailsjs/go/app/App'

// 最近檔案資訊介面（對應後端 RecentFile）
interface RecentFile {
  path: string
  name: string
  size: number
  openedAt: string  // ISO 8601 格式
  totalLines: number
}

interface RecentFilesProps {
  onFileSelect: (path: string) => void  // 選擇檔案的回調
}

/**
 * RecentFiles 組件
 * 顯示最近開啟的檔案列表，支援快速重新開啟
 */
export default function RecentFiles({ onFileSelect }: RecentFilesProps) {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [recentFiles, setRecentFiles] = useState<RecentFile[]>([])
  const [loading, setLoading] = useState(false)

  const open = Boolean(anchorEl)

  // 載入最近檔案列表
  const loadRecentFiles = async () => {
    try {
      setLoading(true)
      const response = await AppAPI.GetRecentFiles()
      if (response.success) {
        setRecentFiles(response.files || [])
      }
    } catch (error) {
      console.error('載入最近檔案失敗:', error)
    } finally {
      setLoading(false)
    }
  }

  // 開啟選單時載入列表
  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
    loadRecentFiles()
  }

  const handleClose = () => {
    setAnchorEl(null)
  }

  // 選擇檔案
  const handleFileSelect = (path: string) => {
    handleClose()
    onFileSelect(path)
  }

  // 清空列表
  const handleClear = async (event: React.MouseEvent) => {
    event.stopPropagation()  // 防止關閉選單
    try {
      const response = await AppAPI.ClearRecentFiles()
      if (response.success) {
        setRecentFiles([])
      }
    } catch (error) {
      console.error('清空最近檔案失敗:', error)
    }
  }

  // 格式化檔案大小
  const formatSize = (bytes: number): string => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  // 格式化日期時間
  const formatDateTime = (isoString: string): string => {
    try {
      const date = new Date(isoString)
      return date.toLocaleString('zh-TW', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      })
    } catch {
      return isoString
    }
  }

  return (
    <>
      <Tooltip title="最近開啟的檔案">
        <IconButton
          color="inherit"
          onClick={handleClick}
          aria-controls={open ? 'recent-files-menu' : undefined}
          aria-haspopup="true"
          aria-expanded={open ? 'true' : undefined}
        >
          <HistoryIcon />
        </IconButton>
      </Tooltip>

      <Menu
        id="recent-files-menu"
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
        PaperProps={{
          sx: {
            maxWidth: 500,
            maxHeight: 400,
          },
        }}
      >
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            px: 2,
            py: 1,
          }}
        >
          <Typography variant="subtitle2" color="text.secondary">
            最近開啟的檔案
          </Typography>
          {recentFiles.length > 0 && (
            <Tooltip title="清空列表">
              <IconButton size="small" onClick={handleClear}>
                <ClearIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          )}
        </Box>
        <Divider />

        {loading ? (
          <MenuItem disabled>
            <ListItemText primary="載入中..." />
          </MenuItem>
        ) : recentFiles.length === 0 ? (
          <MenuItem disabled>
            <ListItemText 
              primary="無最近開啟的檔案"
              secondary="開啟檔案後將顯示在此列表"
            />
          </MenuItem>
        ) : (
          recentFiles.map((file, index) => (
            <MenuItem
              key={`${file.path}-${index}`}
              onClick={() => handleFileSelect(file.path)}
            >
              <ListItemText
                primary={file.name}
                secondary={
                  <Box component="span">
                    <Typography variant="caption" display="block">
                      {formatSize(file.size)} · {file.totalLines.toLocaleString()} 行
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {formatDateTime(file.openedAt)}
                    </Typography>
                  </Box>
                }
                sx={{
                  '& .MuiListItemText-primary': {
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                  },
                }}
              />
            </MenuItem>
          ))
        )}
      </Menu>
    </>
  )
}
