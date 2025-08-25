import { useAppDispatch } from '../../../shared/hooks/redux';
import { setUser, setTokens } from '../store/authSlice';

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