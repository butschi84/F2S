import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';

import {
	getFunctions as serviceGetFunctions
} from '../services/functions'


// get all f2s functions
export const getAllFunctions = () => async dispatch => {
	const allFunctions = await serviceGetFunctions();

    return dispatch(setFunctions(allFunctions));
};

const functionsSlice = createSlice({
	name: 'functions/manage',
	initialState: {
        functions: []
    },
	reducers: {
		setFunctions: (state, action) => {
			return {...state, functions: action.payload};
		}
	},
	extraReducers: {}
});

export const { setFunctions } = functionsSlice.actions;

export default functionsSlice.reducer;