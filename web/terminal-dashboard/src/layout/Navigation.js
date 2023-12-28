import React from 'react';
import { Link } from 'react-router-dom';

const Navigation = ({ isLoggedIn, onLogout }) => {
    return (
        <div className="bg-white">
            <header className="absolute inset-x-0 top-0 z-50">
                <nav className="flex items-center justify-between p-6 lg:px-8" aria-label="Global">
                    <div className="flex lg:flex-1">
                        <Link to="/" className="-m-1.5 p-1.5">
                            <span className="sr-only">Your Company</span>
                            <img className="h-8 w-auto"
                                 src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=600" alt=""/>
                        </Link>
                    </div>
                    {/*<div className="hidden lg:flex lg:gap-x-12">*/}
                    {/*    <Link to="/product" className="text-sm font-semibold leading-6 text-gray-900">Product</Link>*/}
                    {/*    <Link to="/features" className="text-sm font-semibold leading-6 text-gray-900">Features</Link>*/}
                    {/*    <Link to="/marketplace"*/}
                    {/*          className="text-sm font-semibold leading-6 text-gray-900">Marketplace</Link>*/}
                    {/*    <Link to="/company" className="text-sm font-semibold leading-6 text-gray-900">Company</Link>*/}
                    {/*</div>*/}
                    <div className="hidden lg:flex lg:flex-1 lg:justify-end">
                        {isLoggedIn ? (
                            <>
                                {/*<div className="relative">*/}
                                {/*    <img className="w-10 h-10 rounded"*/}
                                {/*         src="../../../assets/soft-ui-flowbite/images/people/profile-picture-5.jpg" alt=""/>*/}
                                {/*    <span*/}
                                {/*        className="absolute bottom-0 left-8 transform translate-y-1/4 w-3.5 h-3.5 bg-green-400 border-2 border-white dark:border-gray-800 rounded-full"></span>*/}
                                {/*</div>*/}
                                {/*<Link to="/login" onClick={onLogout}*/}
                                {/*      className="text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded">Logout</Link>*/}
                            </>
                        ) : (
                            <>
                                <Link to="/login"
                                      className="text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded">Log
                                    in</Link>
                                <Link to="/register"
                              className="text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 ml-4 px-4 py-2 rounded">Register</Link>
                            </>
                            )}
                    </div>
                </nav>
            </header>
        </div>
    );
};

export default Navigation;
