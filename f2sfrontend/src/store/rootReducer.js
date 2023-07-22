import { combineReducers } from '@reduxjs/toolkit';
import functionsSlice from './functionsSlice';
import connectivitySlice from './connectivitySlice';
import configSlice from './configSlice';

const createReducer = asyncReducers =>
	combineReducers({
		functionsSlice,
		configSlice,
		connectivitySlice,
		...asyncReducers
	});

export default createReducer;
