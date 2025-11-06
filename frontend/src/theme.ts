// Material-UI 主題配置
// 文件路徑: frontend/src/theme.ts

import { createTheme, zhTW } from '@mui/material/styles'

/**
 * 應用程式主題配置
 * 使用 Material Design 3 色彩系統和繁體中文本地化
 */
export const theme = createTheme(
  {
    palette: {
      mode: 'light',
      primary: {
        main: '#1976d2', // 藍色 - 主要操作按鈕
        light: '#42a5f5',
        dark: '#1565c0',
        contrastText: '#fff',
      },
      secondary: {
        main: '#9c27b0', // 紫色 - 次要強調
        light: '#ba68c8',
        dark: '#7b1fa2',
        contrastText: '#fff',
      },
      error: {
        main: '#d32f2f', // 紅色 - 錯誤訊息
        light: '#ef5350',
        dark: '#c62828',
      },
      warning: {
        main: '#ed6c02', // 橘色 - 警告訊息
        light: '#ff9800',
        dark: '#e65100',
      },
      info: {
        main: '#0288d1', // 淺藍 - 資訊提示
        light: '#03a9f4',
        dark: '#01579b',
      },
      success: {
        main: '#2e7d32', // 綠色 - 成功訊息
        light: '#4caf50',
        dark: '#1b5e20',
      },
      background: {
        default: '#fafafa',
        paper: '#fff',
      },
      text: {
        primary: 'rgba(0, 0, 0, 0.87)',
        secondary: 'rgba(0, 0, 0, 0.6)',
        disabled: 'rgba(0, 0, 0, 0.38)',
      },
    },
    typography: {
      fontFamily: [
        '-apple-system',
        'BlinkMacSystemFont',
        '"Microsoft JhengHei"', // 繁體中文字型
        '"Segoe UI"',
        'Roboto',
        '"Helvetica Neue"',
        'Arial',
        'sans-serif',
      ].join(','),
      h1: {
        fontSize: '2.5rem',
        fontWeight: 500,
        lineHeight: 1.2,
      },
      h2: {
        fontSize: '2rem',
        fontWeight: 500,
        lineHeight: 1.3,
      },
      h3: {
        fontSize: '1.75rem',
        fontWeight: 500,
        lineHeight: 1.4,
      },
      h4: {
        fontSize: '1.5rem',
        fontWeight: 500,
        lineHeight: 1.4,
      },
      h5: {
        fontSize: '1.25rem',
        fontWeight: 500,
        lineHeight: 1.5,
      },
      h6: {
        fontSize: '1rem',
        fontWeight: 500,
        lineHeight: 1.6,
      },
      body1: {
        fontSize: '1rem',
        lineHeight: 1.5,
      },
      body2: {
        fontSize: '0.875rem',
        lineHeight: 1.43,
      },
      button: {
        textTransform: 'none', // 保持按鈕文字原始大小寫
        fontWeight: 500,
      },
    },
    shape: {
      borderRadius: 8, // 統一圓角
    },
    components: {
      MuiButton: {
        styleOverrides: {
          root: {
            borderRadius: 8,
            padding: '8px 16px',
          },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: {
            borderRadius: 12,
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
          },
        },
      },
      MuiAppBar: {
        styleOverrides: {
          root: {
            boxShadow: '0 1px 3px rgba(0,0,0,0.12)',
          },
        },
      },
      MuiTabs: {
        styleOverrides: {
          root: {
            minHeight: 48,
          },
        },
      },
      MuiTab: {
        styleOverrides: {
          root: {
            minHeight: 48,
            textTransform: 'none',
            fontWeight: 500,
          },
        },
      },
    },
  },
  zhTW // Material-UI 繁體中文本地化
)

// 深色模式主題（備用）
export const darkTheme = createTheme(
  {
    palette: {
      mode: 'dark',
      primary: {
        main: '#90caf9',
      },
      secondary: {
        main: '#ce93d8',
      },
      background: {
        default: '#121212',
        paper: '#1e1e1e',
      },
    },
    typography: theme.typography,
    shape: theme.shape,
    components: theme.components,
  },
  zhTW
)

export default theme
