// 主應用程式組件
// 文件路徑: frontend/src/App.tsx

import { useState } from 'react'
import {
  Container,
  AppBar,
  Toolbar,
  Typography,
  Box,
  Tabs,
  Tab,
  Button,
  Alert,
  CircularProgress,
} from '@mui/material'
import FolderOpenIcon from '@mui/icons-material/FolderOpen'
import { SelectFile, ParseFile } from '../wailsjs/go/main/App'
// Note: 這些函式目前還未在 Go 後端實作完成，先使用佔位符
// TODO: 在後端實作完成後，這裡會自動綁定到 Go 函式

import type { LogEntry } from './types/log'
import TabPanel from './components/TabPanel'
import LogTable from './components/LogTable'
import ErrorSummary from './components/ErrorSummary'
import ProgressIndicator from './components/ProgressIndicator'

// 檔案資訊介面
interface FileInfo {
  path: string
  name: string
  entries: LogEntry[]
  errorCount: number
  errorSamples: string[]
  isLoading: boolean
  error: string | null
}

function App() {
  const [files, setFiles] = useState<FileInfo[]>([])
  const [currentTab, setCurrentTab] = useState(0)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // 處理檔案選擇
  const handleSelectFile = async () => {
    try {
      setError(null)
      setLoading(true)

      // 呼叫 Wails API 選擇檔案
      const filePath = await SelectFile()
      
      if (!filePath) {
        setLoading(false)
        return
      }

      // 檢查檔案是否已經開啟
      const existingIndex = files.findIndex(f => f.path === filePath)
      if (existingIndex >= 0) {
        setCurrentTab(existingIndex)
        setLoading(false)
        return
      }

      // 解析檔案
      const fileName = filePath.split(/[/\\]/).pop() || filePath
      
      // 新增檔案到列表（載入中狀態）
      const newFile: FileInfo = {
        path: filePath,
        name: fileName,
        entries: [],
        errorCount: 0,
        errorSamples: [],
        isLoading: true,
        error: null,
      }
      
      const newFiles = [...files, newFile]
      setFiles(newFiles)
      setCurrentTab(newFiles.length - 1)

      // 呼叫 Wails API 解析檔案
      const result = await ParseFile(filePath)
      
      // 更新檔案資訊
      const updatedFiles = newFiles.map((f, i) => {
        if (i === newFiles.length - 1) {
          return {
            ...f,
            entries: result.entries || [],
            errorCount: result.errorCount || 0,
            errorSamples: result.errorSamples || [],
            isLoading: false,
            error: null,
          }
        }
        return f
      })
      
      setFiles(updatedFiles)
      setLoading(false)

    } catch (err) {
      console.error('載入檔案失敗:', err)
      setError(err instanceof Error ? err.message : '未知錯誤')
      setLoading(false)
      
      // 移除載入失敗的檔案
      setFiles(files.filter(f => !f.isLoading))
    }
  }

  // 處理分頁切換
  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue)
  }

  // 關閉分頁
  // TODO: 實作分頁關閉功能
  // const handleCloseTab = (index: number) => {
  //   const newFiles = files.filter((_, i) => i !== index)
  //   setFiles(newFiles)
  //   
  //   // 調整當前分頁
  //   if (currentTab >= newFiles.length) {
  //     setCurrentTab(Math.max(0, newFiles.length - 1))
  //   }
  // }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>
      {/* 頂部工具列 */}
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Apache Access Log Analyzer
          </Typography>
          <Button
            color="inherit"
            startIcon={<FolderOpenIcon />}
            onClick={handleSelectFile}
            disabled={loading}
          >
            開啟檔案
          </Button>
        </Toolbar>
      </AppBar>

      {/* 錯誤提示 */}
      {error && (
        <Alert severity="error" onClose={() => setError(null)} sx={{ m: 2 }}>
          {error}
        </Alert>
      )}

      {/* 分頁列 */}
      {files.length > 0 && (
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={currentTab} onChange={handleTabChange}>
            {files.map((file) => (
              <Tab
                key={file.path}
                label={
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    {file.name}
                    {file.isLoading && <CircularProgress size={16} />}
                  </Box>
                }
              />
            ))}
          </Tabs>
        </Box>
      )}

      {/* 內容區域 */}
      <Box sx={{ flexGrow: 1, overflow: 'auto' }}>
        {files.length === 0 ? (
          // 歡迎畫面
          <Container maxWidth="md" sx={{ mt: 8, textAlign: 'center' }}>
            <Typography variant="h4" gutterBottom>
              歡迎使用 Apache Access Log Analyzer
            </Typography>
            <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
              點擊上方的「開啟檔案」按鈕來載入 Apache access log 檔案
            </Typography>
            <Button
              variant="contained"
              size="large"
              startIcon={<FolderOpenIcon />}
              onClick={handleSelectFile}
              disabled={loading}
            >
              選擇檔案
            </Button>
          </Container>
        ) : (
          // 顯示檔案內容
          files.map((file, index) => (
            <TabPanel key={file.path} value={currentTab} index={index}>
              {file.isLoading ? (
                <ProgressIndicator
                  loading={true}
                  message={`正在解析 ${file.name}...`}
                />
              ) : file.error ? (
                <Alert severity="error">{file.error}</Alert>
              ) : (
                <Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="h6" gutterBottom>
                      {file.name}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      總記錄數: {file.entries.length.toLocaleString()} | 錯誤數: {file.errorCount.toLocaleString()}
                    </Typography>
                  </Box>
                  
                  {/* 錯誤摘要 */}
                  {file.errorCount > 0 && (
                    <ErrorSummary
                      errorCount={file.errorCount}
                      errorSamples={file.errorSamples}
                    />
                  )}
                  
                  {/* 日誌表格 */}
                  <LogTable entries={file.entries} />
                </Box>
              )}
            </TabPanel>
          ))
        )}
      </Box>
    </Box>
  )
}

export default App
