import { combineReducers } from '@reduxjs/toolkit';
import functionsSlice from './functionsSlice';
import connectivitySlice from './connectivitySlice';

const createReducer = asyncReducers =>
	combineReducers({
		functionsSlice,
		connectivitySlice,
		...asyncReducers
	});

export default createReducer;
