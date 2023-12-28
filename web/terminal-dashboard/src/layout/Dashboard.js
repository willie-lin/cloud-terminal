import React from 'react';
import Sidebar from "./Sidebar";
import HomePage from "./HomePage";

// Dashboard组件
const Dashboard = ({ isLoggedIn, onLogout, email }) => {
        return (
            <div className="flex min-h-screen bg-blue-gray-50">
                <div className="w-64 bg-white shadow-lg">
                    <Sidebar onLogout={onLogout}/>
                </div>
                <div className="flex-grow p-8">
                    <HomePage email={ email }/>
                </div>
            </div>
        );
};


export default Dashboard;
