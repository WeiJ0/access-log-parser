import React from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  Box,
  Typography,
  LinearProgress,
  Alert,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
} from '@mui/material'
import WarningIcon from '@mui/icons-material/Warning'

/**
 * 匯出進度對話框的屬性
 */
interface ExportProgressProps {
  /** 是否顯示對話框 */
  open: boolean
  /** 進度百分比 (0-100) */
  progress: number
  /** 當前狀態訊息 */
  message: string
  /** 警告訊息列表 */
  warnings: string[]
}

/**
 * 匯出進度對話框組件
 * 
 * 顯示 Excel 匯出的即時進度、狀態訊息和警告資訊
 */
export const ExportProgress: React.FC<ExportProgressProps> = ({
  open,
  progress,
  message,
  warnings,
}) => {
  return (
    <Dialog open={open} maxWidth="sm" fullWidth disableEscapeKeyDown>
      <DialogTitle>正在匯出至 Excel</DialogTitle>
      <DialogContent>
        <Box sx={{ py: 2 }}>
          {/* 進度條 */}
          <Box sx={{ mb: 2 }}>
            <LinearProgress 
              variant="determinate" 
              value={progress} 
              sx={{ height: 8, borderRadius: 1 }}
            />
            <Typography 
              variant="body2" 
              color="text.secondary" 
              align="center" 
              sx={{ mt: 1 }}
            >
              {progress.toFixed(0)}%
            </Typography>
          </Box>

          {/* 狀態訊息 */}
          <Typography variant="body1" sx={{ mb: 2 }}>
            {message}
          </Typography>

          {/* 警告訊息 */}
          {warnings.length > 0 && (
            <Alert severity="warning" sx={{ mt: 2 }}>
              <Typography variant="subtitle2" gutterBottom>
                注意事項：
              </Typography>
              <List dense>
                {warnings.map((warning, index) => (
                  <ListItem key={index} sx={{ py: 0.5 }}>
                    <ListItemIcon sx={{ minWidth: 32 }}>
                      <WarningIcon fontSize="small" />
                    </ListItemIcon>
                    <ListItemText 
                      primary={warning} 
                      primaryTypographyProps={{ variant: 'body2' }}
                    />
                  </ListItem>
                ))}
              </List>
            </Alert>
          )}
        </Box>
      </DialogContent>
    </Dialog>
  )
}

export default ExportProgress
