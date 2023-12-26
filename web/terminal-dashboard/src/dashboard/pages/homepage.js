import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom'; // 首先，你需要从 'react-router-dom' 导入 useNavigate



function HomePage({ email }) {
    const [currentTime, setCurrentTime] = useState(new Date());
    const navigate = useNavigate(); // 然后，在你的组件中使用 useNavigate

    // 更新当前时间
    useEffect(() => {
        const timer = setInterval(() => {
            setCurrentTime(new Date());
        }, 1000); // 每秒更新一次

        // 清理函数
        return () => {
            clearInterval(timer);
        };
    }, []);

    // 退出登录的函数
    const logout = () => {
        function clearUserInfo() {
            // 清除localStorage中的用户信息
            localStorage.removeItem('username');
            localStorage.removeItem('token');
            // 你可以添加更多的代码来清除其他的用户信息
        }

        // 在这里清除用户信息
        clearUserInfo(); // 清除用户信息
        // 然后跳转到登录页面
        navigate('/login');
    };

    return (
        <div>
            <h1>Welcome, { email }!</h1>
            <p>The current time is {currentTime.toLocaleTimeString()}.</p>
            <button onClick={logout}>退出登录</button> {/* 添加一个退出登录的按钮 */}
        </div>
    );
}

export default HomePage;
