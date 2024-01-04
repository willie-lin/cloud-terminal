import axios from 'axios';

const api = axios.create({
    baseURL: 'http://0.0.0.0:2023',
});

// login
export const login = async (email, password, totp_Secret) => {
    try {
        const response = await api.post('/api/login', { email, password, totp_Secret });
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

// checkEmail
export const checkEmail = async (email) => {
    try {
        const response = await api.post('/api/check-email', { email });
        return response.data.exists;
    } catch (error) {
        console.error(error);
    }
};

// /api/users getAllUsers

export const getAllUsers = async () => {
    try {
        const response = await api.get('/api/users', {});
        return response.data;
    } catch (error) {
        console.error(error);
    }
};



//getUserByEmail
export const getUserByEmail = async (email) => {
    try {
        const response = await api.post('/api/user/email', { email });
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

// Enable2FA
export const enable2FA = async (email) => {
    try {
        const response = await api.post('/api/enable-2fa', { email });
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

// Confirm2FA
export const confirm2FA = async (email, totp_Secret) => {
    try {
        const response = await api.post('/api/confirm-2FA', { email, totp_Secret });
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


// Validate2FA
export const validate2FA = async (email, token) => {
    try {
        const response = await api.post(`/api/validate-2fa`, { passcode: token });
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

        const response = await api.post(`/api/uploads`, formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            }
        );
        return response.data;
    } catch (error) {
        console.error(error);
    }
}
