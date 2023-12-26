// App.js
import React, {useState} from 'react';
import {BrowserRouter as Router, Navigate, Route, Routes} from 'react-router-dom';
import Navigation from "./dashboard/Navigation";
import Login from "./dashboard/pages/Login";
import Register from "./dashboard/pages/Register";
import HomePage from "./dashboard/pages/homepage";




const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    return (
        <Router>
            <Navigation />
            <Routes>
                {/*<Route path="/login" element={<Login />} />*/}
                {/*<Route path="/register" element={<Register />} />*/}

                <Route path="/login" element={isLoggedIn ? <Navigate to="/" /> : <Login onLogin={() => setIsLoggedIn(true)} />} />
                <Route path="/register" element={isLoggedIn ? <Navigate to="/" /> : <Register onRegister={() => setIsLoggedIn(true)} />} />
                {/*<Route path="/register" element={<Register />} />*/}
                <Route path="/" element={isLoggedIn ? <HomePage /> : <Navigate to="/login" />} />

                {/* 其他路由... */}
            </Routes>
        </Router>
    );
};

export default App;
