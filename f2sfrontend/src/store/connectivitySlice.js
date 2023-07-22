import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios'

async function check(url) {
	const reqUrl = `${url}/health`
	try{
		const response = await axios.get(reqUrl);
		const status = response.data.status
		return status == "ok"
	}catch(ex) {
		return false
	}
}
export const checkConnectivity = () => async (dispatch, getState) => {
	const { apiURL } = getState().connectivitySlice;
	
	const connected = await check(apiURL);
	if(connected){
		localStorage.setItem("apiurl", apiURL)
	}else{
		localStorage.removeItem("apiurl")
	}
	return dispatch(setConnectivity(connected));
};
export const setBackendURL = (url) => async (dispatch, getState) => {
	await dispatch(setURL(url))
	dispatch(checkConnectivity())
};

function getInitialApiURL() {
	const localStoreApiUrl = localStorage.getItem("apiurl")
	if(localStoreApiUrl != "") return localStoreApiUrl
	return window.location.href
}

const connectivitySlice = createSlice({
	name: 'connectivity/manage',
	initialState: {
        ApiConnectionEstablished: false,
		apiURL: getInitialApiURL()
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