import React, {useEffect, useState} from 'react';
import {confirm2FA, enable2FA, getUserByEmail} from "../../../api/api";
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
            const totp_Secret = otp
            // 将电子邮件地址和totp_Secret传递给后端
            await confirm2FA({ email, totp_Secret });
            setIsConfirmed(true); // 设置状态变量为true
        } catch (error) {
            console.error('Error:', error);
            setError('确认二次验证失败');
            setIsConfirmed(false);
        } finally {
            setLoading(false);
        }
    };
    return { userInfo, qrCode, qrGenerated, isConfirmed, loading, error, otp, setOtp, generateQRCode, confirm2FAHandler };
}
function TwoFactorAuthPage({ email }) {
    const { userInfo, qrCode, qrGenerated, isConfirmed, loading, error, otp, setOtp, generateQRCode, confirm2FAHandler } = useTwoFactorAuth(email);
    return (
        <div>
            <div className="flex flex-col items-center justify-center h-screen">
                {error && <p>{error}</p>}
                {loading && <p>加载中...</p>}
                {userInfo && <Typography color="blue-gray" className="font-medium" textGradient>
                    用户名: {userInfo.email}。
                </Typography>}
                {!qrGenerated && <p>你还没有开启二次验证，开启二次验证可以提高账户的安全性。</p>}
                {!qrGenerated && <Button
                    color="lightBlue"
                    buttonType="filled"
                    size="regular"
                    rounded={false}
                    block={false}
                    iconOnly={false}
                    ripple="light"
                    onClick={generateQRCode}
                >
                    生成二次验证二维码
                </Button>}
                {qrGenerated && !isConfirmed && qrCode && <img src={`data:image/png;base64,${qrCode}`} alt="二次验证二维码"/>}
                {qrGenerated && !isConfirmed && qrCode && <input
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
                    }}
                />}
                {qrGenerated && !isConfirmed && qrCode && <Button
                    color="lightBlue"
                    buttonType="filled"
                    size="regular"
                    rounded={false}
                    block={false}
                    iconOnly={false}
                    ripple="light"
                    onClick={confirm2FAHandler}
                >
                    确认二次验证
                </Button>}
                {isConfirmed && <p>✅ 二次验证已成功绑定！</p>}
            </div>
        </div>
    );
}
export default TwoFactorAuthPage;