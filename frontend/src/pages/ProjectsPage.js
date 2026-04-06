// Страница проектов: сетка карточек, диалог создания/редактирования, счётчик задач
import React, { useEffect, useState } from 'react';
import {
  Box,
  Typography,
  Grid2 as Grid,
  Card,
  CardContent,
  CardActions,
  Button,
  IconButton,
  Chip,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import FolderIcon from '@mui/icons-material/Folder';
import { useDispatch, useSelector } from 'react-redux';
import {
  fetchProjects,
  createProject,
  updateProject,
  deleteProject,
} from '../store/projectSlice';

// Компонент выбора цвета
const ColorPicker = ({ value, onChange }) => {
  const colors = [
    '#1976d2',
    '#42a5f5',
    '#26a69a',
    '#66bb6a',
    '#ffa726',
    '#ef5350',
    '#ab47bc',
    '#7e57c2',
  ];
  return (
    <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
      {colors.map((c) => (
        <Box
          key={c}
          onClick={() => onChange(c)}
          sx={{
            width: 32,
            height: 32,
            borderRadius: '50%',
            bgcolor: c,
            cursor: 'pointer',
            border: value === c ? 3 : 1,
            borderColor: value === c ? 'primary.main' : 'divider',
            '&:hover': { transform: 'scale(1.1)' },
          }}
        />
      ))}
    </Box>
  );
};

export default function ProjectsPage() {
  const dispatch = useDispatch();
  const { projects, loading } = useSelector((state) => state.project);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingProject, setEditingProject] = useState(null);
  const [form, setForm] = useState({ name: '', description: '', color: '#1976d2' });

  useEffect(() => {
    dispatch(fetchProjects());
  }, [dispatch]);

  const handleOpenCreate = () => {
    setEditingProject(null);
    setForm({ name: '', description: '', color: '#1976d2' });
    setDialogOpen(true);
  };

  const handleOpenEdit = (project) => {
    setEditingProject(project);
    setForm({
      name: project.name || '',
      description: project.description || '',
      color: project.color || '#1976d2',
    });
    setDialogOpen(true);
  };

  const handleSubmit = () => {
    if (editingProject) {
      dispatch(updateProject({ id: editingProject.id, ...form }));
    } else {
      dispatch(createProject(form));
    }
    setDialogOpen(false);
  };

  const handleDelete = (id) => {
    if (window.confirm('Удалить проект?')) {
      dispatch(deleteProject(id));
    }
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Box>
          <Typography variant="h4" fontWeight={700} gutterBottom>
            Проекты
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Организация задач по проектам
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={handleOpenCreate}
          sx={{ borderRadius: 2 }}
        >
          Новый проект
        </Button>
      </Box>

      {loading ? (
        <Box display="flex" justifyContent="center" py={4}>
          <CircularProgress />
        </Box>
      ) : (
        <Grid container spacing={3}>
          {projects.map((project) => (
            <Grid size={{ xs: 12, sm: 6, md: 4 }} key={project.id}>
              <Card
                sx={{
                  height: '100%',
                  display: 'flex',
                  flexDirection: 'column',
                  borderLeft: 4,
                  borderColor: project.color || 'primary.main',
                  transition: 'transform 0.2s, box-shadow 0.2s',
                  '&:hover': {
                    transform: 'translateY(-2px)',
                    boxShadow: 4,
                  },
                }}
              >
                <CardContent sx={{ flexGrow: 1 }}>
                  <Box sx={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', mb: 1 }}>
                    <FolderIcon sx={{ color: project.color || 'primary.main', fontSize: 40 }} />
                    <Box>
                      <IconButton size="small" onClick={() => handleOpenEdit(project)}>
                        <EditIcon fontSize="small" />
                      </IconButton>
                      <IconButton size="small" color="error" onClick={() => handleDelete(project.id)}>
                        <DeleteIcon fontSize="small" />
                      </IconButton>
                    </Box>
                  </Box>
                  <Typography variant="h6" fontWeight={600} gutterBottom>
                    {project.name}
                  </Typography>
                  {project.description && (
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                      {project.description}
                    </Typography>
                  )}
                  <Chip
                    label={`${project.taskCount ?? 0} задач`}
                    size="small"
                    sx={{ mt: 1 }}
                  />
                </CardContent>
                <CardActions>
                  <Button size="small" onClick={() => handleOpenEdit(project)}>
                    Редактировать
                  </Button>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      {projects.length === 0 && !loading && (
        <Box sx={{ textAlign: 'center', py: 6 }}>
          <Typography color="text.secondary">Нет проектов. Создайте первый!</Typography>
        </Box>
      )}

      <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>{editingProject ? 'Редактировать проект' : 'Новый проект'}</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
            <TextField
              label="Название"
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
              fullWidth
              required
            />
            <TextField
              label="Описание"
              value={form.description}
              onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
              multiline
              rows={3}
              fullWidth
            />
            <Box>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                Цвет
              </Typography>
              <ColorPicker value={form.color} onChange={(c) => setForm((f) => ({ ...f, color: c }))} />
            </Box>
          </Box>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button onClick={() => setDialogOpen(false)}>Отмена</Button>
          <Button variant="contained" onClick={handleSubmit} disabled={!form.name.trim()}>
            {editingProject ? 'Сохранить' : 'Создать'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
