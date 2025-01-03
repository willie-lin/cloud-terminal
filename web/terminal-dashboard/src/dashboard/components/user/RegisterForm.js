// RegisterForm.js
import React, { useState } from 'react';
import {checkEmail, checkOrganizationName, register} from "../../../api/api";
import {Alert, Button, Checkbox, Input, Typography} from "@material-tailwind/react";
import {Link} from "react-router-dom";

function RegisterForm({ onRegister }) {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [organization, setOrganization] = useState('');
    const [organizationError, setOrganizationError] = useState('')
    const [emailError, setEmailError] = useState('');// 添加一个新的状态来保存邮箱错误信息
    const [passwordError, setPasswordError] = useState(''); // 添加一个新的状态来保存密码错误信息
    const [registerError, setRegisterError] = useState(''); // 添加一个新的状态来保存密码错误信息

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



    const handleOrganizationNameChange = async (e) => {
        const organization = e.target.value;
        setOrganization(organization)
        try {
            const exists = await checkOrganizationName(organization);
            setOrganizationError(exists ? 'Organization already registered' : '');
        } catch (error) {
            console.error(error)
        }
    };

    // const handleEmailChange = (e) => setEmail(e.target.value);
    // const handlePasswordChange = (e) => setPassword(e.target.value);
    // const handleConfirmPasswordChange = (e) => setConfirmPassword(e.target.value);
    // const handleOrganizationNameChange = (e) => setOrganizationName(e.target.value);


    const CryptoJS = require("crypto-js");
    const handleSubmit = async (e) => {
        e.preventDefault();
        // 验证组织名称是否填写
        // if (!organization || !password) {
        //     setRegisterError('请填写所有必填字段');
        //     setTimeout(() => setRegisterError(''), 1000); // 1秒后清除错误信息
        //     return;
        // }
        // 验证电子邮件和密码是否已填写
        if (!email || !password || !organization) {
            setRegisterError('请填写所有必填字段');
            setTimeout(() => setRegisterError(''), 1000); // 1秒后清除错误信息
            return;
        }
        // 验证密码和确认密码是否匹配
        if (password !== confirmPassword) {
            setPasswordError("Passwords don't match"); // 设置密码错误信息
            setTimeout(() => setPasswordError(''), 1000); // 1秒后清除错误信息
            return;
        }
        try {
            // 对密码进行哈希处理
            const hashedPassword = CryptoJS.SHA256(password).toString();
            const data = {
                email: email,
                password: hashedPassword,
                tenant_name: organization, // 将租户名称包含在请求数据中

            }
           await register(data); // 使用 register 函数
            // console.log(data);
            onRegister(email);
        } catch (error) {
            // console.error(error);
            setRegisterError("Registration failed");
            setTimeout(() => setRegisterError(''), 1000);
        }
    };

    return (
        <section className="m-8 flex">
            <div className="w-2/5 h-full hidden lg:block">
                <img src="/pattern.png" className="h-full w-full object-cover rounded-3xl" alt="/"/>
            </div>
            <div className="w-full lg:w-3/5 flex flex-col items-center justify-center">
                <div className="text-center">
                    <Typography variant="h2" className="font-bold mb-4">Join Us Today</Typography>
                    <Typography variant="paragraph" color="blue-gray" className="text-lg font-normal">Enter your email and password to register.</Typography>
                </div>
                <form onSubmit={handleSubmit} className="mt-8 mb-2 mx-auto w-80 max-w-screen-lg lg:w-1/2">
                    <div className="mb-1 flex flex-col gap-6">
                        <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                            Your Organization
                        </Typography>
                        <Input
                            variant="outlined"
                            label="Organization Name"
                            size="lg"
                            type="text"
                            color="lightBlue"
                            outline={true}
                            // placeholder="Email"
                            value={organization}
                            onChange={ handleOrganizationNameChange }
                            error={!!organizationError}
                        />
                        <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                            Your Email
                        </Typography>
                        <Input
                            variant="outlined"
                            label="Email"
                            size="lg"
                            type="email"
                            color="lightBlue"
                            outline={true}
                            // placeholder="Email"
                            value={email}
                            onChange={ handleEmailChange }
                            error={!!emailError}
                        />
                        {/*{organizationError && (*/}
                        {/*    <Alert color="red" className="mb-4">*/}
                        {/*        <div className="flex items-center justify-between">*/}
                        {/*            <div className="flex items-center">*/}
                        {/*                <i className="fas fa-info-circle mr-2"></i>*/}
                        {/*                <span className="text-sm">{organizationError}</span>*/}
                        {/*            </div>*/}
                        {/*        </div>*/}
                        {/*    </Alert>*/}
                        {/*)}*/}
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
                            onChange={(e) => setPassword(e.target.value)}
                        />
                        {registerError && (
                            <Alert color="red" className="mb-4">
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center">
                                        <i className="fas fa-info-circle mr-2"></i>
                                        <span className="text-sm">{registerError}</span>
                                    </div>
                                </div>
                            </Alert>
                        )}

                        <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                            Confirm Password
                        </Typography>
                        <Input
                            variant="outlined"
                            label="Confirm Password"
                            type="password"
                            color="lightBlue"
                            size="regular"
                            outline={true}
                            // placeholder="ConfirmPassword"
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
                        containerProps={{className: "-ml-2.5"}}
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
                        Already have an account?
                        <Link to="/login" className="text-gray-900 ml-1">Sign in</Link>
                    </Typography>
                </form>
            </div>
        </section>
    );
}

export default RegisterForm;
