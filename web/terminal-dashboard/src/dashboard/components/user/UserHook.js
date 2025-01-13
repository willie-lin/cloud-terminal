// 自定义钩子函数，用于获取用户信息
import {useEffect, useState} from "react";
import {getAllUsers, getUserByEmail} from "../../../api/api";

// 自定义钩子函数，用于获取所有用户
export const useFetchUsers = () => {
    const [users, setUsers] = useState([]);

    useEffect(() => {
        getAllUsers()
            .then(data => setUsers(data))
            .catch(error => console.error('Error:', error));
    }, []);
    return users;
};
// 自定义钩子函数，用于获取用户信息
export const useFetchUserInfo = (email) => {
    const [userInfo, setUserInfo] = useState(null);

    useEffect(() => {
        if (email) {
            getUserByEmail(email)
                .then(data => setUserInfo(data))
                .catch(error => console.error('Error:', error));
        }
    }, [email]);

    return userInfo;
};



// 自定义钩子函数，用于更新当前时间
export const useCurrentTime = () => {
    const [currentTime, setCurrentTime] = useState(new Date());

    useEffect(() => {
        const timer = setInterval(() => {
            setCurrentTime(new Date());
        }, 1000); // 每秒更新一次
        return () => {
            clearInterval(timer);
        };
    }, []);

    return currentTime;
}


// 自定义钩子函数，用于获取当前地区的格式化年月日和时间
export const useCurrentDateTime = () => {
    const [currentDateTime, setCurrentDateTime] = useState('');

    useEffect(() => {
        const updateDateTime = () => {
            const now = new Date();
            const formattedDateTime = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')} ${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}:${String(now.getSeconds()).padStart(2, '0')}`;
            setCurrentDateTime(formattedDateTime);
        };

        updateDateTime(); // 初始化调用一次

        const timer = setInterval(updateDateTime, 1000); // 每秒更新一次

        return () => {
            clearInterval(timer);
        };
    }, []);

    return currentDateTime;
};



