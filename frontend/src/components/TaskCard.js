import React from 'react';
import { Card, CardContent, Typography, Box, Chip } from '@mui/material';
import EventIcon from '@mui/icons-material/Event';
import FolderIcon from '@mui/icons-material/Folder';
import { STATUS_COLORS, PRIORITY_COLORS } from '../constants';

export default function TaskCard({ task, onClick }) {
  const status = task.status || 'new';
  const priority = task.priority || 'medium';
  const deadline = task.deadline ? new Date(task.deadline).toLocaleDateString('ru') : null;
  const tags = Array.isArray(task.tags) ? task.tags : [];
  const projectName = task.project?.name || task.projectName || '';

  return (
    <Card
      onClick={onClick}
      sx={{
        cursor: 'pointer',
        transition: 'transform 0.2s, box-shadow 0.2s',
        '&:hover': {
          transform: 'translateY(-2px)',
          boxShadow: 4,
        },
      }}
    >
      <CardContent>
        <Typography variant="h6" fontWeight={600} gutterBottom noWrap>
          {task.title}
        </Typography>
        {task.description && (
          <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }} noWrap>
            {task.description}
          </Typography>
        )}
        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, alignItems: 'center', mt: 1 }}>
          <Chip label={status} size="small" color={STATUS_COLORS[status] || 'default'} />
          <Chip label={priority} size="small" color={PRIORITY_COLORS[priority] || 'default'} variant="outlined" />
          {deadline && (
            <Chip
              icon={<EventIcon sx={{ fontSize: 14 }} />}
              label={deadline}
              size="small"
              variant="outlined"
            />
          )}
          {projectName && (
            <Chip
              icon={<FolderIcon sx={{ fontSize: 14 }} />}
              label={projectName}
              size="small"
              variant="outlined"
            />
          )}
          {tags.map((tag) => (
            <Chip key={tag} label={typeof tag === 'string' ? tag : tag.name} size="small" />
          ))}
        </Box>
      </CardContent>
    </Card>
  );
}
