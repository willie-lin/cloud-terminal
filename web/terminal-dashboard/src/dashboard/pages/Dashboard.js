import React, { useState, useEffect } from 'react';

function Dashboard({ email }) {
    const [currentTime, setCurrentTime] = useState(new Date());

    // 更新当前时间
    useEffect(() => {
        const timer = setInterval(() => {
            setCurrentTime(new Date());
        }, 1000); // 每秒更新一次
        return () => {
            clearInterval(timer);
        };
    }, []);

    return (
        <div className="flex flex-col items-center justify-center h-screen bg-blue-50">
            <h1 className="text-4xl font-bold text-blue-600 mb-4">欢迎，{email}!</h1>
            <p className="text-xl text-blue-500">当前时间是 {currentTime.toLocaleTimeString()}。</p>
        </div>
    );
}

export default Dashboard;
