package exporter

import (
	"fmt"
	"sort"
	"strconv"
	
	"access-log-analyzer/internal/stats"
)

// FormatStatsStatistics 格式化 stats.Statistics 為二維字串陣列
// 這是臨時適配器函數，用於支援 stats.Statistics 類型
func (f *Formatter) FormatStatsStatistics(s *stats.Statistics) [][]string {
	if s == nil {
		return [][]string{{"統計項目", "數值"}}
	}
	
	result := make([][]string, 0)
	
	// 基本統計資料區塊
	result = append(result, []string{"===== 基本統計 ====="})
	result = append(result, []string{"統計項目", "數值"})
	result = append(result, []string{"總請求數", strconv.Itoa(s.TotalRequests)})
	result = append(result, []string{"唯一IP數量", strconv.Itoa(s.UniqueIPs)})
	result = append(result, []string{"唯一路徑數量", strconv.Itoa(s.UniquePaths)})
	result = append(result, []string{"總傳輸量 (位元組)", strconv.FormatInt(s.TotalBytes, 10)})
	result = append(result, []string{"總傳輸量 (MB)", fmt.Sprintf("%.2f", float64(s.TotalBytes)/(1024*1024))})
	result = append(result, []string{"平均回應大小 (位元組)", strconv.FormatInt(s.AverageResponseSize, 10)})
	
	// Top IP統計
	result = append(result, []string{""}) // 空行分隔
	result = append(result, []string{"===== Top 10 IP 位址 ====="})
	result = append(result, []string{"IP位址", "請求次數", "總流量(位元組)"})
	topIPsCount := len(s.TopIPs)
	if topIPsCount > 10 {
		topIPsCount = 10
	}
	for i := 0; i < topIPsCount; i++ {
		ip := s.TopIPs[i]
		result = append(result, []string{
			ip.IP,
			strconv.Itoa(ip.RequestCount),
			strconv.FormatInt(ip.TotalBytes, 10),
		})
	}
	
	// Top路徑統計
	result = append(result, []string{""})
	result = append(result, []string{"===== Top 10 請求路徑 ====="})
	result = append(result, []string{"路徑", "請求次數", "平均大小", "錯誤率(%)"})
	topPathsCount := len(s.TopPaths)
	if topPathsCount > 10 {
		topPathsCount = 10
	}
	for i := 0; i < topPathsCount; i++ {
		path := s.TopPaths[i]
		result = append(result, []string{
			path.Path,
			strconv.Itoa(path.RequestCount),
			strconv.FormatInt(path.AverageSize, 10),
			fmt.Sprintf("%.2f", path.ErrorRate),
		})
	}
	
	// 狀態碼分布
	result = append(result, []string{""})
	result = append(result, []string{"===== 狀態碼分布 ====="})
	result = append(result, []string{"類別", "次數"})
	result = append(result, []string{"成功 (2xx)", strconv.Itoa(s.StatusCodeDistribution.Success)})
	result = append(result, []string{"重定向 (3xx)", strconv.Itoa(s.StatusCodeDistribution.Redirection)})
	result = append(result, []string{"客戶端錯誤 (4xx)", strconv.Itoa(s.StatusCodeDistribution.ClientError)})
	result = append(result, []string{"伺服器錯誤 (5xx)", strconv.Itoa(s.StatusCodeDistribution.ServerError)})
	
	// 機器人統計
	result = append(result, []string{""})
	result = append(result, []string{"===== 機器人統計 ====="})
	result = append(result, []string{"統計項目", "數值"})
	result = append(result, []string{"總請求數", strconv.Itoa(s.BotStats.Total)})
	result = append(result, []string{"機器人請求數", strconv.Itoa(s.BotStats.BotRequests)})
	result = append(result, []string{"人類請求數", strconv.Itoa(s.BotStats.HumanRequests)})
	result = append(result, []string{"機器人百分比", fmt.Sprintf("%.2f%%", s.BotStats.BotPercentage)})
	
	// 機器人類型分布
	if len(s.BotStats.BotTypes) > 0 {
		result = append(result, []string{""})
		result = append(result, []string{"===== 機器人類型分布 ====="})
		result = append(result, []string{"類型", "次數"})
		
		// 排序機器人類型
		var botTypes []string
		for botType := range s.BotStats.BotTypes {
			botTypes = append(botTypes, botType)
		}
		sort.Slice(botTypes, func(i, j int) bool {
			return s.BotStats.BotTypes[botTypes[i]] > s.BotStats.BotTypes[botTypes[j]]
		})
		
		for _, botType := range botTypes {
			count := s.BotStats.BotTypes[botType]
			result = append(result, []string{
				botType,
				strconv.Itoa(count),
			})
		}
	}
	
	// Top Bots
	if len(s.BotStats.TopBots) > 0 {
		result = append(result, []string{""})
		result = append(result, []string{"===== Top 機器人 ====="})
		result = append(result, []string{"名稱", "請求次數", "百分比(%)"})
		
		topBotsCount := len(s.BotStats.TopBots)
		if topBotsCount > 10 {
			topBotsCount = 10
		}
		for i := 0; i < topBotsCount; i++ {
			bot := s.BotStats.TopBots[i]
			result = append(result, []string{
				bot.Name,
				strconv.Itoa(bot.Count),
				fmt.Sprintf("%.2f", bot.Percentage),
			})
		}
	}
	
	return result
}
