import axios from 'axios';

// Add a 401 response interceptor
axios.interceptors.response.use(function (response) {
    return response;
}, function (error) {
    if (error.response && 401 === error.response.status) {
        console.log("interceptor: token expired")
        localStorage.removeItem("token");
        window.location = "/login"; 
    } else {
        return Promise.reject(error);
    }
});

export function get(url) {
    return new Promise(async (resolve, reject) => {
        try{
            let {data} = await axios.get(url);
            resolve(data);
        }catch(ex) {
            reject(ex)
        }
    });
}

export function del(url, data=null) {
    return new Promise(async (resolve, reject) => {
        try{
            let {result} = await axios.delete(url, {
                data: data ? data : null
            });
            resolve(result);
        }catch(ex) {
            reject(ex)
        }
    });
}

export function post(url, postData) {
    return new Promise(async (resolve, reject) => {
        try{
            let {data} = await axios.post(url, postData);
            resolve(data);
        }catch(ex) {
            reject(ex)
        }
    });
}

export function put(url, postData) {
    return new Promise(async (resolve, reject) => {
        try{
            let {data} = await axios.put(url, postData);
            resolve(data);
        }catch(ex) {
            reject(ex)
        }
    });
}
