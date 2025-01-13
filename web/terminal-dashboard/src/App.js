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
import AccessPolicyList from "./dashboard/components/accesspolicy/AccessPolicyList";


export const AuthContext = React.createContext(null)

const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [user, setUser] = useState(null);

    useEffect(() => {
        const storedUser = JSON.parse(localStorage.getItem('user'));
        const accessToken = localStorage.getItem('accessToken');
        const refreshToken = localStorage.getItem('refreshToken');

        if (refreshToken && storedUser && accessToken) {
            setIsLoggedIn(true);
            setUser(storedUser);
        }
    }, []);


    const onLogin = (user, accessToken, refreshToken) => {
        setUser(user);
        localStorage.setItem('user', JSON.stringify(user));
        localStorage.setItem('accessToken', accessToken);
        localStorage.setItem('refreshToken', refreshToken);
        setIsLoggedIn(true);
    };

    const onLogout = async () => {
        try {
            await logout(); // 调用API.js中的退出登录函数
            setIsLoggedIn(false);
            setUser(null);
            localStorage.removeItem('user');
            localStorage.removeItem('accessToken');
            localStorage.removeItem('refreshToken');
            window.location.href = '/login'; // 重定向到登录页面
        } catch (error) {
            console.error('Logout error:', error);
        }
    };

    return (
        <ThemeProvider>
            <AuthContext.Provider value={{isLoggedIn, currentUser: user, onLogin, onLogout}}>
                <Router>
                    {!isLoggedIn && <Navigation/>}
                    <Routes>
                        <Route path="/login"
                               element={!isLoggedIn ? <Login onLogin={onLogin}/> : <Navigate to="/dashboard"/>}/>
                        <Route path="/register"
                               element={!isLoggedIn ? <Register onRegister={onLogin}/> : <Navigate to="/dashboard"/>}/>
                        <Route path="/reset-password"
                               element={!isLoggedIn ? <ResetPassword onResetPassword={onLogin}/> :
                                   <Navigate to="/dashboard"/>}/>

                        <Route path="/" element={isLoggedIn ? <Dashboard/> : <Navigate to="/login"/>}>
                            <Route path="dashboard" element={<HomePage user={user}/>}/>
                            <Route path="users" element={<UserList user={user}/>}/>
                            <Route path="roles" element={<RoleList user={user}/>}/>
                            <Route path="policies" element={<AccessPolicyList user={user}/>}/>
                            <Route path="permissions" element={<PermissionList user={user}/>}/>
                            <Route path="userinfo" element={<UserInfo user={user}/>}/>
                            <Route path="open-user-2fa" element={<TwoFactorAuthPage user={user}/>}/>
                            <Route path="/" element={<Navigate to="dashboard"/>}/>
                        </Route>
                        <Route path="*" element={<NotFoundPage/>}/>
                    </Routes>
                </Router>
            </AuthContext.Provider>
        </ThemeProvider>
    );
}
export default App;
