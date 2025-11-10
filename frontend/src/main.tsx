// 應用程式入口點
// 文件路徑: frontend/src/main.tsx

import React from 'react'
import ReactDOM from 'react-dom/client'
import { ThemeProvider } from '@mui/material/styles'
import CssBaseline from '@mui/material/CssBaseline'
import App from './App'
import theme from './theme'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-material.css'

// 確保 root 元素存在
const rootElement = document.getElementById('root')
if (!rootElement) {
  throw new Error('無法找到 root 元素')
}

// 創建 React 根節點並渲染應用程式
ReactDOM.createRoot(rootElement).render(
  <React.StrictMode>
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <App />
    </ThemeProvider>
  </React.StrictMode>
)
