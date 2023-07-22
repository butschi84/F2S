import axios from 'axios';

// Add a 401 response interceptor
axios.interceptors.response.use(function (response) {
    return response;
}, function (error) {
    if (error.response && 401 === error.response.status) {
        localStorage.removeItem("token");
        window.location = "/login"; 
    } else {
        return Promise.reject(error);
    }
});

export function get(url) {
    const apiurl = localStorage.getItem("apiurl")
    return new Promise(async (resolve, reject) => {
        try{
            let {data} = await axios.get(`${apiurl}${url}`);
            resolve(data);
        }catch(ex) {
            reject(ex)
        }
    });
}

export function del(url, data=null) {
    const apiurl = localStorage.getItem("apiurl")
    return new Promise(async (resolve, reject) => {
        try{
            let {result} = await axios.delete(`${apiurl}${url}`, {
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
    return new Promise(async (resolve, reject) => {
        try{
            let {data} = await axios.post(`${apiurl}${url}`, postData);
            resolve(data);
        }catch(ex) {
            reject(ex)
        }
    });
}

export function put(url, postData) {
    const apiurl = localStorage.getItem("apiurl")
    return new Promise(async (resolve, reject) => {
        try{
            let {data} = await axios.put(`${apiurl}${url}`, postData);
            resolve(data);
        }catch(ex) {
            reject(ex)
        }
    });
}
