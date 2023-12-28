// App.js
import React, {useState} from 'react';
import {BrowserRouter as Router, Navigate, Route, Routes} from 'react-router-dom';
import Navigation from "./layout/Navigation";
import Login from "./dashboard/pages/Login";
import Register from "./dashboard/pages/Register";
import Dashboard from "./layout/Dashboard";
const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(!!localStorage.getItem('token'));
    const [email, setEmail] = useState(localStorage.getItem('email') || ''); // 修改这一行来从localStorage中恢复email
    const onLogin = (email, token) => {
        setEmail(email);
        localStorage.setItem('email', email); // 添加这一行来将email保存到localStorage中
        localStorage.setItem('token', token);
        setIsLoggedIn(true);
    };
    const onLogout = () => {
        setIsLoggedIn(false);
        setEmail('');
        localStorage.removeItem('email'); // 添加这一行来从localStorage中删除email
        localStorage.removeItem('token');
    };
    return (
        <Router>
            {/*<Dashboard />*/}
            <Navigation isLoggedIn={isLoggedIn} onLogout={onLogout} />
            <Routes>
                {/* 保存用户的电子邮件 */}
                <Route path="/login" element={isLoggedIn ? <Navigate to="/dashboard" /> : <Login onLogin={onLogin} />} />
                {/* 保存用户的电子邮件 */}
                <Route path="/register" element={isLoggedIn ? <Navigate to="/dashboard" /> : <Register onRegister={onLogin} />} />

                <Route path="/" element={isLoggedIn ? <Dashboard /> : <Navigate to="/login" />} />
                {/* 传递电子邮件作为一个属性 */}
                <Route path="/dashboard" element={isLoggedIn ? <Dashboard email={email} /> : <Navigate to="/login" />} />

                <Route path="/" element={isLoggedIn ? <Dashboard /> : <Navigate to="/login" />} />
            </Routes>
        </Router>
    );
};

export default App;
