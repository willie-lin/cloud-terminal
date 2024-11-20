import React, {useContext} from 'react';
import Sidebar from "./Sidebar";
import {Outlet} from "react-router-dom";
import TopNavbar from "./TopNavbar";
import Footer from "./Footer";
import {AuthContext} from "../App";

// Dashboard组件
function Dashboard() {
    const { email, onLogout } = useContext(AuthContext);


    return (
        <div className="flex h-screen overflow-hidden bg-gray-100">
            {/*<div className="min-h-screen bg-white dark:bg-gray-800">*/}
            {/*    <div className="w-64 bg-white shadow-lg fixed h-full z-10">*/}
            <div className="w-64 flex-shrink-0 bg-white overflow-y-auto">
                <Sidebar email={email} onLogout={onLogout}/>
            </div>
            {/*<div className="flex flex-col flex-grow ml-64" style={{maxWidth: 'calc(100% - 64px)', height: '100vh'}}>*/}
            {/*    <div className="sticky top-0 right-0 left-64 z-10">*/}
            {/*        <TopNavbar/>*/}
            {/*    </div>*/}
            {/*    <div className="sticky p-4 mt-0 flex-grow overflow-auto">*/}
            {/*        /!*<div className="flex-grow p-4 overflow-auto">*!/*/}
            {/*        <Outlet/>*/}
            {/*    </div>*/}
            {/*    <div className="sticky bottom-0">*/}
            {/* Main content area */}
            <div className="flex flex-col flex-grow overflow-hidden">
                <div className="bg-white shadow-md z-10">
                    <TopNavbar/>
                </div>

                {/* Main content */}
                <main className="flex-grow overflow-x-hidden overflow-y-auto">
                    <div className="container mx-auto px-4 py-8 max-w-full">
                        <Outlet/>
                    </div>
                </main>

                {/* Footer */}
                <div className="bg-white w-full">
                    <Footer/>
                </div>
            </div>
        </div>
    );
}

export default Dashboard;
