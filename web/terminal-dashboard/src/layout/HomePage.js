import React, { useState, useEffect } from 'react';
import {getAllUsers, getUserByEmail} from "../api/api";
import UserInfo from "../dashboard/components/UserInfo";
import UserList from "../dashboard/components/UserList";

function HomePage({ email }) {
    const [currentTime, setCurrentTime] = useState(new Date());
    const [users, setUsers] = useState([]);
    const [userInfo, setUserInfo] = useState(null);

    // 获取用户信息
    useEffect(() => {
        if (email) {  // 添加这一行来检查email是否存在
            getUserByEmail(email)
                .then(data => setUserInfo(data))
                .catch(error => console.error('Error:', error));
        }
    }, [email]);

    // 更新当前时间
    useEffect(() => {
        const timer = setInterval(() => {
            setCurrentTime(new Date());
        }, 1000); // 每秒更新一次
        return () => {
            clearInterval(timer);
        };
    }, []);

    // 获取所有用户
    useEffect(() => {
        getAllUsers()
            .then(data => setUsers(data))
            .catch(error => console.error('Error:', error));
    }, []);



    return (
        <div className="flex flex-col items-center justify-center h-screen bg-blue-50">
            <UserInfo userInfo={userInfo} currentTime={currentTime} />
            <UserList users={users} />
        </div>

    );
}

export default HomePage;
