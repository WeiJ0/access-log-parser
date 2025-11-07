// TopPathsList 元件 - 顯示 Top 10 路徑統計
// 文件路徑: frontend/src/components/TopPathsList.tsx
// 用途: User Story 2 - Top 10 路徑列表（T074）

import {
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Box,
} from '@mui/material'

interface PathStatistics {
  path: string
  requestCount: number
  averageSize: number
  errorRate: number
}

interface TopPathsListProps {
  topPaths: PathStatistics[]
}

/**
 * TopPathsList 元件 - 顯示訪問次數最多的前 10 個路徑
 * 
 * @param topPaths - 路徑統計資料陣列（已排序）
 */
function TopPathsList({ topPaths }: TopPathsListProps) {
  if (!topPaths || topPaths.length === 0) {
    return (
      <Paper sx={{ p: 3, textAlign: 'center' }}>
        <Typography variant="body2" color="text.secondary">
          無路徑統計資料
        </Typography>
      </Paper>
    )
  }

  // 根據錯誤率決定顯示顏色
  const getErrorRateColor = (errorRate: number): 'success' | 'warning' | 'error' | 'default' => {
    if (errorRate === 0) return 'success'
    if (errorRate < 5) return 'default'
    if (errorRate < 20) return 'warning'
    return 'error'
  }

  return (
    <Paper sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Top 10 路徑
      </Typography>
      <TableContainer>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>排名</TableCell>
              <TableCell>路徑</TableCell>
              <TableCell align="right">請求次數</TableCell>
              <TableCell align="right">平均回應 (KB)</TableCell>
              <TableCell align="center">錯誤率</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {topPaths.map((path, index) => {
              return (
                <TableRow key={`${path.path}-${index}`} hover>
                  <TableCell>{index + 1}</TableCell>
                  <TableCell>
                    <Typography
                      variant="body2"
                      sx={{
                        fontFamily: 'monospace',
                        maxWidth: '300px',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                      }}
                      title={path.path}
                    >
                      {path.path}
                    </Typography>
                  </TableCell>
                  <TableCell align="right">
                    {path.requestCount.toLocaleString()}
                  </TableCell>
                  <TableCell align="right">
                    {(path.averageSize / 1024).toFixed(2)}
                  </TableCell>
                  <TableCell align="center">
                    <Box sx={{ display: 'flex', justifyContent: 'center', gap: 0.5 }}>
                      <Chip
                        label={`${path.errorRate.toFixed(2)}%`}
                        size="small"
                        color={getErrorRateColor(path.errorRate)}
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

export default TopPathsList
