import React from 'react';
import Sidebar from "./Sidebar";
import {Outlet} from "react-router-dom";
import TopNavbar from "./TopNavbar";
import Footer from "./Footer";

// Dashboard组件
function Dashboard({ isLoggedIn, onLogout, email }) {

    return (
        <div className="min-h-screen">
            <div className="w-64 bg-white shadow-lg fixed h-full z-10">
            <Sidebar onLogout={onLogout} email={email}/>
            </div>
            <div className="flex flex-col flex-grow ml-64"  style={{maxWidth: 'calc(100% - 64px)', height: '100vh'}}>
                <div className="sticky top-0 right-0 left-64 z-10">
                    <TopNavbar/>
                </div>
                <div className="sticky p-8 mt-0 flex-grow overflow-auto">
                    <Outlet/>
                </div>
                <div className="sticky bottom-0">
                    <Footer />
                </div>
            </div>
        </div>
    );
}

export default Dashboard;
