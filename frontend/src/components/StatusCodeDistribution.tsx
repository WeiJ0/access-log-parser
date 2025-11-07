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
} from '@mui/material'

interface StatusCodeStat {
  Code: number
  Count: number
  Percentage: number
}

interface StatusCodeDistributionProps {
  statusCodes: StatusCodeStat[]
}

/**
 * StatusCodeDistribution 元件 - 顯示 HTTP 狀態碼分布統計
 * 
 * @param statusCodes - 狀態碼統計資料陣列（已排序）
 */
function StatusCodeDistribution({ statusCodes }: StatusCodeDistributionProps) {
  if (!statusCodes || statusCodes.length === 0) {
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
  const totalRequests = statusCodes.reduce((sum, stat) => sum + stat.Count, 0)

  return (
    <Paper sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        HTTP 狀態碼分布
      </Typography>
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
            {statusCodes.map((stat) => {
              const statusInfo = getStatusInfo(stat.Code)
              return (
                <TableRow key={stat.Code} hover>
                  <TableCell>
                    <Typography variant="body2" sx={{ fontFamily: 'monospace', fontWeight: 'bold' }}>
                      {stat.Code}
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
                    {stat.Count.toLocaleString()}
                  </TableCell>
                  <TableCell align="right">
                    {stat.Percentage.toFixed(2)}%
                  </TableCell>
                  <TableCell>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <LinearProgress
                        variant="determinate"
                        value={stat.Percentage}
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
