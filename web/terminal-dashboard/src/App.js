// App.js
import React, {useState} from 'react';
import {BrowserRouter as Router, Navigate, Route, Routes} from 'react-router-dom';
import Navigation from "./dashboard/Navigation";
import Login from "./dashboard/pages/Login";
import Register from "./dashboard/pages/Register";
import Dashboard from "./dashboard/pages/Dashboard";




const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    const [email, setEmail] = useState(''); // 保存用户的电子邮件

    return (
        <Router>
            <Navigation isLoggedIn={isLoggedIn} />
            <Routes>
                {/*<Route path="/login" element={isLoggedIn ? <Navigate to="/" /> : <Login onLogin={() => setIsLoggedIn(true)} />} />*/}
                {/*<Route path="/" element={isLoggedIn ? <Navigate to="/" /> : <Login onLogin={() => setIsLoggedIn(true)} />} />*/}
                {/*<Route path="/register" element={isLoggedIn ? <Navigate to="/" /> : <Register onRegister={() => setIsLoggedIn(true)} />} />*/}
                <Route path="/login" element={isLoggedIn ? <Navigate to="/dashboard" /> : <Login onLogin={(email) => {setIsLoggedIn(true); setEmail(email);}} />} /> {/* 保存用户的电子邮件 */}
                <Route path="/register" element={isLoggedIn ? <Navigate to="/dashboard" /> : <Register onRegister={(email) => {setIsLoggedIn(true); setEmail(email);}} />} /> {/* 保存用户的电子邮件 */}
                <Route path="/" element={isLoggedIn ? <Dashboard /> : <Navigate to="/login" />} />
                <Route path="/dashboard" element={isLoggedIn ? <Dashboard email={email} /> : <Navigate to="/login" />} /> {/* 传递电子邮件作为一个属性 */}

                {/*<Route path="/login" element={isLoggedIn ? <Navigate to="/dashboard" /> : <Login onLogin={() => setIsLoggedIn(true)} />} />*/}
                {/*<Route path="/register" element={isLoggedIn ? <Navigate to="/dashboard" /> : <Register onRegister={() => setIsLoggedIn(true)} />} />*/}
                {/*<Route path="/" element={isLoggedIn ? <Dashboard /> : <Navigate to="/login" />} />*/}
                {/*<Route path="/dashboard" element={isLoggedIn ? <Dashboard /> : <Navigate to="/login" />} /> /!* 添加新的路由 *!/*/}


                {/*<Route path="/register" element={<Register />} />*/}
                <Route path="/" element={isLoggedIn ? <Dashboard /> : <Navigate to="/login" />} />

                {/* 其他路由... */}
            </Routes>
        </Router>
    );
};

export default App;
