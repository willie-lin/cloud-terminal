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


// export const register = (email, password) => api.post('/register', { email, password });
// export const getUser = (id) => api.get(`/users/${id}`);
// export const updateUser = (id, data) => api.put(`/users/${id}`, data);