import React from 'react';
import Sidebar from "./Sidebar";
import HomePage from "./HomePage";

// Dashboard组件
function Dashboard({ isLoggedIn, onLogout, email }) {

    return (
        <div className="flex min-h-screen bg-blue-gray-50">
            <div className="w-64 bg-white shadow-lg fixed h-full">
                <Sidebar onLogout={onLogout}/>
            </div>
            <div className="flex-grow p-8 ml-64">
                <HomePage email={email}/>
            </div>
        </div>
    );
};


export default Dashboard;
