import { combineReducers } from '@reduxjs/toolkit';
import functionsSlice from './functionsSlice';

const createReducer = asyncReducers =>
	combineReducers({
		functionsSlice,
		...asyncReducers
	});

export default createReducer;
