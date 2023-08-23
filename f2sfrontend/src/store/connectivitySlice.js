import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { 
	getAuthType as serviceGetAuthType,
	signin as serviceAuthSignin } from '../services/auth';
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

export const signinWithUsernamePassword = (username, password) => async (dispatch, getState) => {
	// save info in local store
	localStorage.setItem("username", username);
	localStorage.setItem("password", password);

	// try accessing the /auth/signin endpoint
	await serviceAuthSignin();

	// => interceptor will delete localstore token when response is 401 and redirect to '/'

	// authentication successful
	dispatch(setAuthenticated(true));
}

export const signinWithToken = (token) => async (dispatch, getState) => {
	// save token in local store
	localStorage.setItem("token", token);

	// try accessing the /auth/signin endpoint
	await serviceAuthSignin();

	// => interceptor will delete localstore token when response is 401 and redirect to '/'

	// authentication successful
	dispatch(setAuthenticated(true));
};

export const logout = () => async (dispatch, getState) => {
	localStorage.removeItem("token")
	localStorage.removeItem("apiurl")
	localStorage.removeItem("authtype")
	localStorage.removeItem("username", "password")
	dispatch(setAuthenticationType("none"))
	dispatch(setConnectivity(false))
	dispatch(setAuthenticated(false))
}

export const checkConnectivity = () => async (dispatch, getState) => {
	const { apiURL } = getState().connectivitySlice;
	
	const connected = await check(apiURL);
	if(connected){
		// save current apiurl in localstore
		localStorage.setItem("apiurl", apiURL)

		// check backend authentication type
		const authType = await serviceGetAuthType();
		localStorage.setItem("authtype", authType)
		dispatch(setAuthenticationType(authType));
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
		apiURL: getInitialApiURL(),
		authenticated: false,
		authenticationType: 'none'
    },
	reducers: {
		setAuthenticated: (state, action) => {
			return {...state, authenticated: action.payload};
		},
		setAuthenticationType: (state, action) => {
			return {...state, authenticationType: action.payload};
		},
		setURL: (state, action) => {
			return {...state, apiURL: action.payload};
		},
		setConnectivity: (state, action) => {
			return {...state, ApiConnectionEstablished: action.payload};
		}
	},
	extraReducers: {}
});

export const { setConnectivity, setURL, setAuthenticationType, setAuthenticated } = connectivitySlice.actions;

export default connectivitySlice.reducer;