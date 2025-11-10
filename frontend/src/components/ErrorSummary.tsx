// ErrorSummary 錯誤摘要組件
// 顯示解析過程中的錯誤數量和樣本
// 文件路徑: frontend/src/components/ErrorSummary.tsx

import {
  Alert,
  AlertTitle,
  Box,
  Typography,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  List,
  ListItem,
  ListItemText,
} from '@mui/material'
import ExpandMoreIcon from '@mui/icons-material/ExpandMore'
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline'
import type { parser } from '../../wailsjs/wailsjs/go/models'

interface ErrorSummaryProps {
  errorCount: number
  errorSamples: parser.ParseError[]
  maxSamples?: number
}

/**
 * ErrorSummary 組件
 * 顯示解析錯誤的摘要資訊和錯誤樣本
 * 
 * @param errorCount - 錯誤總數
 * @param errorSamples - 錯誤樣本陣列（原始錯誤行）
 * @param maxSamples - 最多顯示的樣本數量（預設：10）
 */
export default function ErrorSummary({
  errorCount,
  errorSamples,
  maxSamples = 10,
}: ErrorSummaryProps) {
  // 如果沒有錯誤，不顯示組件
  if (errorCount === 0) {
    return null
  }

  // 取得要顯示的錯誤樣本
  const samplesToShow = errorSamples.slice(0, maxSamples)
  const hasMoreSamples = errorSamples.length > maxSamples

  return (
    <Box sx={{ mb: 2 }}>
      <Accordion defaultExpanded={errorCount > 0}>
        <AccordionSummary
          expandIcon={<ExpandMoreIcon />}
          aria-controls="error-summary-content"
          id="error-summary-header"
        >
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <ErrorOutlineIcon color="warning" />
            <Typography variant="h6">
              解析錯誤摘要
            </Typography>
          </Box>
        </AccordionSummary>
        <AccordionDetails>
          <Alert severity="warning" icon={false}>
            <AlertTitle>
              發現 {errorCount} 個無法解析的日誌行
            </AlertTitle>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              這些行已被跳過，不會出現在統計資料中。
              {errorSamples.length > 0 && '以下是部分錯誤樣本：'}
            </Typography>

            {samplesToShow.length > 0 && (
              <List dense sx={{ bgcolor: 'background.paper', borderRadius: 1 }}>
                {samplesToShow.map((sample, index) => (
                  <ListItem
                    key={index}
                    sx={{
                      borderLeft: 3,
                      borderColor: 'warning.main',
                      mb: 1,
                      bgcolor: 'grey.50',
                    }}
                  >
                    <ListItemText
                      primary={
                        <Typography
                          variant="body2"
                          component="pre"
                          sx={{
                            fontFamily: 'monospace',
                            fontSize: '0.875rem',
                            whiteSpace: 'pre-wrap',
                            wordBreak: 'break-all',
                            margin: 0,
                          }}
                        >
                          {sample.line}
                        </Typography>
                      }
                      secondary={`第 ${sample.lineNumber} 行：${sample.error}`}
                    />
                  </ListItem>
                ))}
              </List>
            )}

            {hasMoreSamples && (
              <Typography
                variant="body2"
                color="text.secondary"
                sx={{ mt: 2, fontStyle: 'italic' }}
              >
                還有 {errorSamples.length - maxSamples} 個錯誤樣本未顯示
              </Typography>
            )}

            {errorCount > 100 && (
              <Alert severity="info" sx={{ mt: 2 }}>
                <Typography variant="body2">
                  💡 提示：錯誤數量較多，建議檢查日誌檔案格式是否正確。
                  系統預期的格式為 Apache Combined Log Format。
                </Typography>
              </Alert>
            )}
          </Alert>
        </AccordionDetails>
      </Accordion>
    </Box>
  )
}
