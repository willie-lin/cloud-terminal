// 用户信息组件
import {useEffect, useState} from "react";
import {getUserByEmail} from "../../../api/api";

function UserInfo({ email }) {
    const [currentTime, setCurrentTime] = useState(new Date());
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


    return (
        <div className="w-full max-w-2xl p-4 bg-white rounded shadow-md">
            {/*<div className="flex-grow w-full max-w-2xl p-4 bg-white rounded shadow-md">*/}
            <h1 className="text-4xl font-bold text-blue-600 mb-4">欢迎，{userInfo?.nickname}!</h1>
            {userInfo && <p className="text-xl text-blue-500">你的用户名是 {userInfo.username}。</p>}
            {userInfo && <p className="text-xl text-blue-500">你的电子邮件是 {userInfo.email}。</p>}
            <p className="text-xl text-blue-500">当前时间是 {currentTime.toLocaleTimeString()}。</p>
        </div>
    );
}

export default UserInfo;