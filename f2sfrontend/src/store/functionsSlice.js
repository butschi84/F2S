import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';

import {
	getFunctions as serviceGetFunctions,
	createFunction as serviceCreateFunction
} from '../services/functions'

// create a new f2s function
export const createNewF2SFunction = (func) => async dispatch => {
	await serviceCreateFunction(func)
	dispatch(getAllFunctions())
}

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