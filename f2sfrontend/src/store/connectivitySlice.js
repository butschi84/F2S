import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios'

async function check(url) {
	const reqUrl = `${url}/docs/`
	try{
		const response = await axios.get(reqUrl);
		const statusCode = response.status;
		console.log("status", statusCode, reqUrl)
		return statusCode == 200
	}catch(ex) {
		return false
	}
}
export const checkConnectivity = () => async (dispatch, getState) => {
	const { apiURL } = getState().connectivitySlice;
	const connected = await check(apiURL);
	return dispatch(setConnectivity(connected));
};
export const setBackendURL = (url) => async (dispatch, getState) => {
	console.log("set url", url)
	await dispatch(setURL(url))
	dispatch(checkConnectivity())
};

const connectivitySlice = createSlice({
	name: 'connectivity/manage',
	initialState: {
        ApiConnectionEstablished: false,
		apiURL: window.location.href
    },
	reducers: {
		setURL: (state, action) => {
			return {...state, apiURL: action.payload};
		},
		setConnectivity: (state, action) => {
			return {...state, ApiConnectionEstablished: action.payload};
		}
	},
	extraReducers: {}
});

export const { setConnectivity, setURL } = connectivitySlice.actions;

export default connectivitySlice.reducer;