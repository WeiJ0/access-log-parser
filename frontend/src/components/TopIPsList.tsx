// TopIPsList 元件 - 顯示 Top 10 IP 統計
// 文件路徑: frontend/src/components/TopIPsList.tsx
// 用途: User Story 2 - Top 10 IP 列表（T073）

import {
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from '@mui/material'

interface IPStatistics {
  IP: string
  Count: number
  TotalBytes: number
  UniqueURLs: number
}

interface TopIPsListProps {
  topIPs: IPStatistics[]
}

/**
 * TopIPsList 元件 - 顯示訪問次數最多的前 10 個 IP 位址
 * 
 * @param topIPs - IP 統計資料陣列（已排序）
 */
function TopIPsList({ topIPs }: TopIPsListProps) {
  if (!topIPs || topIPs.length === 0) {
    return (
      <Paper sx={{ p: 3, textAlign: 'center' }}>
        <Typography variant="body2" color="text.secondary">
          無 IP 統計資料
        </Typography>
      </Paper>
    )
  }

  return (
    <Paper sx={{ p: 2 }}>
      <Typography variant="h6" gutterBottom>
        Top 10 IP 位址
      </Typography>
      <TableContainer>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>排名</TableCell>
              <TableCell>IP 位址</TableCell>
              <TableCell align="right">請求次數</TableCell>
              <TableCell align="right">流量 (MB)</TableCell>
              <TableCell align="right">唯一路徑</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {topIPs.map((ip, index) => (
              <TableRow key={ip.IP} hover>
                <TableCell>{index + 1}</TableCell>
                <TableCell>
                  <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                    {ip.IP}
                  </Typography>
                </TableCell>
                <TableCell align="right">
                  {ip.Count.toLocaleString()}
                </TableCell>
                <TableCell align="right">
                  {(ip.TotalBytes / (1024 * 1024)).toFixed(2)}
                </TableCell>
                <TableCell align="right">
                  {ip.UniqueURLs}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Paper>
  )
}

export default TopIPsList
