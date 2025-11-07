// BotDetection 元件 - 顯示機器人流量偵測結果
// 文件路徑: frontend/src/components/BotDetection.tsx
// 用途: User Story 2 - 機器人偵測統計（T076）

import {
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Box,
  LinearProgress,
  Alert,
} from '@mui/material'
import SmartToyIcon from '@mui/icons-material/SmartToy'
import PersonIcon from '@mui/icons-material/Person'

interface BotStat {
  name: string
  count: number
  percentage: number
}

interface BotDetectionProps {
  botRequests: number
  botPercentage: number
  topBots: BotStat[]
}

/**
 * BotDetection 元件 - 顯示機器人流量偵測和分析結果
 * 
 * @param botRequests - 機器人請求總數
 * @param botPercentage - 機器人流量百分比
 * @param topBots - Top 10 機器人列表
 */
function BotDetection({ botRequests, botPercentage, topBots }: BotDetectionProps) {
  // 提供預設值以避免 undefined 錯誤
  const safeBotRequests = botRequests ?? 0
  const safeBotPercentage = botPercentage ?? 0
  const safeTopBots = topBots ?? []
  
  // 判斷機器人流量是否異常高
  const isHighBotTraffic = safeBotPercentage > 50

  return (
    <Paper sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        機器人流量偵測
      </Typography>

      {/* 機器人流量摘要 */}
      <Box sx={{ mb: 3 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <SmartToyIcon color="action" />
            <Box>
              <Typography variant="body2" color="text.secondary">
                機器人請求
              </Typography>
              <Typography variant="h6">
                {safeBotRequests.toLocaleString()} ({safeBotPercentage.toFixed(2)}%)
              </Typography>
            </Box>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <PersonIcon color="primary" />
            <Box>
              <Typography variant="body2" color="text.secondary">
                人類請求
              </Typography>
              <Typography variant="h6">
                {((100 - safeBotPercentage).toFixed(2))}%
              </Typography>
            </Box>
          </Box>
        </Box>

        {/* 流量比例視覺化 */}
        <Box sx={{ mb: 2 }}>
          <Box sx={{ display: 'flex', height: 24, borderRadius: 1, overflow: 'hidden' }}>
            <Box
              sx={{
                width: `${safeBotPercentage}%`,
                backgroundColor: isHighBotTraffic ? '#ff9800' : '#9e9e9e',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              {safeBotPercentage > 10 && (
                <Typography variant="caption" sx={{ color: 'white', fontWeight: 'bold' }}>
                  機器人 {safeBotPercentage.toFixed(1)}%
                </Typography>
              )}
            </Box>
            <Box
              sx={{
                width: `${100 - safeBotPercentage}%`,
                backgroundColor: '#2196f3',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              {(100 - safeBotPercentage) > 10 && (
                <Typography variant="caption" sx={{ color: 'white', fontWeight: 'bold' }}>
                  人類 {(100 - safeBotPercentage).toFixed(1)}%
                </Typography>
              )}
            </Box>
          </Box>
        </Box>

        {/* 異常警告 */}
        {isHighBotTraffic && (
          <Alert severity="warning" sx={{ mb: 2 }}>
            機器人流量異常偏高（&gt; 50%），建議檢查是否有爬蟲或攻擊行為
          </Alert>
        )}
      </Box>

      {/* Top 10 機器人列表 */}
      {safeTopBots && safeTopBots.length > 0 && (
        <>
          <Typography variant="subtitle2" gutterBottom>
            Top 10 機器人 User-Agent
          </Typography>
          <TableContainer>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>排名</TableCell>
                  <TableCell>User-Agent</TableCell>
                  <TableCell align="right">請求次數</TableCell>
                  <TableCell align="right">百分比</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {safeTopBots.map((bot, index) => (
                  <TableRow key={`${bot.name}-${index}`} hover>
                    <TableCell>{index + 1}</TableCell>
                    <TableCell>
                      <Typography
                        variant="body2"
                        sx={{
                          fontFamily: 'monospace',
                          fontSize: '0.75rem',
                          maxWidth: '250px',
                          overflow: 'hidden',
                          textOverflow: 'ellipsis',
                          whiteSpace: 'nowrap',
                        }}
                        title={bot.name}
                      >
                        {bot.name}
                      </Typography>
                    </TableCell>
                    <TableCell align="right">
                      {bot.count.toLocaleString()}
                    </TableCell>
                    <TableCell align="right">
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Typography variant="body2">
                          {bot.percentage.toFixed(2)}%
                        </Typography>
                        <LinearProgress
                          variant="determinate"
                          value={bot.percentage}
                          sx={{
                            width: 60,
                            height: 6,
                            borderRadius: 1,
                          }}
                        />
                      </Box>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </>
      )}

      {(!safeTopBots || safeTopBots.length === 0) && (
        <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', mt: 2 }}>
          未偵測到機器人流量
        </Typography>
      )}
    </Paper>
  )
}

export default BotDetection
