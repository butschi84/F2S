import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios'

async function check(url) {
	const reqUrl = `${url}/docs/`
	try{
		let {data} = await axios.get(reqUrl);
		return true
	}catch(ex) {
		return false
	}
}

export const checkConnectivity = () => async (dispatch, getState) => {
	
	const { apiURL } = getState().connectivitySlice;
	const connected = await check(apiURL);
	return dispatch(setConnectivity(connected));
  };

const connectivitySlice = createSlice({
	name: 'connectivity/manage',
	initialState: {
        ApiConnectionEstablished: false,
		apiURL: "http://localhost:57336"
    },
	reducers: {
		setConnectivity: (state, action) => {
			return {...state, ApiConnectionEstablished: action.payload};
		}
	},
	extraReducers: {}
});

export const { setConnectivity } = connectivitySlice.actions;

export default connectivitySlice.reducer;