import React, { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import {
    Typography,
    Input,
    Button,
    Checkbox,
    Alert,
} from "@material-tailwind/react";
import { EnvelopeIcon, LockClosedIcon } from "@heroicons/react/24/solid";
import {checkEmail, resetPassword} from "../../../api/api";

const CryptoJS = require("crypto-js");

export function ResetPasswordForm({ onResetPassword }) {
    const [step, setStep] = useState(1);
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [emailError, setEmailError] = useState("");
    const [resetPasswordError, setResetPasswordError] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const [agreeTerms, setAgreeTerms] = useState(false);

    const navigate = useNavigate();

    const handleEmailChange = async (e) => {
        const email = e.target.value;
        setEmail(email);
        try {
            const exists = await checkEmail(email);
            setEmailError(exists ? "" : "Email not registered");
        } catch (error) {
            console.error(error);
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        setResetPasswordError("");

        try {
            if (step === 1) {
                const exists = await checkEmail(email);
                if (exists) {
                    setStep(2);
                } else {
                    setEmailError("Email not registered");
                }
            } else if (step === 2) {
                if (password === confirmPassword) {
                    const hashedPassword = CryptoJS.SHA256(password).toString();
                    // await onResetPassword({
                    //     email,
                    //     password: hashedPassword,
                    // });
                    const data = {
                        email: email,
                        password: hashedPassword,
                    }
                    await resetPassword(data);
                    console.log(data)
                    alert("reset Password ✅");
                    navigate('/login'); // 跳转到登录页面
                } else {
                    setResetPasswordError("Passwords don't match");
                }
            }
        } catch (error) {
            setResetPasswordError(error.message);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <section className="m-8 flex gap-4">
            <div className="w-full lg:w-3/5 mt-24">
                <div className="text-center">
                    <Typography variant="h2" className="font-bold mb-4">Reset Password</Typography>
                    <Typography variant="paragraph" color="blue-gray" className="text-lg font-normal">
                        {step === 1 ? "Enter your email to reset your password." : "Enter your new password."}
                    </Typography>
                </div>
                <form onSubmit={handleSubmit} className="mt-8 mb-2 mx-auto w-80 max-w-screen-lg lg:w-1/2">
                    <div className="mb-1 flex flex-col gap-6">
                        {step === 1 ? (
                            <>
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    Your Email
                                </Typography>
                                <Input
                                    type="email"
                                    label="Email"
                                    value={email}
                                    onChange={handleEmailChange}
                                    error={!!emailError}
                                    icon={<EnvelopeIcon className="h-5 w-5" />}
                                />
                                {emailError && (
                                    <Typography variant="small" color="red" className="flex items-center gap-1 font-normal">
                                        {emailError}
                                    </Typography>
                                )}
                            </>
                        ) : (
                            <>
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    New Password
                                </Typography>
                                <Input
                                    type="password"
                                    label="New Password"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    icon={<LockClosedIcon className="h-5 w-5" />}
                                />
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    Confirm Password
                                </Typography>
                                <Input
                                    type="password"
                                    label="Confirm Password"
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    icon={<LockClosedIcon className="h-5 w-5" />}
                                />
                            </>
                        )}

                        {resetPasswordError && (
                            <Alert color="red" className="mb-4">
                                <div className="flex items-center">
                                    <span className="text-sm">{resetPasswordError}</span>
                                </div>
                            </Alert>
                        )}
                    </div>

                    {step === 2 && (
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
                            containerProps={{ className: "-ml-2.5" }}
                            checked={agreeTerms}
                            onChange={(e) => setAgreeTerms(e.target.checked)}
                        />
                    )}

                    <Button
                        type="submit"
                        color="lightBlue"
                        className="mt-6"
                        fullWidth
                        disabled={isLoading || (step === 2 && (!password || !confirmPassword || !agreeTerms))}
                    >
                        {isLoading ? "Please wait..." : (step === 1 ? "Next Step" : "Reset Password")}
                    </Button>

                    <Typography variant="paragraph" className="text-center text-blue-gray-500 font-medium mt-4">
                        Remember your password?
                        <Link to="/login" className="text-gray-900 ml-1">Sign in</Link>
                    </Typography>
                </form>
            </div>
            <div className="w-2/5 h-full hidden lg:block">
                <img src="/pattern.png" className="h-full w-full object-cover rounded-3xl" alt=""/>
            </div>
        </section>
    );
}


export default ResetPasswordForm;


