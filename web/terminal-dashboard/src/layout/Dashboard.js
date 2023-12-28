import React from 'react';
import Sidebar from "./Sidebar";
import HomePage from "./HomePage";

// Dashboard组件
const Dashboard = ({ isLoggedIn, onLogout }) => {
        return (
            <div className="Dashboard flex">
                <div className="min-h-screen bg-blue-gray-50/50 w-64">
                    <Sidebar isLoggedIn={isLoggedIn} onLogout={onLogout}/>
                </div>
                <div className="p-4 flex-grow">
                    <HomePage/>
                </div>
            </div>
        );
};


export default Dashboard;
