import React, {useEffect, useState} from 'react';
import {check2FA, confirm2FA, enable2FA, getUserByEmail} from "../../../api/api";
import {Button, Typography} from "@material-tailwind/react";

// 自定义Hook，用于处理二次验证的逻辑
function useTwoFactorAuth(email) {
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
    }, []); // 空数组作为依赖，意味着这个 useEffect 只会在组件挂载时运行一次

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
            await confirm2FA(email, otp, secret);
            setIsConfirmed(true); // 设置状态变量为true
        } catch (error) {
            console.error('Error:', error);
            setError('确认二次验证失败');
            setIsConfirmed(false); // 在这里添加这行代码
        } finally {
            setLoading(false);
        }
    };
    return { userInfo, qrCode, qrGenerated, isConfirmed, loading, error, otp, setOtp, generateQRCode, confirm2FAHandler };
}

function TwoFactorAuthPage({ email }) {
    const { userInfo, qrCode, qrGenerated, isConfirmed, loading, error, otp, setOtp, generateQRCode, confirm2FAHandler } = useTwoFactorAuth(email);
    return (
        <div className="flex flex-col items-center justify-center flex-grow bg-gray-200">
            <div
                className="p-6 mt-6 text-center border w-96 rounded-xl hover:shadow-xl transition-shadow duration-300 ease-in-out">
                {error && <p className="text-red-500">{error}</p>}
                {loading && <p className="text-blue-500">加载中...</p>}
                {userInfo && <Typography color="blue-gray" className="font-medium text-black textGradient">
                    用户名: {userInfo.email}。
                </Typography>}
                {!loading && (
                    isConfirmed ? (
                        <p className="text-green-500 text-lg">✅ 已开启二次认证防护！</p>
                    ) : (
                        <>
                            {!qrGenerated &&
                                <p className="text-black text-lg">你还没有开启二次验证，开启二次验证可以提高账户的安全性。</p>}
                            {!qrGenerated && <Button
                                color="lightBlue"
                                buttonType="filled"
                                size="regular"
                                rounded={false}
                                block={false}
                                iconOnly={false}
                                ripple="light"
                                onClick={generateQRCode}
                                className="mt-4"
                            >
                                生成二次验证二维码
                            </Button>}
                            {qrGenerated && qrCode &&
                                <img src={`data:image/png;base64,${qrCode}`} alt="二次验证二维码"
                                     className="mt-4 mx-auto"/>}
                            {qrGenerated && qrCode && <input
                                type="text"
                                value={otp}
                                onChange={e => setOtp(e.target.value)}
                                placeholder="请输入你的一次性密码"
                                style={{
                                    border: '1px solid lightBlue', // 添加边框
                                    backgroundColor: '#f8f9fa', // 改变背景颜色
                                    width: '250px', // 增加宽度
                                    padding: '5px', // 添加内边距
                                    borderRadius: '5px', // 添加边框圆角
                                    marginBottom: '20px', // 增加下边距
                                    textAlign: 'center', // 文字居中
                                }}
                                className="mt-4 mx-auto"
                            />}
                            {qrGenerated && qrCode && <Button
                                color="lightBlue"
                                buttonType="filled"
                                size="regular"
                                rounded={false}
                                block={false}
                                iconOnly={false}
                                ripple="light"
                                onClick={confirm2FAHandler}
                                className="mt-4 mx-auto"
                                style={{width: '250px'}}
                            >
                                确认二次验证
                            </Button>}
                        </>
                    )
                )}
            </div>
        </div>
    );
}

export default TwoFactorAuthPage;