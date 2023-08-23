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
	serviceGetFunctions().then(allFunctions => {
		return dispatch(setFunctions(allFunctions));
	}).catch(ex => {
		console.log("error when trying to get all functions:", ex);
		return dispatch(setFunctions([]));
	})
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