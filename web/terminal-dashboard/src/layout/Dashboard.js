import React, {createContext, useContext, useState} from 'react';
import Sidebar from "./Sidebar";
import {Outlet} from "react-router-dom";
import TopNavbar from "./TopNavbar";

// Dashboard组件
function Dashboard({ isLoggedIn, onLogout, email }) {

    // 创建一个context
    const UserContext = createContext();

// 创建一个Provider组件
    function UserProvider({ children }) {
        const [userInfo, setUserInfo] = useState(null);

        const handleUpdate = (newInfo) => {
            // 更新用户信息
            setUserInfo(newInfo);
        };

        // 将状态和更新函数作为value传递给Provider
        return (
            <UserContext.Provider value={{ userInfo, onUpdate: handleUpdate }}>
                {children}
            </UserContext.Provider>
        );
    }

// 创建一个自定义hook，用于在子组件中访问context
    function useUser() {
        return useContext(UserContext);
    }

    return (
        <div>
            <div className="w-64 bg-white shadow-lg fixed h-full z-10">
                <Sidebar onLogout={onLogout} email={email}/>
            </div>
            <div className="flex flex-col flex-grow ml-64" style={{maxWidth: 'calc(100% - 64px)'}}>
                <div className="fixed top-0 right-0 left-64 z-10">
                    <TopNavbar/>
                </div>
                <div className="flex-grow p-4 mt-12">
                    <Outlet/>
                </div>
            </div>
        </div>
        // </UserProvider>
    );
}

export default Dashboard;
