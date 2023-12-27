// App.js
import React, {useState} from 'react';
import {BrowserRouter as Router, Navigate, Route, Routes} from 'react-router-dom';
import Navigation from "./layout/Navigation";
import Login from "./dashboard/pages/Login";
import Register from "./dashboard/pages/Register";
import HomePage from "./dashboard/pages/HomePage";




const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    const [email, setEmail] = useState(''); // 保存用户的电子邮件

    const onLogout = () => {
        setIsLoggedIn(false);
        setEmail('');
    };

    return (
        <Router>
            {/*<Dashboard />*/}
            <Navigation isLoggedIn={isLoggedIn} onLogout={onLogout} />
            <Routes>
                {/* 保存用户的电子邮件 */}
                <Route path="/login" element={isLoggedIn ? <Navigate to="/dashboard" /> : <Login onLogin={(email) => {setIsLoggedIn(true); setEmail(email);}} />} />
                {/* 保存用户的电子邮件 */}
                <Route path="/register" element={isLoggedIn ? <Navigate to="/dashboard" /> : <Register onRegister={(email) => {setIsLoggedIn(true); setEmail(email);}} />} />

                <Route path="/" element={isLoggedIn ? <HomePage /> : <Navigate to="/login" />} />
                {/* 传递电子邮件作为一个属性 */}
                <Route path="/dashboard" element={isLoggedIn ? <HomePage email={email} /> : <Navigate to="/login" />} />

                <Route path="/" element={isLoggedIn ? <HomePage /> : <Navigate to="/login" />} />
            </Routes>
        </Router>
    );
};

export default App;
