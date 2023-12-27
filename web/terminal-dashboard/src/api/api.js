import axios from 'axios';

const api = axios.create({
    baseURL: 'http://0.0.0.0:2023',
});

// login
export const login = async (email, password) => {
    try {
        const response = await api.post('/api/login', { email, password });
        if (response.status === 403 || response.data === 'Invalid password') {
        // if (response.data === 'Invalid-password') {
            throw new Error('用户名或密码错误');
        }
        return response.data;
    } catch (error) {
        console.error(error);
        throw error;
    }
};

// export const login = async (email, password) => {
//     try {
//         const response = await api.post('/api/login', { email, password });
//         if (response.status === 401) {
//             throw new Error('用户名或密码错误');
//         }
//         return response.data;
//     } catch (error) {
//         console.error(error);
//         throw error;
//     }
// };

// register
export const register = async (email, password) => {
    try {
        const response = await api.post('/api/register', { email, password });
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
