// RegisterForm.js
import React, { useState } from 'react';
import {checkEmail, register} from "../../api/api";

function RegisterForm({ onRegister }) {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [emailError, setEmailError] = useState('');// 添加一个新的状态来保存邮箱错误信息

// 当 email 状态改变时，调用 checkEmail 函数
//     const handleEmailChange = async (e) => {
//         const email = e.target.value;
//         setEmail(email);
//         try {
//             const response = await checkEmail(email);
//             if (response.status === 200) {
//                 setEmailError('');
//             } else if (response.status === 400) {
//                 setEmailError('Email already registered');
//             }
//         } catch (error) {
//             console.error(error);
//         }
//     };
    const handleEmailChange = async (e) => {
        const email = e.target.value;
        setEmail(email);
        try {
            const exists = await checkEmail(email);
            setEmailError(exists ? 'Email already registered' : '');
        } catch (error) {
            console.error(error);
        }
    };


    const handleSubmit = async (e) => {
        e.preventDefault();
        if (password !== confirmPassword) {
            alert("Passwords don't match");
            return;
        }
        if (emailError) {
            alert(emailError);
            return;
        }
        try {
            const data = await register(email, password); // 使用 register 函数
            console.log(data);
            onRegister(data.email);
        } catch (error) {
            console.error(error);
        }
    };

    const getInputClass = () => {
        if (email.length === 0) {
            return "border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md";
        }
        return emailError ? "border-red-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-red-500 focus:ring-1 rounded-md" : "border-green-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-green-500 focus:ring-1 rounded-md";
    };


    return (
        <form onSubmit={handleSubmit}>
        <div className="relative flex min-h-screen text-gray-800 antialiased flex-col justify-center overflow-hidden bg-gray-50 py-6 sm:py-12">
            <div className="relative py-3 sm:w-96 mx-auto text-center">
                <span className="text-2xl font-light">Register for an account</span>
                <div className="mt-4 bg-white shadow-md rounded-lg text-left">
                    {/*<div className="h-2 bg-purple-600 rounded-t-md"></div>*/}
                    <div className="text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded"></div>
                    <div className="px-8 py-6">
                        <label className="block font-semibold mb-2">Email</label>
                        <input type="email" value={email} onChange={handleEmailChange}
                               className={getInputClass()}/>
                        {emailError && <p className="text-red-500">{emailError}</p>}
                        {/*<input type="email" value={email} onChange={handleEmailChange}*/}
                        {/*       className="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md"/>*/}
                        {/*{emailError && <p className="text-red-500">{emailError}</p>}*/}


                        <label className="block mt-3 font-semibold mb-2">Password</label>
                        <input type="password" value={password} onChange={(e) => setPassword(e.target.value)}
                               className="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md"/>
                        <label className="block mt-3 font-semibold mb-2">Confirm Password</label>
                        <input type="password" value={confirmPassword}
                               onChange={(e) => setConfirmPassword(e.target.value)}
                               className="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md"/>
                        <div className="flex justify-between items-baseline">
                            {/*<button type="submit" className="mt-4 bg-purple-600 text-white py-2 px-6 rounded-md hover:bg-purple-700">Register</button>*/}
                            <button type="submit"
                                    className="mt-4 text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded">Register
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        </form>
    );
}

export default RegisterForm;
