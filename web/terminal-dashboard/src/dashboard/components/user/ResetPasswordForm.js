import React, {useState, } from "react";
import {Button, Card, Input, Typography} from "@material-tailwind/react";
import {checkEmail, login, resetPassword} from "../../../api/api";
import {useNavigate} from "react-router-dom";


function ResetPasswordForm({ onResetPassword }) {

    const [step, setStep] = useState(1);
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [emailError, setEmailError] = useState('');// 添加一个新的状态来保存邮箱错误信息
    const [ResetPasswordError, setResetPasswordError] = useState(''); // 保存登录错误消息

    const navigate = useNavigate();

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

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (step === 1) {
            const exists = await checkEmail(email);
            if (exists) {
                setStep(2);
            } else {
                alert("Email not exists");
            }
        }
        else if (step === 2) {
            if (password === confirmPassword) {
                try {
                    const data = await resetPassword(email, password);
                    console.log(data)
                    alert("reset Password ✅");
                    navigate('/login'); // 跳转到登录页面
                }catch (error) {
                    setResetPasswordError(error.message);
                }
            } else {
                    alert("Passwords don't match");
                }
            }
    };

    return (
        <div className="flex justify-center items-center h-screen">
            <Card className="w-1/2">
                <div className="text-center">
                    <Typography variant="h2" className="font-bold mb-4">Reset Password</Typography>
                    <Typography variant="paragraph" color="blue-gray" className="text-lg font-normal">Enter your email reset password to Sign In.</Typography>
                </div>
                {/*<form onSubmit="">*/}
                <form onSubmit={handleSubmit}>
                    {step === 1 && (
                        <>
                        <div className="mb-1 flex flex-col gap-6">
                    <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                        Email
                    </Typography>
                    <Input
                        type="email"
                        color="lightBlue"
                        size="regular"
                        outline={true}
                        placeholder="Email"
                        onChange={handleEmailChange}
                        error={!!emailError}
                    />
                        </div>
                        </>
                        )}
                    {step === 2 && (
                        <>
                    <div className="mb-1 flex flex-col gap-6">
                    <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                        Password
                    </Typography>
                    <Input
                        type="password"
                        color="lightBlue"
                        size="regular"
                        outline={true}
                        placeholder="Password"
                        // value={password}
                        // onChange={(e) => setPassword(e.target.value)}
                    />
                    <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                        Confirm Password
                    </Typography>
                    <Input
                        type="password"
                        color="lightBlue"
                        size="regular"
                        outline={true}
                        placeholder="Confirm Password"
                        // value={confirmPassword}
                        // onChange={(e) => setConfirmPassword(e.target.value)}
                    />
                </div>
                        </>
                    )}
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
                    {step === 1 ? 'NEXT' : 'Submit'}
                </Button>
            </form>
            </Card>
        </div>
    )
}

export default ResetPasswordForm;