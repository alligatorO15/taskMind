// Слайс аутентификации: пользователь, токены, логин, регистрация
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import api, { setRefreshTokenFn } from '../services/api';

// Регистрация пользователя
export const register = createAsyncThunk(
  'auth/register',
  async ({ username, email, password }, { rejectWithValue, dispatch }) => {
    try {
      const { data } = await api.post('/auth/register', { username, email, password });
      if (data.access_token) {
        localStorage.setItem('accessToken', data.access_token);
        localStorage.setItem('refreshToken', data.refresh_token || '');
      }
      await dispatch(getProfile()).unwrap();
      return data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.error || 'Ошибка регистрации');
    }
  }
);

// Вход в систему
export const login = createAsyncThunk(
  'auth/login',
  async ({ email, password }, { rejectWithValue, dispatch }) => {
    try {
      const { data } = await api.post('/auth/login', { email, password });
      if (data.access_token) {
        localStorage.setItem('accessToken', data.access_token);
        localStorage.setItem('refreshToken', data.refresh_token || '');
      }
      await dispatch(getProfile()).unwrap();
      return data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.error || 'Ошибка входа');
    }
  }
);

// Обновление токена
export const refreshToken = createAsyncThunk(
  'auth/refreshToken',
  async (_, { rejectWithValue }) => {
    try {
      const refresh = localStorage.getItem('refreshToken');
      if (!refresh) throw new Error('Нет refresh токена');
      const { data } = await api.post('/auth/refresh', { refresh_token: refresh });
      localStorage.setItem('accessToken', data.access_token);
      if (data.refresh_token) {
        localStorage.setItem('refreshToken', data.refresh_token);
      }
      return data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.error || 'Ошибка обновления токена');
    }
  }
);

// Получение профиля пользователя
export const getProfile = createAsyncThunk(
  'auth/getProfile',
  async (_, { rejectWithValue }) => {
    try {
      const { data } = await api.get('/profile');
      return data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.error || 'Ошибка загрузки профиля');
    }
  }
);

// Выход из системы
export const logout = createAsyncThunk('auth/logout', async () => {
  localStorage.removeItem('accessToken');
  localStorage.removeItem('refreshToken');
});

const authSlice = createSlice({
  name: 'auth',
  initialState: {
    user: null,
    tokens: { accessToken: null, refreshToken: null },
    loading: false,
    error: null,
  },
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    setUser: (state, action) => {
      state.user = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      // register
      .addCase(register.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(register.fulfilled, (state, action) => {
        state.loading = false;
        state.tokens = {
          accessToken: action.payload.access_token,
          refreshToken: action.payload.refresh_token,
        };
        state.error = null;
      })
      .addCase(register.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      })
      // login
      .addCase(login.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(login.fulfilled, (state, action) => {
        state.loading = false;
        state.tokens = {
          accessToken: action.payload.access_token,
          refreshToken: action.payload.refresh_token,
        };
        state.error = null;
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      })
      // refreshToken
      .addCase(refreshToken.fulfilled, (state, action) => {
        state.tokens.accessToken = action.payload.access_token;
        if (action.payload.refresh_token) {
          state.tokens.refreshToken = action.payload.refresh_token;
        }
      })
      .addCase(refreshToken.rejected, (state) => {
        state.user = null;
        state.tokens = { accessToken: null, refreshToken: null };
      })
      // getProfile
      .addCase(getProfile.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(getProfile.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload;
      })
      .addCase(getProfile.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload;
      })
      // logout
      .addCase(logout.fulfilled, (state) => {
        state.user = null;
        state.tokens = { accessToken: null, refreshToken: null };
        state.error = null;
      });
  },
});

// Регистрируем функцию обновления токена для api.js (вызывается из store после создания)
export const setupRefreshToken = (dispatch) => {
  setRefreshTokenFn(() => dispatch(refreshToken()));
};

export const { clearError, setUser } = authSlice.actions;
export default authSlice.reducer;
