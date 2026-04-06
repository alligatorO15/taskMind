import { Box, Typography, Button } from '@mui/material';
import { useNavigate } from 'react-router';
import SearchOffIcon from '@mui/icons-material/SearchOff';

export default function NotFoundPage() {
  const navigate = useNavigate();

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        gap: 2,
        p: 4,
        textAlign: 'center',
      }}
    >
      <SearchOffIcon sx={{ fontSize: 80, color: 'text.secondary' }} />
      <Typography variant="h3" fontWeight={700}>
        404
      </Typography>
      <Typography variant="h6" color="text.secondary">
        Страница не найдена
      </Typography>
      <Button variant="contained" onClick={() => navigate('/dashboard')} sx={{ mt: 2 }}>
        На главную
      </Button>
    </Box>
  );
}
