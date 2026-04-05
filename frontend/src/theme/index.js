import { createTheme } from '@mui/material';

const getDesignTokens = (mode) => ({
  palette: {
    mode,
    ...(mode === 'light'
      ? {
          primary: { main: '#1976d2' },
          secondary: { main: '#42a5f5' },
          background: { default: '#f5f7fa', paper: '#ffffff' },
          text: { primary: '#1a237e', secondary: '#5c6bc0' },
        }
      : {
          primary: { main: '#42a5f5' },
          secondary: { main: '#90caf9' },
          background: { default: '#0d1117', paper: '#161b22' },
          text: { primary: '#e6edf3', secondary: '#8b949e' },
        }),
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h4: { fontWeight: 600 },
    h5: { fontWeight: 600 },
    h6: { fontWeight: 600 },
  },
  shape: { borderRadius: 12 },
  components: {
    MuiButton: {
      styleOverrides: {
        root: { textTransform: 'none', borderRadius: 10 },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: { borderRadius: 12, boxShadow: '0 4px 20px rgba(0,0,0,0.08)' },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: { borderRadius: 8 },
      },
    },
  },
});

export const buildTheme = (darkMode) =>
  createTheme(getDesignTokens(darkMode ? 'dark' : 'light'));
