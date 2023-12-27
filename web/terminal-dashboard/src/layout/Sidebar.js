import React from 'react';
import { Link } from 'react-router-dom';

const Sidebar = ({ isLoggedIn, onLogout }) => {
    return (
        <aside className="w-64 bg-white shadow-md">
            <div className="p-6">
                <Link to="/" className="-m-1.5 p-1.5">
                    <span className="sr-only">Your Company</span>
                    <img className="h-8 w-auto"
                         src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=600" alt=""/>
                </Link>
                <div className="mt-8">
                    <Link to="/product" className="text-sm font-semibold leading-6 text-gray-900 block py-1">Product</Link>
                    <Link to="/features" className="text-sm font-semibold leading-6 text-gray-900 block py-1">Features</Link>
                    <Link to="/marketplace" className="text-sm font-semibold leading-6 text-gray-900 block py-1">Marketplace</Link>
                    <Link to="/company" className="text-sm font-semibold leading-6 text-gray-900 block py-1">Company</Link>
                </div>
                <div className="mt-8">
                    {isLoggedIn ? (
                        <Link to="/login" onClick={onLogout}
                              className="text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded block">Logout</Link>
                    ) : (
                        <>
                            <Link to="/login"
                                  className="text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded block">Log in</Link>
                            <Link to="/register"
                                  className="text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 ml-4 px-4 py-2 rounded block mt-2">Register</Link>
                        </>
                    )}
                </div>
            </div>
        </aside>
    );
};

export default Sidebar;
