import React from 'react';
import Sidebar from "./Sidebar";
import {Outlet} from "react-router-dom";
import TopNavbar from "./TopNavbar";

// Dashboard组件
function Dashboard({ isLoggedIn, onLogout, email }) {

    return (
        <div className="flex min-h-screen bg-blue-gray-50">
            <div className="w-64 bg-white shadow-lg fixed h-full">
                <Sidebar onLogout={onLogout}/>
            </div>
            <div className="flex flex-col flex-grow ml-64">
                <TopNavbar/>
                <div className="flex-grow p-8">
                    <Outlet/>
                </div>
            </div>
        </div>
    );
}


export default Dashboard;
