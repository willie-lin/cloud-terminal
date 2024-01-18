import axios from 'axios';
import Cookies from 'js-cookie';


const api = axios.create({
    // baseURL: 'http://0.0.0.0:2023',
    baseURL: 'https://127.0.0.1:443',
    withCredentials: true,  // 添加这一行
    timeout: 1000,
    // headers: {'Authorization': `Bearer ${document.cookie.replace(/(?:^|.*;\s*)AccessToken\s*=\s*([^;]*).*$|^.*$/, "$1")}`}
});


// login
export const login = async (email, password, otp) => {
    try {
        const response = await api.post('/api/login', { email, password, otp });
        if (response.status === 403 || response.data === 'Invalid password' || response.data === 'Invalid-OTP') {
        // if (response.data === 'Invalid-password') {
            throw new Error('用户名或密码错误');
        }
        return response.data;
    } catch (error) {
        console.error(error);
        throw error;
    }
};

// register
export const register = async (email, password) => {
    try {
        const response = await api.post('/api/register', { email, password });
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

// forgot-password
export const resetPassword = async (email, password) => {
    try {
        const response = await api.post('/api/reset-password', { email, password });
        return response.data;
    } catch (error) {
        console.error(error);
    }
};



// checkEmail
export const checkEmail = async (email) => {
    try {
        const response = await api.post('/api/check-email', { email });
        return response.data.exists;
    } catch (error) {
        console.error(error);
    }
};

// api/users getAllUsers
export const getAllUsers = async () => {
    try {
        // const token = Cookies.get('user_token');
        // console.log(token)
        const response = await api.get('/admin/users',
            // {headers: {"Authorization": "Bearer " + localStorage.getItem("user_token")}}
            // {headers: {'Authorization': `Bearer ${token}`}}
        );
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

// export const getAllUsers = async () => {
//     try {
//         const response = await api.get('/api/users', {
//         });
//         return response.data;
//     } catch (error) {
//         console.error(error);
//     }
// };


//getUserByEmail
export const getUserByEmail = async (email) => {
    try {
        const response = await api.post('/admin/user/email', { email },{
            headers: {
                "Authorization": "Bearer " + localStorage.getItem("user_token")}
        });
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

// Enable2FA
export const enable2FA = async (email) => {
    try {
        const response = await api.post('/admin/enable-2fa', { email }, {
            headers: {
                "Authorization": "Bearer " + localStorage.getItem("user_token")}
        });
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

// Confirm2FA
export const confirm2FA = async (email, otp, secret) => {
    try {
        const response = await api.post('/admin/confirm-2FA', { email, otp, secret },
            {
                headers: {
                    "Authorization": "Bearer " + localStorage.getItem("user_token")}
            });
        if (response.status === 400) {
            throw new Error('Invalid TOTP secret');
        }
        return response.data;
    } catch (error) {
        console.error(error);
        throw error;
    }
};


//
export const check2FA = async (email) => {
    try {
        const response = await api.post('/api/check-2FA', { email } );
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

// upload file
export const uploadFile = async (file) => {
    try {
        const formData = new FormData();
        formData.append('file', file);

        const response = await api.post(`/admin/uploads`, formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                    "Authorization": "Bearer " + localStorage.getItem("user_token")
                },
            }
        );
        return response.data;
    } catch (error) {
        console.error(error);
    }
}

// edit-user-info
export const editUserInfo = async (data) => {
    try {
        const response = await api.post(`/admin/edit-userinfo`, data,
            {
                headers: {
                    "Authorization": "Bearer " + localStorage.getItem("user_token")}
            });
        return response.data;
    } catch (error) {
        console.error(error);
    }
}
