// Дашборд: статистика задач, последние задачи
import React, { useEffect } from 'react';
import { useNavigate } from 'react-router';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Chip,
  Button,
  CircularProgress,
} from '@mui/material';
import TaskAltIcon from '@mui/icons-material/TaskAlt';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import WarningAmberIcon from '@mui/icons-material/WarningAmber';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import { useDispatch, useSelector } from 'react-redux';
import { fetchTasks } from '../store/taskSlice';

const StatCard = ({ title, value, icon, color }) => (
  <Card sx={{ height: '100%', borderLeft: 4, borderColor: color }}>
    <CardContent>
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
        <Typography color="text.secondary" variant="body2">
          {title}
        </Typography>
        {icon}
      </Box>
      <Typography variant="h4" fontWeight={700}>
        {value}
      </Typography>
    </CardContent>
  </Card>
);

export default function DashboardPage() {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { tasks, loading } = useSelector((state) => state.task);

  useEffect(() => {
    dispatch(fetchTasks());
  }, [dispatch]);

  const total = tasks.length;
  const done = tasks.filter((t) => t.status === 'done' || t.status === 'completed').length;
  const overdue = tasks.filter((t) => {
    if (!t.deadline) return false;
    return new Date(t.deadline) < new Date() && t.status !== 'done' && t.status !== 'completed';
  }).length;

  const recentTasks = tasks.slice(0, 5);

  return (
    <Box>
      <Typography variant="h4" fontWeight={700} gutterBottom>
        Дашборд
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
        Обзор ваших задач
      </Typography>

      {loading ? (
        <Box display="flex" justifyContent="center" py={4}>
          <CircularProgress />
        </Box>
      ) : (
        <>
          <Grid container spacing={3} sx={{ mb: 4 }}>
            <Grid size={{ xs: 12, sm: 4 }}>
              <StatCard
                title="Всего задач"
                value={total}
                icon={<TaskAltIcon color="primary" />}
                color="primary.main"
              />
            </Grid>
            <Grid size={{ xs: 12, sm: 4 }}>
              <StatCard
                title="Выполнено"
                value={done}
                icon={<CheckCircleIcon color="success" />}
                color="success.main"
              />
            </Grid>
            <Grid size={{ xs: 12, sm: 4 }}>
              <StatCard
                title="Просрочено"
                value={overdue}
                icon={<WarningAmberIcon color="warning" />}
                color="warning.main"
              />
            </Grid>
          </Grid>

          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h6" fontWeight={600}>
                  Последние задачи
                </Typography>
                <Button
                  endIcon={<ArrowForwardIcon />}
                  onClick={() => navigate('/tasks')}
                  size="small"
                >
                  Все задачи
                </Button>
              </Box>
              {recentTasks.length === 0 ? (
                <Typography color="text.secondary">Нет задач</Typography>
              ) : (
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                  {recentTasks.map((task) => (
                    <Box
                      key={task.id}
                      sx={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        p: 1.5,
                        borderRadius: 2,
                        bgcolor: 'action.hover',
                        cursor: 'pointer',
                        '&:hover': { bgcolor: 'action.selected' },
                      }}
                      onClick={() => navigate('/tasks')}
                    >
                      <Typography fontWeight={500}>{task.title}</Typography>
                      <Box sx={{ display: 'flex', gap: 0.5 }}>
                        <Chip
                          label={task.status || 'new'}
                          size="small"
                          color={task.status === 'done' ? 'success' : 'default'}
                        />
                        {task.priority && (
                          <Chip label={task.priority} size="small" variant="outlined" />
                        )}
                      </Box>
                    </Box>
                  ))}
                </Box>
              )}
            </CardContent>
          </Card>
        </>
      )}
    </Box>
  );
}
