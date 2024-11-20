import axios from 'axios';


const api = axios.create({
    // baseURL: 'http://0.0.0.0:2023',
    baseURL: 'https://127.0.0.1:443',
    withCredentials: true,  // 添加这一行
    timeout: 1000,
    // headers: {'Authorization': `Bearer ${document.cookie.replace(/(?:^|.*;\s*)AccessToken\s*=\s*([^;]*).*$|^.*$/, "$1")}`}
});


// login
export const login = async (data) => {
    try {
        const response = await api.post('/api/login', data);
        if (response.status === 403 || response.data === 'Invalid password' || response.data === 'Invalid-OTP') {
            throw new Error('用户名或密码错误');
        }
        return response.data;
    } catch (error) {
        console.error(error);
        throw error;
    }
};

// register
export const register = async (data) => {
    try {
        const response = await api.post('/api/register', data);
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

// forgot-password
export const resetPassword = async (data) => {
    try {
        const response = await api.post('/api/reset-password', data);
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
        const response = await api.get('/admin/users',
        );
        return response.data;
    } catch (error) {
        console.error(error);
    }
};


//getUserByEmail
export const getUserByEmail = async (email) => {
    try {
        const response = await api.post('/admin/user/email', { email },
        );
        return response.data;
        return response.data.email
    } catch (error) {
        console.error(error);
    }
};

// Enable2FA
export const enable2FA = async (email) => {
    try {
        const response = await api.post('/admin/enable-2fa', { email },

        );
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

// Confirm2FA
export const confirm2FA = async (data) => {
    try {
        const response = await api.post('/admin/confirm-2FA', data,
            );
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

        const response = await api.post(`/admin/uploads`, formData,
            {
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

// edit-user-info
export const editUserInfo = async (data) => {
    try {
        const response = await api.post(`/admin/edit-userinfo`, data,
        );
        return response.data;
    } catch (error) {
        console.error(error);
    }
}

export const addUser = async (data) => {
    try {
        const response = await api.post('/admin/add-user', data);
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

export const updateUser = async (data) => {
    try {
        const response = await api.post('/admin/update-user', data);
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

export const deleteUser = async (data) => {
    try {
        const response = await api.post('/admin/delete-user', data);
        return response.data;
    } catch (error) {
        console.error(error);
    }
};
// role

// addRole
export const addRole = async (data) => {
    try {
        const response = await api.post('/admin/add-role', data);
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

export const getAllRoles = async () => {
    try {
        const response = await api.get('/admin/roles',);
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

export const deleteRole = async (data) => {
    try {
        const response = await api.post('/admin/delete-role', data);
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

export const checkRoleName = async (data) => {
    try {
        const response = await api.post('/admin/check-role-name', data);
        return response.data.exists;
    } catch (error) {
        console.error(error);
    }
};

// permission

// addPermission
export const addPermission = async (data) => {
    try {
        const response = await api.post('/admin/add-permission', data);
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

export const getAllPermissions = async () => {
    try {
        const response = await api.get('/admin/permissions',);
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

export const deletePermission = async (data) => {
    try {
        const response = await api.post('/admin/delete-permission', data);
        return response.data;
    } catch (error) {
        console.error(error);
    }
};

export const checkPermissionName = async (data) => {
    try {
        const response = await api.post('/admin/check-permission-name', data);
        return response.data.exists;
    } catch (error) {
        console.error(error);
    }
};

