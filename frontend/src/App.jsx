import React, { lazy, Suspense, useState, useMemo } from 'react';
import { Routes, Route, Navigate } from 'react-router';
import { ThemeProvider, CssBaseline, CircularProgress, Box } from '@mui/material';
import { buildTheme } from './theme';
import ErrorBoundary from './components/ErrorBoundary';
import Layout from './components/Layout';

const LoginPage = lazy(() => import('./pages/LoginPage'));
const RegisterPage = lazy(() => import('./pages/RegisterPage'));
const DashboardPage = lazy(() => import('./pages/DashboardPage'));
const TasksPage = lazy(() => import('./pages/TasksPage'));
const ProjectsPage = lazy(() => import('./pages/ProjectsPage'));
const NotFoundPage = lazy(() => import('./pages/NotFoundPage'));

const PageLoader = () => (
  <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '60vh' }}>
    <CircularProgress />
  </Box>
);

const isAuthenticated = () => !!localStorage.getItem('accessToken');

function ProtectedRoute({ children }) {
  return isAuthenticated() ? children : <Navigate to="/login" replace />;
}

function PublicRoute({ children }) {
  return isAuthenticated() ? <Navigate to="/dashboard" replace /> : children;
}

function App() {
  const [darkMode, setDarkMode] = useState(
    () => window.matchMedia?.('(prefers-color-scheme: dark)')?.matches ?? false
  );

  const theme = useMemo(() => buildTheme(darkMode), [darkMode]);

  return (
    <ErrorBoundary>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Suspense fallback={<PageLoader />}>
          <Routes>
            <Route
              path="/login"
              element={<PublicRoute><LoginPage /></PublicRoute>}
            />
            <Route
              path="/register"
              element={<PublicRoute><RegisterPage /></PublicRoute>}
            />
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <Layout darkMode={darkMode} onToggleDarkMode={() => setDarkMode((d) => !d)} />
                </ProtectedRoute>
              }
            >
              <Route index element={<Navigate to="/dashboard" replace />} />
              <Route path="dashboard" element={<DashboardPage />} />
              <Route path="tasks" element={<TasksPage />} />
              <Route path="projects" element={<ProjectsPage />} />
            </Route>
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </Suspense>
      </ThemeProvider>
    </ErrorBoundary>
  );
}

export default App;
