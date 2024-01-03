// LoginForm.js
import React, {useEffect, useState} from 'react';
import {check2FA, checkEmail, login} from "../../../api/api";
import {Alert, Button, Checkbox, Input, Typography} from "@material-tailwind/react";
import {Link} from "react-router-dom";

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
    // const getInputClass = () => {
    //     if (email.length === 0) {
    //         return "border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-indigo-500 focus:ring-1 rounded-md";
    //     }
    //     // return emailError ? "border-green-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-green-500 focus:ring-1 rounded-md" : "border-red-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-red-500 focus:ring-1 rounded-md";
    //     return emailError ? "border-red-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-red-500 focus:ring-1 rounded-md" : "border-green-500 w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-green-500 focus:ring-1 rounded-md";
    // };
    const handleSubmit = async (e) => {
        e.preventDefault();
        // 验证电子邮件和密码是否已填写
        if (!email || !password) {
            setLoginError('请填写所有必填字段');
            return;
        }
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
        <section className="m-8 flex gap-4">
            <div className="w-full lg:w-3/5 mt-24">
                <div className="text-center">
                    <Typography variant="h2" className="font-bold mb-4">Sign In</Typography>
                    <Typography variant="paragraph" color="blue-gray" className="text-lg font-normal">Enter your email
                        and password to Sign In.</Typography>
                </div>
                <form onSubmit={handleSubmit} className="mt-8 mb-2 mx-auto w-80 max-w-screen-lg lg:w-1/2">
                    <div className="mb-1 flex flex-col gap-6">
                        <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                            Your email
                        </Typography>
                        <Input
                            type="email"
                            color="lightBlue"
                            size="regular"
                            outline={true}
                            placeholder="name@mail.com"
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
                            placeholder="********"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                        {isConfirmed && (
                            <>
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    OTP
                                </Typography>
                                <Input
                                    type="text"
                                    color="lightBlue"
                                    size="regular"
                                    outline={true}
                                    placeholder="OTP"
                                    value={otp}
                                    onChange={(e) => setOtp(e.target.value)}
                                />
                            </>
                        )}

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
                    </div>
                    <Checkbox
                        label={
                            <Typography
                                variant="small"
                                color="gray"
                                className="flex items-center justify-start font-medium"
                            >
                                I agree the&nbsp;
                                <a href="#/" className="font-normal text-black transition-colors hover:text-gray-900 underline">
                                    Terms and Conditions
                                </a>
                            </Typography>
                        }
                        containerProps={{className: "-ml-2.5"}}
                    />

                    <Button
                        fullWidth
                        type="submit"
                        color="lightBlue"
                        buttonType="filled"
                        size="regular"
                        rounded={false}
                        block={false}
                        iconOnly={false}
                        ripple="light"
                        className="mt-6"
                    >
                        Sign In
                    </Button>

                    <div className="flex items-center justify-between gap-2 mt-6">
                        <Checkbox
                            label={
                                <Typography
                                    variant="small"
                                    color="gray"
                                    className="flex items-center justify-start font-medium"
                                >
                                    Subscribe me to newsletter
                                </Typography>
                            }
                            containerProps={{className: "-ml-2.5"}}
                        />
                        <Typography variant="small" className="font-medium text-gray-900">
                            <a href="#/">
                                Forgot Password
                            </a>
                        </Typography>
                    </div>

                    <div className="space-y-4 mt-8">
                        <Button size="lg" color="white" className="flex items-center gap-2 justify-center shadow-md"
                                fullWidth>
                            <svg width="17" height="16" viewBox="0 0 17 16" fill="none"
                                 xmlns="http://www.w3.org/2000/svg">
                                <g clipPath="url(#clip0_1156_824)">
                                    <path
                                        d="M16.3442 8.18429C16.3442 7.64047 16.3001 7.09371 16.206 6.55872H8.66016V9.63937H12.9813C12.802 10.6329 12.2258 11.5119 11.3822 12.0704V14.0693H13.9602C15.4741 12.6759 16.3442 10.6182 16.3442 8.18429Z"
                                        fill="#4285F4"/>
                                    <path
                                        d="M8.65974 16.0006C10.8174 16.0006 12.637 15.2922 13.9627 14.0693L11.3847 12.0704C10.6675 12.5584 9.7415 12.8347 8.66268 12.8347C6.5756 12.8347 4.80598 11.4266 4.17104 9.53357H1.51074V11.5942C2.86882 14.2956 5.63494 16.0006 8.65974 16.0006Z"
                                        fill="#34A853"/>
                                    <path
                                        d="M4.16852 9.53356C3.83341 8.53999 3.83341 7.46411 4.16852 6.47054V4.40991H1.51116C0.376489 6.67043 0.376489 9.33367 1.51116 11.5942L4.16852 9.53356Z"
                                        fill="#FBBC04"/>
                                    <path
                                        d="M8.65974 3.16644C9.80029 3.1488 10.9026 3.57798 11.7286 4.36578L14.0127 2.08174C12.5664 0.72367 10.6469 -0.0229773 8.65974 0.000539111C5.63494 0.000539111 2.86882 1.70548 1.51074 4.40987L4.1681 6.4705C4.8001 4.57449 6.57266 3.16644 8.65974 3.16644Z"
                                        fill="#EA4335"/>
                                </g>
                                <defs>
                                    <clipPath id="clip0_1156_824">
                                        <rect width="16" height="16" fill="white" transform="translate(0.5)"/>
                                    </clipPath>
                                </defs>
                            </svg>
                            <span>Sign in With Google</span>
                        </Button>
                        <Button size="lg" color="white" className="flex items-center gap-2 justify-center shadow-md"
                                fullWidth>
                            <img src="/img/twitter-logo.svg" height={24} width={24} alt=""/>
                            <span>Sign in With Twitter</span>
                        </Button>
                    </div>
                    <Typography variant="paragraph" className="text-center text-blue-gray-500 font-medium mt-4">
                        Not registered?
                        <Link to="/auth/sign-up" className="text-gray-900 ml-1">Create account</Link>
                    </Typography>


                </form>
            </div>
            <div className="w-2/5 h-full hidden lg:block">
                <img src="/img/pattern.png" className="h-full w-full object-cover rounded-3xl" alt="/"/>
            </div>
        </section>


// <form onSubmit={handleSubmit} className="flex flex-col items-center justify-center min-h-screen">
//                         <Card className="w-full max-w-md"> {/* 设置最大宽度 */}
//                             <h6 className="text-gray-500 text-lg text-center">Login to your account</h6>
//                             <hr className="mb-6 border-b-1 border-gray-300"/>
//                             <div className="mb-4">
//                                 <Input
//                                     type="email"
//                                     color="lightBlue"
//                                     size="regular"
//                                     outline={true}
//                                     placeholder="Email"
//                                     value={email}
//                                     onChange={handleEmailChange}
//                                     error={!!emailError}
//                                 />
//                             </div>
//                             <div className="mb-4">
//                                 <Input
//                                     type="password"
//                                     color="lightBlue"
//                                     size="regular"
//                                     outline={true}
//                                     placeholder="Password"
//                                     value={password}
//                                     onChange={(e) => setPassword(e.target.value)}
//                                 />
//                             </div>
//                             {isConfirmed && (
//                                 <div className="mb-4">
//                                     <Input
//                                         type="text"
//                                         color="lightBlue"
//                                         size="regular"
//                                         outline={true}
//                                         placeholder="OTP"
//                                         value={otp}
//                                         onChange={(e) => setOtp(e.target.value)}
//                                     />
//                                 </div>
//                             )}
//                             {loginError && <div className="text-red-500 mb-4">{loginError}</div>}
//
//                             {loginError && (
//                                 <Alert color="red" className="mb-4">
//                                     <div className="flex items-center justify-between">
//                                         <div className="flex items-center">
//                                             <i className="fas fa-info-circle mr-2"></i>
//                                             <span className="text-sm">{loginError}</span>
//                                         </div>
//                                     </div>
//                                 </Alert>
//                             )}
//                             <div className="flex flex-col items-stretch"> {/* 改变布局 */}
//                                 <Button
//                                     type="submit"
//                                     color="lightBlue"
//                                     buttonType="filled"
//                                     size="regular"
//                                     rounded={false}
//                                     block={false}
//                                     iconOnly={false}
//                                     ripple="light"
//                                     className="mb-4" // 添加边距
//                                 >
//                                     Login
//                                 </Button>
//                                 <a href="#/" className="text-sm hover:underline">Forgot password?</a>
//                             </div>
//                         </Card>
//                 </form>


//
// {/*//     <form onSubmit={handleSubmit}>*/}
        // <div
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

