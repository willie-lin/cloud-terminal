import React, {useEffect, useState} from 'react';
import {confirm2FA, enable2FA, getUserByEmail, validate2FA} from "../../../api/api";
import {Button, Typography} from "@material-tailwind/react";

function TwoFactorAuthPage({ email }) {
    const [userInfo, setUserInfo] = useState(null);
    const [qrCode, setQrCode] = useState(null);
    const [secret, setSecret] = useState(null);

    // 获取用户信息
    useEffect(() => {
        if (email) {
            getUserByEmail(email)
                .then(data => setUserInfo(data))
                .catch(error => console.error('Error:', error));
        }
    }, [email]);

    // 生成二次验证的二维码
    const generateQRCode = async () => {
        try {
            const data = await enable2FA(email);
            setQrCode(data.qrCode);
            setSecret(data.secret);
        } catch (error) {
            console.error('Error:', error);
        }
    };

    // 确认二次验证
    // 确认二次验证
    const confirm2FAHandler = async () => {
        try {
            const user = { email: email, TotpSecret: secret };
            await confirm2FA(user);
            alert("2FA confirmed");
        } catch (error) {
            console.error('Error:', error);
        }
    };


    return (
        <div>
            <div className="flex flex-col items-center justify-center h-screen">
                {userInfo && <Typography color="blue-gray" className="font-medium" textGradient>
                    用户名: {userInfo.email}。
                </Typography>}
                <p>你还没有开启二次验证，开启二次验证可以提高账户的安全性。</p>
                <Button
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
                </Button>
                {qrCode && <img src={`data:image/png;base64,${qrCode}`} alt="二次验证二维码"/>}
                {qrCode && <Button
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


// function TwoFactorAuthPage(email) {
//     const [userInfo, setUserInfo] = useState(null);
//     // 获取用户信息
//     useEffect(() => {
//         if (email) {  // 添加这一行来检查email是否存在
//             getUserByEmail(email)
//                 .then(data => setUserInfo(data))
//                 .catch(error => console.error('Error:', error));
//         }
//     }, [email]);


// const [step, setStep] = useState(1);
// const [token, setToken] = useState('');
// const [qrCode, setQrCode] = useState(null);
//
// const handleNextStep = async () => {
//     if (step === 1) {
//         // 请求后端生成二维码
//         const data = await enable2FA();
//         setQrCode(data);
//         setStep(2);
//     } else if (step === 2) {
//         setStep(3);
//     } else if (step === 3) {
//         // 验证动态令牌
//         const data = await validate2FA(token);
//         if (data.status === 'valid') {
//             setStep(4);
//         } else {
//             alert('动态令牌无效，请重试。');
//         }
//     }
// };

// return (
//
//     <div className="flex flex-col items-center justify-center h-screen">
//         {userInfo && <Typography color="blue-gray" className="font-medium" textGradient>
//             用户名: {userInfo.email}。
//         </Typography>}
//         {/*<h1 className="mb-4 text-2xl">用户邮箱：{userInfo.email}</h1>*/}
//         <p>你还没有开启二次验证，开启二次验证可以提高账户的安全性。</p>
//     </div>
// <div className="flex flex-col items-center justify-center h-screen">
//     {step === 1 && (
//         <>
//             <h1 className="mb-4 text-2xl">开启二次验证</h1>
//             <Button
//                 color="lightBlue"
//                 buttonType="filled"
//                 size="regular"
//                 rounded={false}
//                 block={false}
        //                 iconOnly={false}
        //                 ripple="light"
        //                 // onClick={handleNextStep}
        //             >
        //                 开启
        //             </Button>
        //         </>
        //     )}
        //     {step === 2 && qrCode && (
        //         <>
        //             <h1 className="mb-4 text-2xl">扫描二维码</h1>
        //             <img src={qrCode} alt="二维码" />
        //             <Button
        //                 color="lightBlue"
        //                 buttonType="filled"
        //                 size="regular"
        //                 rounded={false}
        //                 block={false}
        //                 iconOnly={false}
        //                 ripple="light"
        //                 // onClick={handleNextStep}
        //             >
        //                 已扫描
        //             </Button>
        //         </>
        //     )}
        //     {step === 3 && (
        //         <>
        //             <h1 className="mb-4 text-2xl">输入动态令牌</h1>
        //             <Input
        //                 type="text"
        //                 color="lightBlue"
        //                 size="regular"
        //                 outline={true}
        //                 placeholder="动态令牌"
        //                 value={token}
        //                 // onChange={(e) => setToken(e.target.value)}
        //             />
        //             <Button
        //                 color="lightBlue"
        //                 buttonType="filled"
        //                 size="regular"
        //                 rounded={false}
        //                 block={false}
        //                 iconOnly={false}
        //                 ripple="light"
        //                 // onClick={handleNextStep}
        //             >
        //                 提交
        //             </Button>
        //         </>
        //     )}
        //     {step === 4 && (
        //         <h1 className="mb-4 text-2xl">二次验证已开启！</h1>
        //     )}
        // </div>
//     );
// }

// export default TwoFactorAuthPage;
