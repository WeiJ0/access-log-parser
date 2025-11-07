// Dashboard 元件 - 顯示統計分析結果
// 文件路徑: frontend/src/components/Dashboard.tsx
// 用途: User Story 2 - 統計資訊儀表板（T072）

import { Box, Grid, Paper, Typography } from '@mui/material'
import TopIPsList from './TopIPsList'
import TopPathsList from './TopPathsList'
import StatusCodeDistribution from './StatusCodeDistribution'
import BotDetection from './BotDetection'

// 統計資料介面（對應 Go internal/stats/statistics.go）
export interface Statistics {
  // 基本統計
  TotalRequests: number
  UniqueIPs: number
  TotalBytes: number
  
  // Top 10 IP（對應 IPStatistics）
  TopIPs: Array<{
    IP: string
    Count: number
    TotalBytes: number
    UniqueURLs: number
  }>
  
  // Top 10 路徑（對應 PathStatistics）
  TopPaths: Array<{
    Path: string
    Count: number
    AvgBytes: number
    Methods: Record<string, number>
  }>
  
  // 狀態碼分布（對應 StatusCodeStatistics）
  StatusCodeDist: Array<{
    Code: number
    Count: number
    Percentage: number
  }>
  
  // 機器人偵測
  BotRequests: number
  BotPercentage: number
  TopBots: Array<{
    UserAgent: string
    Count: number
    Percentage: number
  }>
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
              {statistics.TotalRequests.toLocaleString()}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Typography variant="body2" color="text.secondary">
              唯一 IP 數
            </Typography>
            <Typography variant="h5">
              {statistics.UniqueIPs.toLocaleString()}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Typography variant="body2" color="text.secondary">
              總流量
            </Typography>
            <Typography variant="h5">
              {(statistics.TotalBytes / (1024 * 1024)).toFixed(2)} MB
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
          <TopIPsList topIPs={statistics.TopIPs} />
        </Grid>

        {/* Top 10 路徑 */}
        <Grid item xs={12} md={6}>
          <TopPathsList topPaths={statistics.TopPaths} />
        </Grid>

        {/* 狀態碼分布 */}
        <Grid item xs={12} md={6}>
          <StatusCodeDistribution statusCodes={statistics.StatusCodeDist} />
        </Grid>

        {/* 機器人偵測 */}
        <Grid item xs={12} md={6}>
          <BotDetection
            botRequests={statistics.BotRequests}
            botPercentage={statistics.BotPercentage}
            topBots={statistics.TopBots}
          />
        </Grid>
      </Grid>
    </Box>
  )
}

export default Dashboard
