import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface Notification {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  message: string;
  duration?: number;
}

interface UIState {
  notifications: Notification[];
  globalLoading: boolean;
  sidebarOpen: boolean;
  activeTab: string;
}

const initialState: UIState = {
  notifications: [],
  globalLoading: false,
  sidebarOpen: false,
  activeTab: 'home',
};

const uiSlice = createSlice({
  name: 'ui',
  initialState,
  reducers: {
    setGlobalLoading: (state, action: PayloadAction<boolean>) => {
      state.globalLoading = action.payload;
    },
    
    setSidebarOpen: (state, action: PayloadAction<boolean>) => {
      state.sidebarOpen = action.payload;
    },
    
    setActiveTab: (state, action: PayloadAction<string>) => {
      state.activeTab = action.payload;
    },
    
    addNotification: (state, action: PayloadAction<Omit<Notification, 'id'>>) => {
      const id = Date.now().toString();
      const notification: Notification = {
        id,
        ...action.payload,
        duration: action.payload.duration || 5000,
      };
      state.notifications.push(notification);
    },
    
    removeNotification: (state, action: PayloadAction<string>) => {
      state.notifications = state.notifications.filter(
        notification => notification.id !== action.payload
      );
    },
    
    clearNotifications: (state) => {
      state.notifications = [];
    },
    
    showSuccess: (state, action: PayloadAction<string>) => {
      state.notifications.push({
        id: Date.now().toString(),
        type: 'success',
        message: action.payload,
        duration: 5000,
      });
    },
    
    showError: (state, action: PayloadAction<string>) => {
      state.notifications.push({
        id: Date.now().toString(),
        type: 'error',
        message: action.payload,
        duration: 8000,
      });
    },
    
    showWarning: (state, action: PayloadAction<string>) => {
      state.notifications.push({
        id: Date.now().toString(),
        type: 'warning',
        message: action.payload,
        duration: 6000,
      });
    },
    
    showInfo: (state, action: PayloadAction<string>) => {
      state.notifications.push({
        id: Date.now().toString(),
        type: 'info',
        message: action.payload,
        duration: 4000,
      });
    },
  },
});

export const {
  setGlobalLoading,
  setSidebarOpen,
  setActiveTab,
  addNotification,
  removeNotification,
  clearNotifications,
  showSuccess,
  showError,
  showWarning,
  showInfo,
} = uiSlice.actions;

export default uiSlice.reducer; 