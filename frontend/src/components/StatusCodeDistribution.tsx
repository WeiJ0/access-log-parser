// StatusCodeDistribution 元件 - 顯示 HTTP 狀態碼分布
// 文件路徑: frontend/src/components/StatusCodeDistribution.tsx
// 用途: User Story 2 - 狀態碼分布圖表（T075）

import {
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  LinearProgress,
  Box,
  Chip,
  Grid,
} from '@mui/material'

interface StatusCodeDistribution {
  success: number
  redirection: number
  clientError: number
  serverError: number
  details: Record<number, number>
}

interface StatusCodeDistributionProps {
  distribution: StatusCodeDistribution
}

/**
 * StatusCodeDistribution 元件 - 顯示 HTTP 狀態碼分布統計
 * 
 * @param distribution - 狀態碼分布統計資料
 */
function StatusCodeDistribution({ distribution }: StatusCodeDistributionProps) {
  if (!distribution || !distribution.details) {
    return (
      <Paper sx={{ p: 3, textAlign: 'center' }}>
        <Typography variant="body2" color="text.secondary">
          無狀態碼統計資料
        </Typography>
      </Paper>
    )
  }

  // 根據狀態碼取得類別和顏色
  const getStatusInfo = (code: number): { category: string; color: string; chipColor: 'success' | 'info' | 'warning' | 'error' | 'default' } => {
    if (code >= 200 && code < 300) {
      return { category: '成功', color: '#4caf50', chipColor: 'success' }
    }
    if (code >= 300 && code < 400) {
      return { category: '重新導向', color: '#2196f3', chipColor: 'info' }
    }
    if (code >= 400 && code < 500) {
      return { category: '客戶端錯誤', color: '#ff9800', chipColor: 'warning' }
    }
    if (code >= 500) {
      return { category: '伺服器錯誤', color: '#f44336', chipColor: 'error' }
    }
    return { category: '其他', color: '#9e9e9e', chipColor: 'default' }
  }

  // 計算總請求數
  const totalRequests = Object.values(distribution.details).reduce((sum, count) => sum + count, 0)

  // 將 details 轉換為陣列並排序
  const sortedDetails = Object.entries(distribution.details)
    .map(([code, count]) => ({
      code: Number(code),
      count,
      percentage: totalRequests > 0 ? (count / totalRequests * 100) : 0
    }))
    .sort((a, b) => b.count - a.count)

  return (
    <Paper sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        HTTP 狀態碼分布
      </Typography>
      
      {/* 分類摘要 */}
      <Box sx={{ mb: 2 }}>
        <Grid container spacing={1}>
          <Grid item xs={6} sm={3}>
            <Box sx={{ p: 1, bgcolor: '#e8f5e9', borderRadius: 1 }}>
              <Typography variant="caption" color="text.secondary">成功 (2xx)</Typography>
              <Typography variant="h6">{distribution.success.toLocaleString()}</Typography>
            </Box>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Box sx={{ p: 1, bgcolor: '#e3f2fd', borderRadius: 1 }}>
              <Typography variant="caption" color="text.secondary">重新導向 (3xx)</Typography>
              <Typography variant="h6">{distribution.redirection.toLocaleString()}</Typography>
            </Box>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Box sx={{ p: 1, bgcolor: '#fff3e0', borderRadius: 1 }}>
              <Typography variant="caption" color="text.secondary">客戶端錯誤 (4xx)</Typography>
              <Typography variant="h6">{distribution.clientError.toLocaleString()}</Typography>
            </Box>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Box sx={{ p: 1, bgcolor: '#ffebee', borderRadius: 1 }}>
              <Typography variant="caption" color="text.secondary">伺服器錯誤 (5xx)</Typography>
              <Typography variant="h6">{distribution.serverError.toLocaleString()}</Typography>
            </Box>
          </Grid>
        </Grid>
      </Box>

      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        總請求數: {totalRequests.toLocaleString()}
      </Typography>
      <TableContainer>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>狀態碼</TableCell>
              <TableCell>類別</TableCell>
              <TableCell align="right">次數</TableCell>
              <TableCell align="right">百分比</TableCell>
              <TableCell>分布</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {sortedDetails.map((stat) => {
              const statusInfo = getStatusInfo(stat.code)
              return (
                <TableRow key={stat.code} hover>
                  <TableCell>
                    <Typography variant="body2" sx={{ fontFamily: 'monospace', fontWeight: 'bold' }}>
                      {stat.code}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={statusInfo.category}
                      size="small"
                      color={statusInfo.chipColor}
                    />
                  </TableCell>
                  <TableCell align="right">
                    {stat.count.toLocaleString()}
                  </TableCell>
                  <TableCell align="right">
                    {stat.percentage.toFixed(2)}%
                  </TableCell>
                  <TableCell>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <LinearProgress
                        variant="determinate"
                        value={Math.min(stat.percentage, 100)}
                        sx={{
                          flexGrow: 1,
                          height: 8,
                          borderRadius: 1,
                          backgroundColor: '#e0e0e0',
                          '& .MuiLinearProgress-bar': {
                            backgroundColor: statusInfo.color,
                          },
                        }}
                      />
                    </Box>
                  </TableCell>
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </TableContainer>
    </Paper>
  )
}

export default StatusCodeDistribution
