// Слайс проектов: список проектов, CRUD-операции
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import api from '../services/api';

// Загрузка списка проектов
export const fetchProjects = createAsyncThunk(
  'project/fetchProjects',
  async (_, { rejectWithValue }) => {
    try {
      const { data } = await api.get('/projects');
      return data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.message || 'Ошибка загрузки проектов');
    }
  }
);

// Создание проекта
export const createProject = createAsyncThunk(
  'project/createProject',
  async (projectData, { rejectWithValue }) => {
    try {
      const { data } = await api.post('/projects', projectData);
      return data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.message || 'Ошибка создания проекта');
    }
  }
);

// Обновление проекта
export const updateProject = createAsyncThunk(
  'project/updateProject',
  async ({ id, ...projectData }, { rejectWithValue }) => {
    try {
      const { data } = await api.put(`/projects/${id}`, projectData);
      return data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.message || 'Ошибка обновления проекта');
    }
  }
);

// Удаление проекта
export const deleteProject = createAsyncThunk(
  'project/deleteProject',
  async (id, { rejectWithValue }) => {
    try {
      await api.delete(`/projects/${id}`);
      return id;
    } catch (err) {
      return rejectWithValue(err.response?.data?.message || 'Ошибка удаления проекта');
    }
  }
);

const projectSlice = createSlice({
  name: 'project',
  initialState: {
    projects: [],
    loading: false,
    error: null,
  },
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // fetchProjects
      .addCase(fetchProjects.fulfilled, (state, action) => {
        state.projects = Array.isArray(action.payload) ? action.payload : action.payload.projects || [];
      })
      // createProject
      .addCase(createProject.fulfilled, (state, action) => {
        state.projects.push(action.payload);
      })
      // updateProject
      .addCase(updateProject.fulfilled, (state, action) => {
        const idx = state.projects.findIndex((p) => p.id === action.payload.id);
        if (idx >= 0) state.projects[idx] = action.payload;
      })
      // deleteProject
      .addCase(deleteProject.fulfilled, (state, action) => {
        state.projects = state.projects.filter((p) => p.id !== action.payload);
      })
      .addMatcher(
        (action) => action.type.startsWith('project/') && action.type.endsWith('/pending'),
        (state) => {
          state.loading = true;
          state.error = null;
        },
      )
      .addMatcher(
        (action) => action.type.startsWith('project/') && action.type.endsWith('/fulfilled'),
        (state) => {
          state.loading = false;
        },
      )
      .addMatcher(
        (action) => action.type.startsWith('project/') && action.type.endsWith('/rejected'),
        (state, action) => {
          state.loading = false;
          state.error = action.payload;
        },
      );
  },
});

export const { clearError } = projectSlice.actions;
export default projectSlice.reducer;
