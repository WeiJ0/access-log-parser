// ProgressIndicator 進度指示器組件
// 顯示檔案解析的載入狀態
// 文件路徑: frontend/src/components/ProgressIndicator.tsx

import {
  Box,
  CircularProgress,
  Typography,
  LinearProgress,
} from '@mui/material'

interface ProgressIndicatorProps {
  loading: boolean
  message?: string
  progress?: number
  variant?: 'circular' | 'linear'
}

/**
 * ProgressIndicator 組件
 * 顯示載入狀態和可選的進度百分比
 * 
 * @param loading - 是否正在載入
 * @param message - 載入訊息（可選）
 * @param progress - 進度百分比 0-100（可選，僅用於 linear 變體）
 * @param variant - 顯示樣式：circular（圓形）或 linear（線性）
 */
export default function ProgressIndicator({
  loading,
  message = '載入中...',
  progress,
  variant = 'circular',
}: ProgressIndicatorProps) {
  if (!loading) {
    return null
  }

  if (variant === 'linear') {
    return (
      <Box sx={{ width: '100%', p: 2 }}>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
          {message}
        </Typography>
        {progress !== undefined ? (
          <>
            <LinearProgress variant="determinate" value={progress} />
            <Typography
              variant="caption"
              color="text.secondary"
              sx={{ mt: 0.5, display: 'block' }}
            >
              {progress.toFixed(0)}%
            </Typography>
          </>
        ) : (
          <LinearProgress />
        )}
      </Box>
    )
  }

  // 預設：circular 變體
  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        p: 4,
        gap: 2,
      }}
    >
      <CircularProgress size={48} />
      <Typography variant="body1" color="text.secondary">
        {message}
      </Typography>
      {progress !== undefined && (
        <Typography variant="caption" color="text.secondary">
          {progress.toFixed(0)}%
        </Typography>
      )}
    </Box>
  )
}
