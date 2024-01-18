import React, { useState } from 'react';

function AddUserForm({ onAddUser }) {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const handleSubmit = (event) => {
        event.preventDefault();
        onAddUser({ email, password });
        setEmail('');
        setPassword('');
    };

    return (
        <form onSubmit={handleSubmit}>
            <label>
                Email:
                <input type="email" value={email} onChange={e => setEmail(e.target.value)} required />
            </label>
            <label>
                Password:
                <input type="password" value={password} onChange={e => setPassword(e.target.value)} required />
            </label>
            <button type="submit">Add User</button>
        </form>
    );
}

// 在你的 UserList 组件中使用 AddUserForm 组件：
function UserList() {
    // ...你的其他代码...

    const handleAddUser = (newUser) => {
        // 这里是处理新用户的代码，例如发送一个请求到后端服务
    };

    return (
        <div>
            {/* ...你的其他代码... */}
            <AddUserForm onAddUser={handleAddUser} />
        </div>
    );
}
