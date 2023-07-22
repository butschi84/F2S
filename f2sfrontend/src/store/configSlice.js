import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { getConfig as serviceGetConfig } from '../services/config'
import axios from 'axios'

export const getF2SConfig = () => async (dispatch, getState) => {
	const config = await serviceGetConfig()
    dispatch(setConfig(config))
};

const configSlice = createSlice({
	name: 'connectivity/manage',
	initialState: {
        config: {}
    },
	reducers: {
		setConfig: (state, action) => {
			return {...state, config: action.payload};
		}
	},
	extraReducers: {}
});

export const { setConfig } = configSlice.actions;

export default configSlice.reducer;