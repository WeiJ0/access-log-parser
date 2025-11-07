// 主應用程式組件
// 文件路徑: frontend/src/App.tsx

import { useState, useEffect } from 'react'
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
  Snackbar,
  IconButton,
} from '@mui/material'
import FolderOpenIcon from '@mui/icons-material/FolderOpen'
import FileDownloadIcon from '@mui/icons-material/FileDownload'
import CloseIcon from '@mui/icons-material/Close'
// 匯入 Wails 綁定的 API 和類型
import * as AppAPI from '../wailsjs/wailsjs/go/app/App'
import type { models } from '../wailsjs/wailsjs/go/models'
import type { parser } from '../wailsjs/wailsjs/go/models'

// 使用 Wails 生成的類型別名
type LogEntry = models.LogEntry
type ParseError = parser.ParseError
import TabPanel from './components/TabPanel'
import LogTable from './components/LogTable'
import ErrorSummary from './components/ErrorSummary'
import ProgressIndicator from './components/ProgressIndicator'
import Dashboard, { type Statistics } from './components/Dashboard'
import ExportProgress from './components/ExportProgress'

// 檔案資訊介面
interface FileInfo {
  path: string
  name: string
  entries: LogEntry[]
  errorCount: number
  errorSamples: ParseError[]
  statistics: Statistics | null  // User Story 2: 統計資訊
  statTime: number  // 統計計算耗時（毫秒）
  isLoading: boolean
  error: string | null
}

