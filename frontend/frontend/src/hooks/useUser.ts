import { useDispatch, useSelector } from 'react-redux';
import { setUser, clearUser } from '../store/userSlice';
import { RootState } from '../store';

export const useUser = () => {
  const dispatch = useDispatch();
  const { login, isAuthenticated } = useSelector((state: RootState) => state.user);

  const loginUser = (username: string) => {
    dispatch(setUser(username));
  };

  const logoutUser = () => {
    dispatch(clearUser());
  };

  return {
    login,
    isAuthenticated,
    loginUser,
    logoutUser,
  };
};