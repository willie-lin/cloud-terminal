// RegisterForm.js
import React, { useState } from 'react';
import {checkEmail, register} from "../../../api/api";
import {Alert, Button, Card, Checkbox, Input, Typography} from "@material-tailwind/react";

function RegisterForm({ onRegister }) {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [emailError, setEmailError] = useState('');// 添加一个新的状态来保存邮箱错误信息
    const [passwordError, setPasswordError] = useState(''); // 添加一个新的状态来保存密码错误信息
    const [loginError, setLoginError] = useState(''); // 添加一个新的状态来保存密码错误信息

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
        // 验证电子邮件和密码是否已填写
        if (!email || !password) {
            setLoginError('请填写所有必填字段');
            setTimeout(() => setLoginError(''), 1000); // 1秒后清除错误信息
            return;
        }
        // 验证密码和确认密码是否匹配
        if (password !== confirmPassword) {
            setPasswordError("Passwords don't match"); // 设置密码错误信息
            setTimeout(() => setPasswordError(''), 1000); // 1秒后清除错误信息
            return;
        }
        try {
            const data = await register(email, password); // 使用 register 函数
            console.log(data);
            onRegister(email);
        } catch (error) {
            console.error(error);
        }
    };

    return (
        <section className="m-8 flex">
            <div className="w-2/5 h-full hidden lg:block">
                <img src="/img/pattern.png" className="h-full w-full object-cover rounded-3xl" alt="/"/>
            </div>
            <div className="w-full lg:w-3/5 flex flex-col items-center justify-center">
                <div className="text-center">
                    <Typography variant="h2" className="font-bold mb-4">Join Us Today</Typography>
                    <Typography variant="paragraph" color="blue-gray" className="text-lg font-normal">Enter your email and password to register.</Typography>
                </div>
                <form onSubmit={handleSubmit} className="mt-8 mb-2 mx-auto w-80 max-w-screen-lg lg:w-1/2">
                    <div className="mb-1 flex flex-col gap-6">
                        <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                            Your email
                        </Typography>
                        <Input
                            size="lg"
                            placeholder="name@mail.com"
                            className=" !border-t-blue-gray-200 focus:!border-t-gray-900"
                            labelProps={{
                                className: "before:content-none after:content-none",
                            }}
                            type="email"
                            color="lightBlue"
                            outline={true}
                            value={email}
                            onChange={handleEmailChange}
                            error={!!emailError}
                        />
                        <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                            Password
                        </Typography>
                        <Input
                            type="password"
                            color="lightBlue"
                            size="regular"
                            outline={true}
                            placeholder="Password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
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

                        <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                            Confirm Password
                        </Typography>
                        <Input
                            type="password"
                            color="lightBlue"
                            size="regular"
                            outline={true}
                            placeholder="Confirm Password"
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                        />
                    </div>
                    {passwordError && (
                        <Alert color="red" className="mb-4">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center">
                                    <i className="fas fa-info-circle mr-2"></i>
                                    <span className="text-sm">{passwordError}</span>
                                </div>
                            </div>
                        </Alert>
                    )}
                    <Checkbox
                        label={
                            <Typography
                                variant="small"
                                color="gray"
                                className="flex items-center justify-start font-medium"
                            >
                                I agree the&nbsp;
                                <a
                                    href="#/"
                                    className="font-normal text-black transition-colors hover:text-gray-900 underline"
                                >
                                    Terms and Conditions
                                </a>
                            </Typography>
                        }
                        containerProps={{ className: "-ml-2.5" }}
                    />

                    <Button fullWidth
                            type="submit"
                            color="lightBlue"
                            buttonType="filled"
                            size="regular"
                            rounded={false}
                            block={false}
                            iconOnly={false}
                            ripple="light"
                            className="mt-6" // 添加边距
                    >
                        Register Now
                    </Button>
            </form>
            </div>
        </section>


// <form onSubmit={handleSubmit} className="flex flex-col items-center justify-center min-h-screen">
//         <Card className="w-full max-w-md"> {/* 设置最大宽度 */}
//         <h6 className="text-gray-500 text-lg text-center">Register for an account</h6>
//         <hr className="mb-6 border-b-1 border-gray-300"/>
//         <div className="mb-4">
//                     <Input
//                         type="email"
//                         color="lightBlue"
//                         size="regular"
//                         outline={true}
//                         placeholder="Email"
//                         value={email}
//                         onChange={handleEmailChange}
//                         error={!!emailError}
//                     />
//                 </div>
//                 <div className="mb-4">
//                     <Input
//                         type="password"
//                         color="lightBlue"
//                         size="regular"
//                         outline={true}
//                         placeholder="Password"
//                         value={password}
//                         onChange={(e) => setPassword(e.target.value)}
//                     />
//                 </div>
//                 <div className="mb-4">
//                     <Input
//                         type="password"
//                         color="lightBlue"
//                         size="regular"
//                         outline={true}
//                         placeholder="Confirm Password"
//                         value={confirmPassword}
//                         onChange={(e) => setConfirmPassword(e.target.value)}
//                     />
//                 </div>
//                 {passwordError && (
//                     <Alert color="red" className="mb-4">
//                         <div className="flex items-center justify-between">
//                             <div className="flex items-center">
//                                 <i className="fas fa-info-circle mr-2"></i>
//                                 <span className="text-sm">{passwordError}</span>
//                             </div>
//                         </div>
//                     </Alert>
//                 )}
//                 <div className="flex flex-col items-stretch"> {/* 改变布局 */}
//                     <Button
//                         type="submit"
//                         color="lightBlue"
//                         buttonType="filled"
//                         size="regular"
//                         rounded={false}
//                         block={false}
//                         iconOnly={false}
//                         ripple="light"
//                         className="mb-4" // 添加边距
//                     >
//                         Register
//                     </Button>
//                 </div>
//             </Card>
//         </form>


        // <form onSubmit={handleSubmit}>
        // <div className="relative flex min-h-screen text-gray-800 antialiased flex-col justify-center overflow-hidden bg-gray-50 py-6 sm:py-12">
        //     <div className="relative py-3 sm:w-96 mx-auto text-center">
        //         <span className="text-2xl font-light">Register for an account</span>
        //         <div className="mt-4 bg-white shadow-md rounded-lg text-left">
        //             {/*<div className="h-2 bg-purple-600 rounded-t-md"></div>*/}
        //             <div className="text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded"></div>
        //             <div className="px-8 py-6">
        //                 <label className="block font-semibold mb-2">Email</label>
        //                 <input type="email" value={email} onChange={handleEmailChange}
        //                        className={getInputClass()}/>
        //                 {emailError && <p className="text-red-500">{emailError}</p>}
        //                 <label className="block mt-3 font-semibold mb-2">Password</label>
        //                 <input type="password" value={password} onChange={(e) => setPassword(e.target.value)}
        //                        className="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md"/>
        //                 <label className="block mt-3 font-semibold mb-2">Confirm Password</label>
        //                 <input type="password" value={confirmPassword}
        //                        onChange={(e) => setConfirmPassword(e.target.value)}
        //                        className="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md"/>
        //                 <div className="flex justify-between items-baseline">
        //                     {/*<button type="submit" className="mt-4 bg-purple-600 text-white py-2 px-6 rounded-md hover:bg-purple-700">Register</button>*/}
        //                     <button type="submit"
        //                             className="mt-4 text-sm font-semibold leading-6 text-white bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded">Register
        //                     </button>
        //                 </div>
        //             </div>
        //         </div>
        //     </div>
        // </div>
        // </form>
    );
}

export default RegisterForm;
