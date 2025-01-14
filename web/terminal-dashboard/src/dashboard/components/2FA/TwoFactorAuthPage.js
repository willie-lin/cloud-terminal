import React, {useContext, useEffect, useState} from 'react';
import {check2FA, confirm2FA, enable2FA, getUserByEmail} from "../../../api/api";
import {Alert, Button, Card, CardFooter, Input, Typography} from "@material-tailwind/react";
import {AuthContext} from "../../../App";
import {LockOpenIcon, QrCodeIcon, ShieldCheckIcon} from "@heroicons/react/16/solid";

// 自定义Hook，用于处理二次验证的逻辑
function useTwoFactorAuth() {

    const { currentUser } = useContext(AuthContext);
    const email = currentUser.email
    const [userInfo, setUserInfo] = useState(null);
    const [qrCode, setQrCode] = useState(null);
    const [otp, setOtp] = useState(''); // 新增状态变量
    const [qrGenerated, setQrGenerated] = useState(false);
    const [isConfirmed, setIsConfirmed] = useState(false); // 新增状态变量
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    // 在你的组件的状态中添加一个新的状态变量来保存tempSecret
    const [secret, setSecret] = useState(null);

    // 检测是否开启2FA
    useEffect(() => {
        async function checkUser2FA() {
            // const email = 'test2@example.com'; // 替换为你需要检查的邮箱
            const response = await check2FA(email);
            console.log(response); // 在控制台打印响应

            // 根据响应设置 isConfirmed 的值
            if (response && response.isConfirmed !== undefined) {
                setIsConfirmed(response.isConfirmed);
            }
        }
        checkUser2FA();
    }, ); // 空数组作为依赖，意味着这个 useEffect 只会在组件挂载时运行一次

     // 检测 用户是否存在
    useEffect(() => {
        if (email) {
            getUserByEmail(email)
                .then(data => setUserInfo(data))
                .catch(error => {
                    console.error('Error:', error);
                    setError('获取用户信息失败');
                });
        }
    }, [email]);

    // 生成二维码
    const generateQRCode = async () => {
        setLoading(true);
        try {
            const data = await enable2FA(email);
            setQrCode(data.qrCode);
            setSecret(data.secret); // 保存tempSecret
            setQrGenerated(true);
        } catch (error) {
            console.error('Error:', error);
            setError('生成二次验证二维码失败');
        } finally {
            setLoading(false);
        }
    };
    // 扫描二维码，进行绑定
    const confirm2FAHandler = async () => {
        setLoading(true);
        try {
            // 将OTP赋值给totp_Secret
            // 将电子邮件地址和totp_Secret传递给后端
            const data = {
                email: email,
                otp: otp,
                secret: secret
            }
            await confirm2FA(data);
            setIsConfirmed(true); // 设置状态变量为true
        } catch (error) {
            console.error('Error:', error);
            setError('Secondary verification failed！');
            setIsConfirmed(false); // 在这里添加这行代码
        } finally {
            setLoading(false);
        }
    };
    return { userInfo, qrCode, qrGenerated, isConfirmed, loading, error, otp, setOtp, generateQRCode, confirm2FAHandler };
}

function TwoFactorAuthPage({ user }) {
    const { userInfo, qrCode, qrGenerated, isConfirmed, loading, error, otp, setOtp, generateQRCode, confirm2FAHandler } = useTwoFactorAuth(user.email);
    return (
        <div className="w-full ">
            {/*<div className="w-full flex flex-col items-center lg:w-3/5 mt-12 bg-gray-100 p-4">*/}
                <Card className="w-full shadow-lg rounded-lg">
            {/*<Card className="w-full max-w-md p-6 shadow-lg rounded-lg">*/}
                <div className="text-center mb-6">
                    <Typography variant="h2" className="font-bold mb-4">MFA Authentication</Typography>
                    <Typography variant="paragraph" color="blue-gray" className="text-lg font-normal">
                        {isConfirmed ? "Your account is protected with 2FA" : "Enhance your account security with 2FA"}
                    </Typography>
                </div>

                <div className="p-6">
                    {error && (
                        <Alert color="red" className="mb-4" open={true}>
                            {error}
                        </Alert>
                    )}
                    {loading && (
                        <Alert color="blue-gray" className="mb-4" open={true}>
                            Loading...
                        </Alert>
                    )}
                    {userInfo && (
                        <Typography variant="h6" color="blue-gray" className="mb-4">
                            Username: {userInfo.email}
                        </Typography>
                    )}
                    {!loading && (
                        isConfirmed ? (
                            <Alert icon={<ShieldCheckIcon className="h-6 w-6"/>} className="mb-2" open={true}>
                                ✅ Two-factor authentication is enabled!
                            </Alert>
                        ) : (
                            <>
                                {!qrGenerated && (
                                    <Typography color="gray" className="mb-3">
                                        You haven't turned on MFA authentication, which improves the security of your
                                        account.
                                    </Typography>
                                )}
                                {!qrGenerated && (
                                    <Button
                                        type="button"
                                        color="lightBlue"
                                        onClick={generateQRCode}
                                        className="flex items-center justify-center gap-2 mt-2 w-full"
                                    >
                                        <QrCodeIcon className="h-5 w-5"/>
                                        Generate MFA QR Code
                                    </Button>
                                )}
                                {qrGenerated && qrCode && (
                                    <div className="flex flex-col items-center gap-4 mt-4">
                                        <img
                                            src={`data:image/png;base64,${qrCode}`}
                                            alt="Two-factor authentication QR code"
                                            className="w-48 h-48"
                                        />
                                        <Input
                                            type="text"
                                            value={otp}
                                            onChange={(e) => setOtp(e.target.value)}
                                            placeholder="Enter your one-time password"
                                            className="!border !border-gray-300 bg-white text-gray-900 shadow-lg ring-transparent placeholder:text-gray-500 focus:!border-gray-900 focus:ring-4 focus:ring-gray-900"
                                            labelProps={{className: "hidden"}}
                                            containerProps={{className: "min-w-[100px]"}}
                                        />
                                        <Button
                                            type="button"
                                            color="black"
                                            onClick={confirm2FAHandler}
                                            className="flex items-center justify-center gap-2 w-full mt-4"
                                        >
                                            <LockOpenIcon className="h-5 w-5"/>
                                            Confirm MFA Authentication
                                        </Button>
                                    </div>
                                )}
                            </>
                        )
                    )}
                </div>

                <CardFooter className="text-center pt-0">
                    <Typography variant="small" className="mt-6 flex justify-center">
                        Need Help?
                        <Typography
                            as="a"
                            href="#"
                            variant="small"
                            color="lightBlue"
                            className="ml-1 font-bold"
                        >
                            Contact Support.
                        </Typography>
                    </Typography>
                </CardFooter>
            </Card>
        </div>
    );
}

export default TwoFactorAuthPage;