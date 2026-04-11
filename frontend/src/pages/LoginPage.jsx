// Страница входа: форма email/пароль
import React, { useEffect } from 'react';
import { Link, useNavigate } from 'react-router';
import {
  Box,
  Card,
  CardContent,
  TextField,
  Button,
  Typography,
  Alert,
  CircularProgress,
} from '@mui/material';
import { useDispatch, useSelector } from 'react-redux';
import { login, clearError } from '../store/authSlice';

export default function LoginPage() {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { loading, error, user } = useSelector((state) => state.auth);

  useEffect(() => {
    return () => dispatch(clearError());
  }, [dispatch]);

  useEffect(() => {
    if (user) navigate('/dashboard');
  }, [user, navigate]);

  const handleSubmit = (e) => {
    e.preventDefault();
    const form = e.target;
    dispatch(
      login({
        email: form.email.value,
        password: form.password.value,
      })
    );
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #1976d2 0%, #42a5f5 50%, #90caf9 100%)',
        p: 2,
      }}
    >
      <Card
        sx={{
          maxWidth: 420,
          width: '100%',
          boxShadow: '0 20px 60px rgba(0,0,0,0.2)',
          borderRadius: 3,
        }}
      >
        <CardContent sx={{ p: 4 }}>
          <Typography variant="h4" gutterBottom fontWeight={700} color="primary" textAlign="center">
            TaskMind
          </Typography>
          <Typography variant="body2" color="text.secondary" textAlign="center" sx={{ mb: 3 }}>
            Войдите в свой аккаунт
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }} onClose={() => dispatch(clearError())}>
              {error}
            </Alert>
          )}

          <form onSubmit={handleSubmit}>
            <TextField
              fullWidth
              label="Email"
              name="email"
              type="email"
              required
              autoComplete="email"
              sx={{ mb: 2 }}
            />
            <TextField
              fullWidth
              label="Пароль"
              name="password"
              type="password"
              required
              autoComplete="current-password"
              sx={{ mb: 3 }}
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              size="large"
              disabled={loading}
              sx={{ py: 1.5, borderRadius: 2 }}
            >
              {loading ? <CircularProgress size={24} color="inherit" /> : 'Войти'}
            </Button>
          </form>

          <Typography variant="body2" textAlign="center" sx={{ mt: 2 }}>
            Нет аккаунта?{' '}
            <Link to="/register" style={{ color: 'inherit', fontWeight: 600 }}>
              Зарегистрироваться
            </Link>
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );
}
