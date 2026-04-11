import React, { useEffect, useState } from 'react';
import {
  Box,
  Typography,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Fab,
  Grid2 as Grid,
  CircularProgress,
  InputAdornment,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import SearchIcon from '@mui/icons-material/Search';
import { useDispatch, useSelector } from 'react-redux';
import { fetchTasks, createTask, updateTask } from '../store/taskSlice';
import { fetchProjects } from '../store/projectSlice';
import { TASK_STATUSES, TASK_PRIORITIES } from '../constants';
import TaskCard from '../components/TaskCard';
import TaskDialog from '../components/TaskDialog';

const STATUS_FILTER_OPTIONS = ['', ...TASK_STATUSES];
const PRIORITY_FILTER_OPTIONS = ['', ...TASK_PRIORITIES];

export default function TasksPage() {
  const dispatch = useDispatch();
  const { tasks, loading } = useSelector((state) => state.task);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingTask, setEditingTask] = useState(null);
  const [filters, setFilters] = useState({
    status: '',
    priority: '',
    search: '',
  });

  useEffect(() => {
    dispatch(fetchTasks());
    dispatch(fetchProjects());
  }, [dispatch]);

  const filteredTasks = tasks.filter((task) => {
    if (filters.status && task.status !== filters.status) return false;
    if (filters.priority && task.priority !== filters.priority) return false;
    if (filters.search) {
      const q = filters.search.toLowerCase();
      const match =
        (task.title || '').toLowerCase().includes(q) ||
        (task.description || '').toLowerCase().includes(q);
      if (!match) return false;
    }
    return true;
  });

  const handleCreate = () => {
    setEditingTask(null);
    setDialogOpen(true);
  };

  const handleEdit = (task) => {
    setEditingTask(task);
    setDialogOpen(true);
  };

  const handleSubmit = (payload) => {
    if (payload.id) {
      dispatch(updateTask({ id: payload.id, ...payload }));
    } else {
      dispatch(createTask(payload));
    }
    setDialogOpen(false);
    setEditingTask(null);
  };

  return (
    <Box>
      <Typography variant="h4" fontWeight={700} gutterBottom>
        Задачи
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
        Управление задачами
      </Typography>

      {/* Панель фильтров */}
      <Box
        sx={{
          display: 'flex',
          flexWrap: 'wrap',
          gap: 2,
          mb: 3,
          p: 2,
          borderRadius: 2,
          bgcolor: 'background.paper',
          boxShadow: 1,
        }}
      >
        <TextField
          placeholder="Поиск..."
          value={filters.search}
          onChange={(e) => setFilters((f) => ({ ...f, search: e.target.value }))}
          size="small"
          sx={{ minWidth: 200 }}
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon color="action" />
                </InputAdornment>
              ),
            },
          }}
        />
        <FormControl size="small" sx={{ minWidth: 140 }}>
          <InputLabel>Статус</InputLabel>
          <Select
            value={filters.status}
            onChange={(e) => setFilters((f) => ({ ...f, status: e.target.value }))}
            label="Статус"
          >
            {STATUS_FILTER_OPTIONS.map((s) => (
              <MenuItem key={s || 'all'} value={s}>
                {s || 'Все'}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <FormControl size="small" sx={{ minWidth: 140 }}>
          <InputLabel>Приоритет</InputLabel>
          <Select
            value={filters.priority}
            onChange={(e) => setFilters((f) => ({ ...f, priority: e.target.value }))}
            label="Приоритет"
          >
            {PRIORITY_FILTER_OPTIONS.map((p) => (
              <MenuItem key={p || 'all'} value={p}>
                {p || 'Все'}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>

      {loading ? (
        <Box display="flex" justifyContent="center" py={4}>
          <CircularProgress />
        </Box>
      ) : (
        <Grid container spacing={2}>
          {filteredTasks.map((task) => (
            <Grid size={{ xs: 12, sm: 6, md: 4 }} key={task.id}>
              <TaskCard task={task} onClick={() => handleEdit(task)} />
            </Grid>
          ))}
        </Grid>
      )}

      {filteredTasks.length === 0 && !loading && (
        <Box sx={{ textAlign: 'center', py: 6 }}>
          <Typography color="text.secondary">Нет задач. Создайте первую!</Typography>
        </Box>
      )}

      <Fab
        color="primary"
        aria-label="добавить задачу"
        sx={{ position: 'fixed', bottom: 24, right: 24 }}
        onClick={handleCreate}
      >
        <AddIcon />
      </Fab>

      <TaskDialog
        open={dialogOpen}
        onClose={() => setDialogOpen(false)}
        onSubmit={handleSubmit}
        task={editingTask}
      />
    </Box>
  );
}