function App() {
  const [files, setFiles] = useState<FileInfo[]>([])
  const [currentTab, setCurrentTab] = useState(0)
  const [currentSubTab, setCurrentSubTab] = useState(0)  // 二級標籤頁（Dashboard/日誌表格）
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // 匯出狀態
  const [exporting, setExporting] = useState(false)
  const [exportProgress, setExportProgress] = useState(0)
  const [exportMessage, setExportMessage] = useState('')
  const [exportWarnings, setExportWarnings] = useState<string[]>([])
  const [successMessage, setSuccessMessage] = useState<string | null>(null)

  // 檢查 Wails runtime 是否已載入
  useEffect(() => {
    const runtimeCheck = {
      window: typeof window,
      windowGo: typeof (window as any)['go'],
      windowGoKeys: (window as any)['go'] ? Object.keys((window as any)['go']) : [],
      appApi: typeof AppAPI,
      selectFile: typeof AppAPI.SelectFile
    }
    console.log('Wails runtime check:', runtimeCheck)
    
    // 如果 window.go 不存在，顯示錯誤
    if (!(window as any)['go']) {
      setError('Wails runtime 未載入！window.go 不存在')
    }
  }, [])

  // 處理檔案選擇
  const handleSelectFile = async () => {
    try {
      setError(null)
      setLoading(true)

      console.log('開始選擇檔案...')
      console.log('AppAPI.SelectFile 類型:', typeof AppAPI.SelectFile)
      console.log('window.go 存在:', typeof (window as any)['go'])
      
      // 呼叫 Wails API 選擇檔案
      // 不需要傳遞 context 參數
      const response = await AppAPI.SelectFile()
      
      console.log('SelectFile 回應:', response)
      
      if (!response || !response.success || !response.filePath) {
        console.log('檔案選擇失敗或取消:', response?.errorMessage)
        setLoading(false)
        if (response?.errorMessage && response.errorMessage !== '使用者取消選擇') {
          setError(response.errorMessage)
        }
        return
      }
      
      const filePath = response.filePath

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
        statistics: null,  // 統計資訊初始為 null
        statTime: 0,  // 統計耗時初始為 0
        isLoading: true,
        error: null,
      }
      
      const newFiles = [...files, newFile]
      setFiles(newFiles)
      setCurrentTab(newFiles.length - 1)

      // 呼叫 Wails API 解析檔案
      // 不需要傳遞 context 參數
      const result = await AppAPI.ParseFile({ filePath })
      
      // 檢查解析是否成功
      if (!result.success || !result.logFile) {
        throw new Error(result.errorMessage || '解析檔案失敗')
      }
      
      const logFile = result.logFile
      
      // 更新檔案資訊
      const updatedFiles = newFiles.map((f, i) => {
        if (i === newFiles.length - 1) {
          return {
            ...f,
            entries: logFile.entries || [],
            errorCount: logFile.errorLines || 0,
            errorSamples: result.errorSamples || [],
            statistics: logFile.statistics || null,
            statTime: logFile.statTime || 0,
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
      const errorMsg = err instanceof Error ? err.message : String(err)
      setError(`錯誤: ${errorMsg}`)
      setLoading(false)
      
      // 移除載入失敗的檔案
      setFiles(files.filter(f => !f.isLoading))
    }
  }

  // 處理分頁切換
  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue)
  }

  /**
   * 處理匯出至 Excel
   * 將當前檔案的分析結果匯出為 Excel 格式
   */
  const handleExport = async () => {
    // 檢查是否有開啟的檔案
    if (files.length === 0 || currentTab >= files.length) {
      setError('請先開啟一個日誌檔案')
      return
    }

    const currentFile = files[currentTab]
    
    // 檢查檔案是否載入完成
    if (currentFile.isLoading) {
      setError('檔案正在載入中，請稍後再試')
      return
    }

    // 檢查是否有錯誤
    if (currentFile.error) {
      setError('無法匯出有錯誤的檔案')
      return
    }

    try {
      setExporting(true)
      setExportProgress(0)
      setExportMessage('準備匯出...')
      setExportWarnings([])

      // 步驟 1: 選擇儲存位置
      setExportMessage('選擇儲存位置...')
      const defaultFileName = currentFile.name.replace(/\.(log|txt)$/, '') + '.xlsx'
      
      const saveResponse = await AppAPI.SelectSaveLocation(defaultFileName)
      
      if (!saveResponse.success) {
        if (saveResponse.errorMessage) {
          setError(saveResponse.errorMessage)
        }
        // 使用者取消選擇
        return
      }

      setExportProgress(10)
      setExportMessage('開始匯出資料...')

      // 步驟 2: 執行匯出
      const exportResponse = await AppAPI.ExportToExcel({
        filePath: currentFile.path,
        savePath: saveResponse.savePath,
      })

      if (!exportResponse.success) {
        setError(exportResponse.errorMessage || '匯出失敗')
        return
      }

      // 顯示警告訊息（如果有）
      if (exportResponse.warnings && exportResponse.warnings.length > 0) {
        setExportWarnings(exportResponse.warnings)
      }

      setExportProgress(100)
      setExportMessage('匯出完成！')

      // 格式化檔案大小
      const fileSizeStr = formatFileSize(exportResponse.fileSize)
      
      // 顯示成功訊息
      setTimeout(() => {
        setExporting(false)
        setSuccessMessage(`已成功匯出至 ${exportResponse.exportPath} (${fileSizeStr})`)
      }, 500)

    } catch (err) {
      console.error('匯出錯誤:', err)
      setError(err instanceof Error ? err.message : '匯出過程中發生未知錯誤')
    } finally {
      if (!exporting) {
        setExporting(false)
      }
    }
  }

  /**
   * 格式化檔案大小
   */
  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
    if (bytes < 1024 * 1024 * 1024) return `${(bytes / 1024 / 1024).toFixed(2)} MB`
    return `${(bytes / 1024 / 1024 / 1024).toFixed(2)} GB`
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
          <Button
            color="inherit"
            startIcon={<FileDownloadIcon />}
            onClick={handleExport}
            disabled={loading || files.length === 0 || exporting}
            sx={{ ml: 1 }}
          >
            匯出至 Excel
          </Button>
        </Toolbar>
      </AppBar>

      {/* 匯出進度對話框 */}
      <ExportProgress
        open={exporting}
        progress={exportProgress}
        message={exportMessage}
        warnings={exportWarnings}
      />

      {/* 成功通知 */}
      <Snackbar
        open={!!successMessage}
        autoHideDuration={6000}
        onClose={() => setSuccessMessage(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert 
          onClose={() => setSuccessMessage(null)} 
          severity="success" 
          sx={{ width: '100%' }}
          action={
            <IconButton
              size="small"
              aria-label="close"
              color="inherit"
              onClick={() => setSuccessMessage(null)}
            >
              <CloseIcon fontSize="small" />
            </IconButton>
          }
        >
          {successMessage}
        </Alert>
      </Snackbar>

      {/* 錯誤提示 */}
      {error && (
        <Alert severity="error" onClose={() => setError(null)} sx={{ m: 2 }}>
          {error}
        </Alert>
      )}
      
      {/* 除錯資訊 */}
      {!(window as any)['go'] && (
        <Alert severity="warning" sx={{ m: 2 }}>
          <Typography variant="body2">
            <strong>除錯資訊：</strong>Wails runtime 未載入
          </Typography>
          <Typography variant="caption" component="pre" sx={{ mt: 1 }}>
            {JSON.stringify({
              windowGo: typeof (window as any)['go'],
              location: window.location.href
            }, null, 2)}
          </Typography>
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
                  {/* 檔案資訊摘要 */}
                  <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
                    <Typography variant="h6" gutterBottom>
                      {file.name}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      總記錄數: {file.entries.length.toLocaleString()} | 錯誤數: {file.errorCount.toLocaleString()}
                    </Typography>
                  </Box>
                  
                  {/* 二級標籤頁 */}
                  <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                    <Tabs value={currentSubTab} onChange={(_e, v) => setCurrentSubTab(v)}>
                      <Tab label="統計儀表板" />
                      <Tab label="日誌明細" />
                    </Tabs>
                  </Box>
                  
                  {/* 統計儀表板 */}
                  <TabPanel value={currentSubTab} index={0}>
                    <Dashboard statistics={file.statistics} statTime={file.statTime} />
                  </TabPanel>
                  
                  {/* 日誌明細 */}
                  <TabPanel value={currentSubTab} index={1}>
                    {/* 錯誤摘要 */}
                    {file.errorCount > 0 && (
                      <Box sx={{ p: 2 }}>
                        <ErrorSummary
                          errorCount={file.errorCount}
                          errorSamples={file.errorSamples}
                        />
                      </Box>
                    )}
                    
                    {/* 日誌表格 */}
                    <LogTable entries={file.entries} />
                  </TabPanel>
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
