// Wails API Mock - 暫時使用直到 Wails 綁定生成
// 文件路徑: frontend/src/services/wailsMock.ts

export const SelectFile = async (): Promise<any> => {
  console.warn('SelectFile mock called - Wails bindings not yet generated')
  return {
    success: false,
    error: 'Wails bindings not yet generated'
  }
}

// Wails 自動將 Go 的 ParseFile(ctx, req) 轉換為 TypeScript 的 ParseFile(filePath)
export const ParseFile = async (filePath: string): Promise<any> => {
  console.warn('ParseFile mock called - Wails bindings not yet generated')
  return {
    success: false,
    error: 'Wails bindings not yet generated',
    filePath  // 回傳供除錯用
  }
}

