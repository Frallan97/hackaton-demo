import { useDispatch, useSelector, TypedUseSelectorHook } from 'react-redux';
import type { RootState, AppDispatch } from './index';
import { setUser, setTokens } from './slices/authSlice';

export const useAppDispatch: () => AppDispatch = useDispatch;
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;

// Custom hook for Google OAuth login
export const useGoogleOAuth = () => {
  const dispatch = useAppDispatch();
  
  const handleGoogleLoginSuccess = (response: any) => {
    if (response.user && response.access_token) {
      dispatch(setUser(response.user));
      dispatch(setTokens({
        accessToken: response.access_token,
        refreshToken: response.refresh_token,
      }));
    }
  };
  
  return { handleGoogleLoginSuccess };
}; 