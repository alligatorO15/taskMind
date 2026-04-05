// Слайс задач: список, текущая задача, CRUD-операции
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import api from '../services/api';

// Загрузка списка задач
export const fetchTasks = createAsyncThunk(
  'task/fetchTasks',
  async (params = {}, { rejectWithValue }) => {
    try {
      const { data } = await api.get('/tasks', { params });
      return data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.message || 'Ошибка загрузки задач');
    }
  }
);

// Загрузка одной задачи по ID
export const fetchTask = createAsyncThunk(
  'task/fetchTask',
  async (id, { rejectWithValue }) => {
    try {
      const { data } = await api.get(`/tasks/${id}`);
      return data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.message || 'Ошибка загрузки задачи');
    }
  }
);

// Создание задачи
export const createTask = createAsyncThunk(
  'task/createTask',
  async (taskData, { rejectWithValue }) => {
    try {
      const { data } = await api.post('/tasks', taskData);
      return data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.message || 'Ошибка создания задачи');
    }
  }
);

// Обновление задачи
export const updateTask = createAsyncThunk(
  'task/updateTask',
  async ({ id, ...taskData }, { rejectWithValue }) => {
    try {
      const { data } = await api.patch(`/tasks/${id}`, taskData);
      return data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.message || 'Ошибка обновления задачи');
    }
  }
);

// Удаление задачи
export const deleteTask = createAsyncThunk(
  'task/deleteTask',
  async (id, { rejectWithValue }) => {
    try {
      await api.delete(`/tasks/${id}`);
      return id;
    } catch (err) {
      return rejectWithValue(err.response?.data?.message || 'Ошибка удаления задачи');
    }
  }
);

const taskSlice = createSlice({
  name: 'task',
  initialState: {
    tasks: [],
    currentTask: null,
    loading: false,
    error: null,
  },
  reducers: {
    clearCurrentTask: (state) => {
      state.currentTask = null;
    },
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // fetchTasks
      .addCase(fetchTasks.fulfilled, (state, action) => {
        state.tasks = Array.isArray(action.payload) ? action.payload : action.payload.tasks || [];
      })
      // fetchTask
      .addCase(fetchTask.fulfilled, (state, action) => {
        state.currentTask = action.payload;
      })
      // createTask
      .addCase(createTask.fulfilled, (state, action) => {
        state.tasks.unshift(action.payload);
      })
      // updateTask
      .addCase(updateTask.fulfilled, (state, action) => {
        const idx = state.tasks.findIndex((t) => t.id === action.payload.id);
        if (idx >= 0) state.tasks[idx] = action.payload;
        if (state.currentTask?.id === action.payload.id) {
          state.currentTask = action.payload;
        }
      })
      // deleteTask
      .addCase(deleteTask.fulfilled, (state, action) => {
        state.tasks = state.tasks.filter((t) => t.id !== action.payload);
        if (state.currentTask?.id === action.payload) {
          state.currentTask = null;
        }
      })
      // Любой pending → loading + сброс ошибки
      .addMatcher(
        (action) => action.type.startsWith('task/') && action.type.endsWith('/pending'),
        (state) => {
          state.loading = true;
          state.error = null;
        },
      )
      // Любой fulfilled → снять loading
      .addMatcher(
        (action) => action.type.startsWith('task/') && action.type.endsWith('/fulfilled'),
        (state) => {
          state.loading = false;
        },
      )
      // Любой rejected → снять loading + записать ошибку
      .addMatcher(
        (action) => action.type.startsWith('task/') && action.type.endsWith('/rejected'),
        (state, action) => {
          state.loading = false;
          state.error = action.payload;
        },
      );
  },
});

export const { clearCurrentTask, clearError } = taskSlice.actions;
export default taskSlice.reducer;
