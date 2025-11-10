// Dashboard 元件 - 顯示統計分析結果
// 文件路徑: frontend/src/components/Dashboard.tsx
// 用途: User Story 2 - 統計資訊儀表板（T072）

import { Box, Grid, Paper, Typography } from '@mui/material'
import TopIPsList from './TopIPsList'
import TopPathsList from './TopPathsList'
import StatusCodeDistribution from './StatusCodeDistribution'
import BotDetection from './BotDetection'

// 統計資料介面（對應 Go internal/stats/statistics.go）
// 注意：欄位名稱必須與 Go JSON 標籤匹配（小寫開頭）
export interface Statistics {
  // 基本統計
  totalRequests: number          // 總請求數
  uniqueIPs: number              // 唯一 IP 數量
  uniquePaths: number            // 唯一路徑數量
  totalBytes: number             // 總傳輸量（位元組）
  averageResponseSize: number    // 平均回應大小
  
  // Top IP 統計
  topIPs: Array<{
    ip: string
    requestCount: number
    totalBytes: number
  }>
  
  // Top 路徑統計
  topPaths: Array<{
    path: string
    requestCount: number
    averageSize: number
    errorRate: number
  }>
  
  // 狀態碼分布
  statusCodeDistribution: {
    success: number       // 2xx 成功
    redirection: number   // 3xx 重定向
    clientError: number   // 4xx 客戶端錯誤
    serverError: number   // 5xx 伺服器錯誤
    details: Record<number, number>  // 詳細狀態碼分布
  }
  
  // 機器人統計
  botStats: {
    total: number              // 總請求數
    botRequests: number        // 機器人請求數
    humanRequests: number      // 人類請求數
    botPercentage: number      // 機器人百分比
    botTypes: Record<string, number>  // 機器人類型分布
    topBots: Array<{
      name: string
      count: number
      percentage: number
    }>
  }
}

interface DashboardProps {
  statistics: Statistics | null
  statTime: number  // 統計計算耗時（毫秒）
}

/**
 * Dashboard 元件 - 顯示日誌統計分析結果
 * 
 * @param statistics - 統計資料物件
 * @param statTime - 統計計算耗時（毫秒）
 */
function Dashboard({ statistics, statTime }: DashboardProps) {
  if (!statistics) {
    return (
      <Box sx={{ p: 3, textAlign: 'center' }}>
        <Typography variant="body1" color="text.secondary">
          無統計資料
        </Typography>
      </Box>
    )
  }

  return (
    <Box sx={{ p: 3 }}>
      {/* 基本統計摘要 */}
      <Paper sx={{ p: 2, mb: 3 }}>
        <Typography variant="h6" gutterBottom>
          統計摘要
        </Typography>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={6} md={3}>
            <Typography variant="body2" color="text.secondary">
              總請求數
            </Typography>
            <Typography variant="h5">
              {statistics.totalRequests.toLocaleString()}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Typography variant="body2" color="text.secondary">
              唯一 IP 數
            </Typography>
            <Typography variant="h5">
              {statistics.uniqueIPs.toLocaleString()}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Typography variant="body2" color="text.secondary">
              總流量
            </Typography>
            <Typography variant="h5">
              {(statistics.totalBytes / (1024 * 1024)).toFixed(2)} MB
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Typography variant="body2" color="text.secondary">
              計算耗時
            </Typography>
            <Typography variant="h5">
              {statTime} ms
            </Typography>
          </Grid>
        </Grid>
      </Paper>

      {/* 統計圖表區域 */}
      <Grid container spacing={3}>
        {/* Top 10 IP */}
        <Grid item xs={12} md={6}>
          <TopIPsList topIPs={statistics.topIPs} />
        </Grid>

        {/* Top 10 路徑 */}
        <Grid item xs={12} md={6}>
          <TopPathsList topPaths={statistics.topPaths} />
        </Grid>

        {/* 狀態碼分布 */}
        <Grid item xs={12} md={6}>
          <StatusCodeDistribution distribution={statistics.statusCodeDistribution} />
        </Grid>

        {/* 機器人偵測 */}
        <Grid item xs={12} md={6}>
          <BotDetection
            botRequests={statistics.botStats.botRequests}
            botPercentage={statistics.botStats.botPercentage}
            topBots={statistics.botStats.topBots}
          />
        </Grid>
      </Grid>
    </Box>
  )
}

export default Dashboard
