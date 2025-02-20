import React, {useContext} from 'react';
import Sidebar from "./Sidebar";
import {Outlet} from "react-router-dom";
import TopNavbar from "./TopNavbar";
import Footer from "./Footer";
import {AuthContext} from "../App";
import { useTheme } from './ThemeContext';

// Dashboard组件
function Dashboard() {
    const { email, onLogout } = useContext(AuthContext);
    const { isDarkMode } = useTheme()


    return (
        <div className={`flex h-screen overflow-hidden ${
                isDarkMode
                    ? 'dark bg-gray-900 text-gray-100'
                    : 'light bg-gray-100 text-gray-900'
            }`}>
                <div className={`w-64 flex-shrink-0 ${
                    isDarkMode ? 'bg-gray-800' : 'bg-white'
                } overflow-y-auto`}>
                    <Sidebar email={email} onLogout={onLogout} />
                </div>
                <div className="flex flex-col flex-grow overflow-hidden">
                    <div className={`${
                        isDarkMode ? 'bg-gray-800' : 'bg-white'
                    } shadow-md z-10`}>
                        <TopNavbar />
                    </div>

                    <main className="flex-grow overflow-x-hidden overflow-y-auto">
                        {/*<div className="container mx-auto px-4 py-8 max-w-full">*/}
                            <Outlet />
                        {/*</div>*/}
                    </main>

                    <div className={`${
                        isDarkMode ? 'bg-gray-800' : 'bg-white'
                    } w-full`}>
                        <Footer />
                    </div>
                </div>
            </div>
    );
}

export default Dashboard;
