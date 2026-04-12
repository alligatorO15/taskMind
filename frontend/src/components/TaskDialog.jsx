import React, { useEffect, useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Box,
  Chip,
  InputAdornment,
} from '@mui/material';
import { useSelector } from 'react-redux';
import { TASK_STATUSES, TASK_PRIORITIES, STATUS_LABELS, PRIORITY_LABELS } from '../constants';

export default function TaskDialog({ open, onClose, onSubmit, task = null }) {
  const { projects } = useSelector((state) => state.project);
  const [form, setForm] = useState({
    title: '',
    description: '',
    status: 'new',
    priority: 'medium',
    projectId: '',
    tags: [],
    deadline: '',
    reminder: '',
    tagInput: '',
  });

  useEffect(() => {
    if (task) {
      setForm({
        title: task.title || '',
        description: task.description || '',
        status: task.status || 'new',
        priority: task.priority || 'medium',
        projectId: task.project_id || task.projectId || task.project?.id || '',
        tags: Array.isArray(task.tags) ? task.tags.map((t) => (typeof t === 'string' ? t : t.name)) : [],
        deadline: task.deadline ? task.deadline.slice(0, 16) : '',
        reminder: task.reminder ? task.reminder.slice(0, 16) : '',
        tagInput: '',
      });
    } else {
      setForm({
        title: '',
        description: '',
        status: 'new',
        priority: 'medium',
        projectId: '',
        tags: [],
        deadline: '',
        reminder: '',
        tagInput: '',
      });
    }
  }, [task, open]);

  const handleChange = (field) => (e) => {
    setForm((f) => ({ ...f, [field]: e.target.value }));
  };

  const handleAddTag = () => {
    const tag = form.tagInput.trim();
    if (tag && !form.tags.includes(tag)) {
      setForm((f) => ({ ...f, tags: [...f.tags, tag], tagInput: '' }));
    }
  };

  const handleRemoveTag = (tag) => {
    setForm((f) => ({ ...f, tags: f.tags.filter((t) => t !== tag) }));
  };

  const toRFC3339 = (v) => (v ? new Date(v).toISOString() : undefined);

  const handleSubmit = () => {
    const payload = {
      title: form.title,
      description: form.description || undefined,
      status: form.status,
      priority: form.priority,
      project_id: form.projectId || undefined,
      tags: form.tags,
      deadline: toRFC3339(form.deadline),
      reminder_before: form.reminder || undefined,
    };
    if (task) payload.id = task.id;
    onSubmit(payload);
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth slotProps={{ paper: { sx: { borderRadius: 2 } } }}>
      <DialogTitle>{task ? 'Редактировать задачу' : 'Новая задача'}</DialogTitle>
      <DialogContent>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
          <TextField
            label="Название"
            value={form.title}
            onChange={handleChange('title')}
            required
            fullWidth
          />
          <TextField
            label="Описание"
            value={form.description}
            onChange={handleChange('description')}
            multiline
            rows={3}
            fullWidth
          />
          <Box sx={{ display: 'flex', gap: 2 }}>
            <FormControl fullWidth>
              <InputLabel>Статус</InputLabel>
              <Select value={form.status} onChange={handleChange('status')} label="Статус">
                {TASK_STATUSES.map((s) => (
                  <MenuItem key={s} value={s}>
                    {STATUS_LABELS[s] || s}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            <FormControl fullWidth>
              <InputLabel>Приоритет</InputLabel>
              <Select value={form.priority} onChange={handleChange('priority')} label="Приоритет">
                {TASK_PRIORITIES.map((p) => (
                  <MenuItem key={p} value={p}>
                    {PRIORITY_LABELS[p] || p}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>
          <FormControl fullWidth>
            <InputLabel>Проект</InputLabel>
            <Select value={form.projectId} onChange={handleChange('projectId')} label="Проект">
              <MenuItem value="">Без проекта</MenuItem>
              {projects.map((p) => (
                <MenuItem key={p.id} value={p.id}>
                  {p.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <TextField
            label="Теги"
            value={form.tagInput}
            onChange={handleChange('tagInput')}
            onKeyDown={(e) => e.key === 'Enter' && (e.preventDefault(), handleAddTag())}
            slotProps={{
              input: {
                endAdornment: (
                  <InputAdornment position="end">
                    <Button size="small" onClick={handleAddTag}>
                      Добавить
                    </Button>
                  </InputAdornment>
                ),
              },
            }}
          />
          {form.tags.length > 0 && (
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
              {form.tags.map((tag) => (
                <Chip key={tag} label={tag} onDelete={() => handleRemoveTag(tag)} size="small" />
              ))}
            </Box>
          )}
          <TextField
            label="Дедлайн"
            type="datetime-local"
            value={form.deadline}
            onChange={handleChange('deadline')}
            slotProps={{ inputLabel: { shrink: true } }}
            fullWidth
          />
          <TextField
            label="Напоминание"
            type="datetime-local"
            value={form.reminder}
            onChange={handleChange('reminder')}
            slotProps={{ inputLabel: { shrink: true } }}
            fullWidth
          />
        </Box>
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2 }}>
        <Button onClick={onClose}>Отмена</Button>
        <Button variant="contained" onClick={handleSubmit} disabled={!form.title.trim()}>
          {task ? 'Сохранить' : 'Создать'}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
