import React from 'react';
import Navigation from "./Navigation";
import Sidebar from "./Sidebar";

// Dashboard组件
const Dashboard = ({ children }) => {
    return (
        <div className="flex h-screen bg-gray-200">
            {/*<Sidebar />*/}
            <div className="flex flex-col w-full">
                <Navigation />
                <main className="h-full overflow-y-auto">
                    <div className="container mx-auto px-6 py-8">
                        {children}
                    </div>
                </main>
            </div>
        </div>
    );
};
// Sidebar组件


export default Dashboard;
