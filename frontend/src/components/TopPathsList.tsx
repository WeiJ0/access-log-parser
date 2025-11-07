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
  Path: string
  Count: number
  AvgBytes: number
  Methods: Record<string, number>
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

  // 取得主要 HTTP 方法
  const getPrimaryMethod = (methods: Record<string, number>): string => {
    if (!methods || Object.keys(methods).length === 0) return 'N/A'
    
    return Object.entries(methods)
      .sort(([, a], [, b]) => b - a)
      .map(([method]) => method)
      .slice(0, 2)
      .join(', ')
  }

  // 根據 HTTP 方法取得顏色
  const getMethodColor = (methodStr: string): 'primary' | 'success' | 'warning' | 'error' | 'default' => {
    if (methodStr.includes('GET')) return 'primary'
    if (methodStr.includes('POST')) return 'success'
    if (methodStr.includes('PUT') || methodStr.includes('PATCH')) return 'warning'
    if (methodStr.includes('DELETE')) return 'error'
    return 'default'
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
              <TableCell align="center">主要方法</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {topPaths.map((path, index) => {
              const primaryMethod = getPrimaryMethod(path.Methods)
              return (
                <TableRow key={`${path.Path}-${index}`} hover>
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
                      title={path.Path}
                    >
                      {path.Path}
                    </Typography>
                  </TableCell>
                  <TableCell align="right">
                    {path.Count.toLocaleString()}
                  </TableCell>
                  <TableCell align="right">
                    {(path.AvgBytes / 1024).toFixed(2)}
                  </TableCell>
                  <TableCell align="center">
                    <Box sx={{ display: 'flex', justifyContent: 'center', gap: 0.5 }}>
                      <Chip
                        label={primaryMethod}
                        size="small"
                        color={getMethodColor(primaryMethod)}
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
