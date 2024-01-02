import React, {useEffect, useState} from 'react';
import {confirm2FA, enable2FA, getUserByEmail, validate2FA} from "../../../api/api";
import {Button, Typography} from "@material-tailwind/react";

// 自定义Hook，用于处理二次验证的逻辑
function useTwoFactorAuth(email) {
    const [userInfo, setUserInfo] = useState(null);
    const [qrCode, setQrCode] = useState(null);
    const [qrGenerated, setQrGenerated] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    // 获取用户信息
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

    // 生成二次验证的二维码
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

    // 确认二次验证
    // 确认二次验证
    const confirm2FAHandler = async () => {
        setLoading(true);
        try {
            await confirm2FA(email);
            alert("2FA confirmed");
        } catch (error) {
            console.error('Error:', error);
            setError('确认二次验证失败');
        } finally {
            setLoading(false);
        }
    };
    return { userInfo, qrCode, qrGenerated, loading, error, generateQRCode, confirm2FAHandler };
}
function TwoFactorAuthPage({ email }) {
    const { userInfo, qrCode, qrGenerated, loading, error, generateQRCode, confirm2FAHandler } = useTwoFactorAuth(email);

    return (
        <div>
            <div className="flex flex-col items-center justify-center h-screen">
                {error && <p>{error}</p>}
                {loading && <p>加载中...</p>}
                {userInfo && <Typography color="blue-gray" className="font-medium" textGradient>
                    用户名: {userInfo.email}。
                </Typography>}
                <p>你还没有开启二次验证，开启二次验证可以提高账户的安全性。</p>
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
                {qrGenerated && qrCode && <img src={`data:image/png;base64,${qrCode}`} alt="二次验证二维码"/>}
                {qrGenerated && qrCode && <Button
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
            </div>
        </div>
    );
}
export default TwoFactorAuthPage;