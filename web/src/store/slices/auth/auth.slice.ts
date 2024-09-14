import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import { AuthStateType } from "./auth.slice.types";

const initialState: AuthStateType = {
	user: null,
};
export const authSlice = createSlice({
	name: "auth",
	initialState,
	reducers: {
		login: (state: AuthStateType, action: PayloadAction<AuthStateType>) => {
			state.user = action.payload.user;
		},
		logout: (state: AuthStateType) => {
			state.user = null;
		},
		updateUser: (state: AuthStateType, action: PayloadAction<Partial<AuthStateType>>) => {
			if (state.user) {
				state.user = { ...state.user, ...action.payload.user };
			}
		},
	},
});

export const { updateUser, login, logout } = authSlice.actions;

export default authSlice.reducer;
