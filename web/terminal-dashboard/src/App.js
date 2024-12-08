// App.js
import React, {useEffect, useState} from 'react';
import {BrowserRouter as Router, Navigate, Route, Routes} from 'react-router-dom';
import Navigation from "./layout/Navigation";
import Login from "./dashboard/pages/Login";
import Register from "./dashboard/pages/Register";
import Dashboard from "./layout/Dashboard";
import UserInfo from "./dashboard/components/user/UserInfo";
import HomePage from "./layout/HomePage";
import TwoFactorAuthPage from "./dashboard/components/2FA/TwoFactorAuthPage";
import ResetPassword from "./dashboard/pages/ResetPassword";
import NotFoundPage from "./dashboard/pages/404";
import RoleList from "./dashboard/components/role/RoleList";
import PermissionList from "./dashboard/components/permission/PermissionList";
import UserList from "./dashboard/components/user/UserList";
import { ThemeProvider } from './layout/ThemeContext';
import {logout} from "./api/api";


export const AuthContext = React.createContext(null)

const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [email, setEmail] = useState('');


    useEffect(() => {
        const storedEmail = localStorage.getItem('email');
        const refreshToken = localStorage.getItem('refreshToken');

        if (refreshToken && storedEmail) {
            setIsLoggedIn(true);
            setEmail(storedEmail);
        }
    }, []);


    const onLogin = (email, refreshToken) => {
        setEmail(email);
        localStorage.setItem('email', email); // 添加这一行来将email保存到localStorage中
        localStorage.setItem('refreshToken', refreshToken);
        setIsLoggedIn(true);
    };
    const onLogout = async () => {
        try {
        await logout(); // 调用API.js中的退出登录函数
        setIsLoggedIn(false);
        setEmail('');
        localStorage.removeItem('email'); // 添加这一行来从localStorage中删除email
        localStorage.removeItem('refreshToken');
        window.location.href = '/login'; // 重定向到登录页面
            } catch (error) {
            console.error('Logout error:', error);
        }
    };


// 在组件中使用 onLogout 方法
//     <button onClick={onLogout}>Logout</button>


    return (
        // 路由组件
        <ThemeProvider>
            <AuthContext.Provider value={{ isLoggedIn, email, onLogin, onLogout }}>
                <Router>
                    {!isLoggedIn && <Navigation />}
                    <Routes>
                        <Route path="/login" element={!isLoggedIn ? <Login onLogin={onLogin} /> : <Navigate to="/dashboard" />} />
                        <Route path="/register" element={!isLoggedIn ? <Register onRegister={onLogin} /> : <Navigate to="/dashboard" />} />
                        <Route path="/reset-password" element={!isLoggedIn ? <ResetPassword onResetPassword={onLogin} /> : <Navigate to="/dashboard" />} />

                        {/*<Route path="/" element={isLoggedIn ? <Dashboard onLogout={onLogout} email={email} /> : <Navigate to="/login" />}>*/}

                        <Route
                            path="/"
                            element={isLoggedIn ? <Dashboard /> : <Navigate to="/login" />}
                        >
                            <Route path="dashboard" element={<HomePage email={email}/>}/>
                            <Route path="users" element={<UserList email={email}/>}/>
                            <Route path="roles" element={<RoleList email={email}/>}/>
                            <Route path="permissions" element={<PermissionList email={email}/>}/>
                            <Route path="userinfo" element={<UserInfo email={email}/>}/>
                            <Route path="open-user-2fa" element={<TwoFactorAuthPage email={email} />}/>
                            <Route path="/" element={<Navigate to="dashboard" />} />
                            {/*<Route path="*" element={<NotFoundPage />} />*/}
                        </Route>
                        <Route path="*" element={<NotFoundPage />} />
                    </Routes>
                </Router>
            </AuthContext.Provider>
        </ThemeProvider>
    );
}
export default App;
