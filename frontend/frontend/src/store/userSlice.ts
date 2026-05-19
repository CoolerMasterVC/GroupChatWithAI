import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface UserState {
  login: string | null;
  isAuthenticated: boolean;
}

const initialState: UserState = {
  login: null,
  isAuthenticated: false,
};

const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setUser: (state, action: PayloadAction<string>) => {
      state.login = action.payload;
      state.isAuthenticated = true;
    },
    clearUser: (state) => {
      state.login = null;
      state.isAuthenticated = false;
    },
  },
});

export const { setUser, clearUser } = userSlice.actions;
export default userSlice.reducer;