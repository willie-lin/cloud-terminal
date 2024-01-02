// LoginForm.js
import React, {useEffect, useState} from 'react';
import {check2FA, checkEmail, login} from "../../../api/api";

function LoginForm({ onLogin }) {
    const [email, setEmail] = React.useState('');
    const [password, setPassword] = React.useState('');
    const [emailError, setEmailError] = useState('');// 添加一个新的状态来保存邮箱错误信息
    const [loginError, setLoginError] = useState(''); // 保存登录错误消息
    const [isConfirmed, setIsConfirmed] = useState(false); // 新增状态变量
    const [otp, setOtp] = useState(''); // OTP 输入框的状态变量

    useEffect(() => {
        async function checkUser2FA() {
            if (email) { // 只有当 email 不为空时才发送请求
                const response = await check2FA(email);
                console.log(response); // 在控制台打印响应

                // 根据响应设置 isConfirmed 的值
                if (response && response.isConfirmed !== undefined) {
                    setIsConfirmed(response.isConfirmed);
                }
            }
        }

        checkUser2FA();
    }, [email]); // 将 email 添加到依赖数组中



    const handleEmailChange = async (e) => {
        const email = e.target.value;
        setEmail(email);
        try {
            const exists = await checkEmail(email);
            // setEmailError(exists ? 'Email already registered' : '');
            setEmailError(exists ? '' : 'Email not registered');
        } catch (error) {
            console.error(error);
        }
    };
    const getInputClass = () => {
        if (email.length === 0) {
            return "border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md";
        }
        // return emailError ? "border-green-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-green-500 focus:ring-1 rounded-md" : "border-red-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-red-500 focus:ring-1 rounded-md";
        return emailError ? "border-red-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-red-500 focus:ring-1 rounded-md" : "border-green-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-green-500 focus:ring-1 rounded-md";
    };
    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            // 将OTP赋值给totp_Secret
            const data = await login(email, password, otp); // 使用 login 函数
            console.log(data);
            localStorage.setItem('token', data.token);
            localStorage.setItem('refreshToken', data.refreshToken);
            // 将 token 存储到 localStorage 中
            onLogin(email);
        } catch (error) {
            // setLoginError(error.message);
            console.error(error);
            setLoginError('用户名或密码错误');
        }
    };
    return (
        <form onSubmit={handleSubmit}>
            <div
                className="relative flex min-h-screen text-gray-800 antialiased flex-col justify-center overflow-hidden bg-gray-50 py-6 sm:py-12">
                <div className="relative py-3 sm:w-96 mx-auto text-center">
                    <span className="text-2xl font-light">Login to your account</span>
                    <div className="mt-4 bg-white shadow-md rounded-lg text-left">
                        <div
                            className="text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded"></div>
                        <div className="px-8 py-6">
                            <label className="block font-semibold mb-2">Email</label>
                            <input type="email" value={email} onChange={handleEmailChange}
                                   className={getInputClass()}/>
                            {emailError && <p className="text-red-500">{emailError}</p>}
                            <label className="block mt-3 font-semibold mb-2">Password</label>
                            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)}
                                   className="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md"/>
                            {isConfirmed && (
                                <div>
                                    <label className="block mt-3 font-semibold mb-2">OTP</label>
                                    <input type="text" value={otp} onChange={(e) => setOtp(e.target.value)}
                                           className="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md"/>
                                </div>
                            )}
                            {loginError && <div className="text-red-500">{loginError}</div>}
                            <div className="flex justify-between items-baseline">
                                <button type="submit"
                                        className="mt-4 text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded">Login
                                </button>
                                <a href="#/" className="text-sm hover:underline">Forgot password?</a>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </form>

        // <form onSubmit={handleSubmit}>
        //     <div className="relative flex min-h-screen text-gray-800 antialiased flex-col justify-center overflow-hidden bg-gray-50 py-6 sm:py-12">
        //         <div className="relative py-3 sm:w-96 mx-auto text-center">
        //             <span className="text-2xl font-light">Login to your account</span>
        //             <div className="mt-4 bg-white shadow-md rounded-lg text-left">
        //             <div className="text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded"></div>
        //                 <div className="px-8 py-6">
        //                     <label className="block font-semibold mb-2">Email</label>
        //                     <input type="email" value={email} onChange={handleEmailChange}
        //                            className={getInputClass()}/>
        //                     {emailError && <p className="text-red-500">{emailError}</p>}
        //                     {/*<input type="email" value={email} onChange={(e) => setEmail(e.target.value)}*/}
        //                     {/*       className="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md"/>*/}
        //                     <label className="block mt-3 font-semibold mb-2">Password</label>
        //                     <input type="password" value={password} onChange={(e) => setPassword(e.target.value)}
        //                            className="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md"/>
        //                     {loginError && <div className="text-red-500">{loginError}</div>}
        //                     <div className="flex justify-between items-baseline">
        //                         <button type="submit"
        //                             // className="mt-4 bg-purple-600 text-white py-2 px-6 rounded-md hover:bg-purple-700">Login
        //                                 className="mt-4 text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded">Login
        //                         </button>
        //                         <a href="#/" className="text-sm hover:underline">Forgot password?</a>
        //                     </div>
        //                 </div>
        //             </div>
        //         </div>
        //     </div>
        // </form>
    );
}

export default LoginForm;

