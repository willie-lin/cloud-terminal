import React, { useState, useEffect } from 'react';
import {getAllUsers, getUserByEmail} from "../../api/api";

function HomePage({ email }) {
    const [currentTime, setCurrentTime] = useState(new Date());
    const [users, setUsers] = useState([]);


    const [userInfo, setUserInfo] = useState(null);


    useEffect(() => {
        getAllUsers()
            .then(data => setUsers(data))
            .catch(error => console.error('Error:', error));
    }, []);



    // 更新当前时间
    useEffect(() => {
        const timer = setInterval(() => {
            setCurrentTime(new Date());
        }, 1000); // 每秒更新一次
        return () => {
            clearInterval(timer);
        };
    }, []);

    // 获取用户信息
    useEffect(() => {
        getUserByEmail(email)
            .then(data => setUserInfo(data))
            .catch(error => console.error('Error:', error));
    }, [email]);

    return (
        <div className="flex flex-col items-center justify-center h-screen bg-blue-50">
            <div className="w-full max-w-2xl p-4 bg-white rounded shadow-md">
                <h1 className="text-4xl font-bold text-blue-600 mb-4">欢迎，{userInfo?.nickname}!</h1>
                {userInfo && <p className="text-xl text-blue-500">你的用户名是 {userInfo.username}。</p>}
                {userInfo && <p className="text-xl text-blue-500">你的电子邮件是 {userInfo.email}。</p>}
                <p className="text-xl text-blue-500">当前时间是 {currentTime.toLocaleTimeString()}。</p>
            </div>

            <div className="w-full mt-4 overflow-x-auto shadow-md sm:rounded-lg">
                <table className="min-w-full text-sm text-left text-gray-500 dark:text-gray-400">
                    <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                    <tr>
                        <th scope="col" className="px-6 py-3">ID</th>
                        <th scope="col" className="px-6 py-3">创建时间</th>
                        <th scope="col" className="px-6 py-3">更新时间</th>
                        <th scope="col" className="px-6 py-3">用户名</th>
                        <th scope="col" className="px-6 py-3">邮箱</th>
                        <th scope="col" className="px-6 py-3">在线状态</th>
                        <th scope="col" className="px-6 py-3">启用类型</th>
                        <th scope="col" className="px-6 py-3">最后登录时间</th>
                    </tr>
                    </thead>
                    <tbody>
                    {users.map(user => (
                        <tr key={user.id} className="bg-white border-b dark:bg-gray-800 dark:border-gray-700">
                            <td className="px-6 py-4">{user.id}</td>
                            <td className="px-6 py-4">{user.created_at}</td>
                            <td className="px-6 py-4">{user.updated_at}</td>
                            <td className="px-6 py-4">{user.username}</td>
                            <td className="px-6 py-4">{user.email}</td>
                            <td className="px-6 py-4">{user.online ? '在线' : '离线'}</td>
                            <td className="px-6 py-4">{user.enable_type ? '启用' : '禁用'}</td>
                            <td className="px-6 py-4">{user.last_login_time}</td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>
        </div>

    );
}

export default HomePage;
