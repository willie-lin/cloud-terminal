// LoginForm.js
import React, {useEffect, useState} from 'react';
import {check2FA, checkEmail, login} from "../../../api/api";
import {Alert, Button, Card, Input} from "@material-tailwind/react";

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

    useEffect(() => {
        if (loginError) {
            const timer = setTimeout(() => {
                setLoginError('');
            }, 500); // 设置3秒后自动隐藏

            return () => clearTimeout(timer); // 清除定时器
        }
    }, [loginError]);

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
    // const getInputClass = () => {
    //     if (email.length === 0) {
    //         return "border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md";
    //     }
    //     // return emailError ? "border-green-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-green-500 focus:ring-1 rounded-md" : "border-red-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-red-500 focus:ring-1 rounded-md";
    //     return emailError ? "border-red-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-red-500 focus:ring-1 rounded-md" : "border-green-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-green-500 focus:ring-1 rounded-md";
    // };
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

        <form onSubmit={handleSubmit} className="flex flex-col items-center justify-center min-h-screen">
            <Card className="w-full max-w-md"> {/* 设置最大宽度 */}
                <h6 className="text-gray-500 text-lg text-center">Login to your account</h6>
                <hr className="mb-6 border-b-1 border-gray-300"/>
                <div className="mb-4">
                    <Input
                        type="email"
                        color="lightBlue"
                        size="regular"
                        outline={true}
                        placeholder="Email"
                        value={email}
                        onChange={handleEmailChange}
                        error={!!emailError}
                    />
                </div>
                <div className="mb-4">
                    <Input
                        type="password"
                        color="lightBlue"
                        size="regular"
                        outline={true}
                        placeholder="Password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                </div>
                {isConfirmed && (
                    <div className="mb-4">
                        <Input
                            type="text"
                            color="lightBlue"
                            size="regular"
                            outline={true}
                            placeholder="OTP"
                            value={otp}
                            onChange={(e) => setOtp(e.target.value)}
                        />
                    </div>
                )}
                {/*{loginError && <div className="text-red-500 mb-4">{loginError}</div>}*/}

                {loginError && (
                    <Alert color="red" className="mb-4">
                        <div className="flex items-center justify-between">
                            <div className="flex items-center">
                                <i className="fas fa-info-circle mr-2"></i>
                                <span className="text-sm">{loginError}</span>
                            </div>
                        </div>
                    </Alert>
                )}

                <div className="flex flex-col items-stretch"> {/* 改变布局 */}
                    <Button
                        type="submit"
                        color="lightBlue"
                        buttonType="filled"
                        size="regular"
                        rounded={false}
                        block={false}
                        iconOnly={false}
                        ripple="light"
                        className="mb-4" // 添加边距
                    >
                        Login
                    </Button>
                    <a href="#/" className="text-sm hover:underline">Forgot password?</a>
                </div>
            </Card>
        </form>

        // <form onSubmit={handleSubmit}>
        //     <div
        //         className="relative flex min-h-screen text-gray-800 antialiased flex-col justify-center overflow-hidden bg-gray-50 py-6 sm:py-12">
        //         <div className="relative py-3 sm:w-96 mx-auto text-center">
        //             <span className="text-2xl font-light">Login to your account</span>
        //             <div className="mt-4 bg-white shadow-md rounded-lg text-left">
        //                 <div
        //                     className="text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded"></div>
    //                 <div className="px-8 py-6">
    //                     <label className="block font-semibold mb-2">Email</label>
    //                     <input type="email" value={email} onChange={handleEmailChange}
    //                            className={getInputClass()}/>
    //                     {emailError && <p className="text-red-500">{emailError}</p>}
    //                     <label className="block mt-3 font-semibold mb-2">Password</label>
    //                     <input type="password" value={password} onChange={(e) => setPassword(e.target.value)}
    //                            className="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md"/>
    //                     {isConfirmed && (
    //                         <div>
    //                             <label className="block mt-3 font-semibold mb-2">OTP</label>
    //                             <input type="text" value={otp} onChange={(e) => setOtp(e.target.value)}
    //                                    className="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md"/>
    //                         </div>
    //                     )}
    //                     {loginError && <div className="text-red-500">{loginError}</div>}
    //                     <div className="flex justify-between items-baseline">
    //                         <button type="submit"
    //                                 className="mt-4 text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded">Login
    //                             </button>
    //                             <a href="#/" className="text-sm hover:underline">Forgot password?</a>
    //                         </div>
    //                     </div>
    //                 </div>
    //             </div>
    //         </div>
    //     </form>
    );
}

export default LoginForm;

