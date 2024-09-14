import { configureStore } from "@reduxjs/toolkit";
import { TypedUseSelectorHook, useDispatch as useAppDispatch, useSelector as useAppSelector } from "react-redux";
import { combineReducers } from "redux";
import { apiMiddleware } from "./api.middleware";
import { APIsReducers } from "./api.reducers";
import authSlice from "./slices/auth/auth.slice";

import modalSlice from "./slices/modal/modal.slice";

// Combine your slices with the RTK Query API service reducer
const rootReducer = combineReducers({
	auth: authSlice,
	modal: modalSlice,
	...APIsReducers,
});

// Configure your store to include the API service's middleware
export const store = configureStore({
	reducer: rootReducer,
	middleware: (getDefaultMiddleware) =>
		getDefaultMiddleware({
			serializableCheck: false,
		}).concat(apiMiddleware),
});

type Store = typeof store;
export type RootState = ReturnType<Store["getState"]>;
export type AppDispatch = Store["dispatch"];
export const { dispatch } = store;
export const useDispatch = () => useAppDispatch<AppDispatch>();
export const useSelector: TypedUseSelectorHook<RootState> = useAppSelector;

export default store;
