import axios from 'axios';

// Add a 401 response interceptor
axios.interceptors.response.use(function (response) {
    return response;
}, function (error) {
    if (error.response && 401 === error.response.status) {
        console.log("got 401 response. => we are unauthorized")
        localStorage.removeItem("token");
        window.location = "/"; 
        return Promise.reject("unauthorized");
    } else {
        return Promise.reject(error);
    }
});

function fetchAuthHeader() {
    const authType = localStorage.getItem("authtype");

    switch(authType) {
        case "token":
            const token = localStorage.getItem("token");
            return `Bearer ${token}`;
        case "basic":
            const username = localStorage.getItem("username");
            const password = localStorage.getItem("password");
            const basic = btoa(`${username}:${password}`); // Encode username and password in Base64
            return `Basic ${basic}`;
    }

    return null;
}

export function get(url) {
    const apiurl = localStorage.getItem("apiurl")

    // fetch headers
    const authHeader = fetchAuthHeader();
    const headers = authHeader ? { Authorization: authHeader } : {};

    return new Promise(async (resolve, reject) => {
        axios.get(`${apiurl}${url}`, { headers }).then(result => {
            return resolve(result.data);
        }).catch(ex => {
            reject(ex)
        })
    });
}

export function del(url, data=null) {
    const apiurl = localStorage.getItem("apiurl")

    // fetch headers
    const authHeader = fetchAuthHeader();
    const headers = authHeader ? { Authorization: authHeader } : {};

    return new Promise(async (resolve, reject) => {
        try{
            let {result} = await axios.delete(`${apiurl}${url}`, { headers }, {
                data: data ? data : null
            });
            resolve(result);
        }catch(ex) {
            reject(ex)
        }
    });
}

export function post(url, postData) {
    const apiurl = localStorage.getItem("apiurl")

    // fetch headers
    const authHeader = fetchAuthHeader();
    const headers = authHeader ? { Authorization: authHeader } : {};

    return new Promise(async (resolve, reject) => {
        try{
            let {data} = await axios.post(`${apiurl}${url}`, postData, { headers });
            resolve(data);
        }catch(ex) {
            reject(ex)
        }
    });
}

export function put(url, postData) {
    const apiurl = localStorage.getItem("apiurl")

    // fetch headers
    const authHeader = fetchAuthHeader();
    const headers = authHeader ? { Authorization: authHeader } : {};

    return new Promise(async (resolve, reject) => {
        try{
            let {data} = await axios.put(`${apiurl}${url}`, postData, { headers });
            resolve(data);
        }catch(ex) {
            reject(ex)
        }
    });
}
