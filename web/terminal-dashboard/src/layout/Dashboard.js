import React from 'react';
import Sidebar from "./Sidebar";
import {Outlet} from "react-router-dom";
import TopNavbar from "./TopNavbar";

// Dashboard组件
function Dashboard({ isLoggedIn, onLogout, email }) {

    return (
        <div className="min-h-screen">
            <div className="w-64 bg-white shadow-lg fixed h-full z-10">
                <Sidebar onLogout={onLogout} email={email}/>
            </div>
            <div className="flex flex-col flex-grow ml-64" style={{maxWidth: 'calc(100% - 64px)'}}>
                <div className="fixed top-0 right-0 left-64 z-10">
                    <TopNavbar/>
                </div>
                <div className="flex-grow p-8 mt-8">
                    <Outlet/>
                </div>
            </div>
        </div>
    );
}

export default Dashboard;
