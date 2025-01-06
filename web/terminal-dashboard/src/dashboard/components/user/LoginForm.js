// LoginForm.js
import React, {useEffect, useState} from 'react';
import {check2FA, checkEmail, login} from "../../../api/api";
import {Alert, Button, Checkbox, Input, Typography} from "@material-tailwind/react";
import {Link, useNavigate} from "react-router-dom";
import CryptoJS from 'crypto-js';
import {EnvelopeIcon, LockClosedIcon} from "@heroicons/react/24/solid";

function LoginForm({ onLogin }) {
    const [email, setEmail] = React.useState('');
    const [password, setPassword] = React.useState('');
    const [emailError, setEmailError] = useState('');// 添加一个新的状态来保存邮箱错误信息
    const [loginError, setLoginError] = useState(''); // 保存登录错误消息
    const [isConfirmed, setIsConfirmed] = useState(false); // 新增状态变量
    const [otp, setOtp] = useState(''); // OTP 输入框的状态变量
    const navigate = useNavigate();

    useEffect(() => {
        async function checkUser2FA() {
            if (email) { // 只有当 email 不为空时才发送请求
                const response = await check2FA(email);
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
            setEmailError(exists ? '' : 'Email not registered');
        } catch (error) {
            console.error(error);
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        // 验证电子邮件和密码是否已填写
        if (!email || !password) {
            setLoginError('请填写所有必填字段');
            return;
        }
        try {
            // 对密码进行哈希处理
            const hashedPassword = CryptoJS.SHA256(password).toString();
            // 将OTP赋值给totp_Secret
            const data = {
                email: email,
                password: hashedPassword,
                otp: otp
            }
            // await login(data); // 使用 login 函数
            const response = await login(data);
            if (response?.data) {
                const { accessToken, refreshToken, user } = response.data;

                localStorage.setItem('accessToken', accessToken);
                localStorage.setItem('refreshToken', refreshToken);
                localStorage.setItem('user', JSON.stringify(user));

                onLogin(user, refreshToken);
                navigate('/dashboard');
            } else {
                console.error('Login response is invalid:', response);
                setLoginError('服务器响应无效');
            }
        } catch (error) {
            console.error("Error during login:", error);
            setLoginError('用户名或密码错误或者OTP错误');
            setTimeout(() => setLoginError(''), 1000);
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
                            Your Email
                        </Typography>
                        <Input
                            variant="outlined"
                            label="Email"
                            type="email"
                            color="lightBlue"
                            size="regular"
                            outline={true}
                            // placeholder="email"
                            value={email}
                            onChange={handleEmailChange}
                            error={!!emailError}
                            icon={<EnvelopeIcon className="h-5 w-5" />}
                        />
                        <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                            Password
                        </Typography>
                        <Input
                            variant="outlined"
                            label="Password"
                            type="password"
                            color="lightBlue"
                            size="regular"
                            outline={true}
                            // placeholder="Password"
                            value={password}
                            icon={<LockClosedIcon className="h-5 w-5" />}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                        {isConfirmed && (
                            <>
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    OTP
                                </Typography>
                                <Input
                                    variant="outlined"
                                    label="OTP"
                                    type="text"
                                    color="lightBlue"
                                    size="regular"
                                    outline={true}
                                    // placeholder="OTP"
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
                            <a href="/reset-password">
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
                            <img src="/twitter-logo.svg" height={24} width={24} alt=""/>
                            <span>Sign in With Twitter</span>
                        </Button>
                    </div>
                    <Typography variant="paragraph" className="text-center text-blue-gray-500 font-medium mt-4">
                        Not registered?
                        <Link to="/register" className="text-gray-900 ml-1">Create account</Link>
                    </Typography>


                </form>
            </div>
            <div className="w-2/5 h-full hidden lg:block">
                <img src="/pattern.png" className="h-full w-full object-cover rounded-3xl" alt="/"/>
            </div>
        </section>
    );
}

export default LoginForm;

