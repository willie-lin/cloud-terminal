import axios from 'axios';

const api = axios.create({
    baseURL: 'http://0.0.0.0:2023',
});

// login
export const login = async (email, password) => {
    try {
        const response = await api.post('/api/login', { email, password });
        return response.data;
    } catch (error) {
        console.error(error);
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

// checkEmail
export const checkEmail = async (email) => {
    try {
        const response = await api.post('/api/check-email', { email });
        return response.data === 'Email already registered';
    } catch (error) {
        console.error(error);
    }
};
