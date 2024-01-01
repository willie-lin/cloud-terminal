// 用户信息组件
import {useEffect, useState} from "react";
import {getUserByEmail} from "../../../api/api";
import {
    Avatar, Button,
    Card,
    CardBody,
    CardHeader,
    Typography
} from "@material-tailwind/react";
import {useNavigate} from "react-router-dom";

function UserInfo({ email }) {
    const [currentTime, setCurrentTime] = useState(new Date());
    const [userInfo, setUserInfo] = useState(null);

    // 获取用户信息
    useEffect(() => {
        if (email) {  // 添加这一行来检查email是否存在
            getUserByEmail(email)
                .then(data => setUserInfo(data))
                .catch(error => console.error('Error:', error));
        }
    }, [email]);

    // 更新当前时间
    useEffect(() => {
        const timer = setInterval(() => {
            setCurrentTime(new Date());
        }, 1000); // 每秒更新一次
        return () => {
            clearInterval(timer);
        };
    }, []);

    const navigate = useNavigate();

    const handleEdit = () => {
        navigate('/edit-user-info');  // 跳转到用户信息编辑页面
        // navigate(0);  // 强制刷新页面
    };

    return (
        <Card className="w-full max-w-2xl p-4 bg-white rounded shadow-md">
            <CardHeader floated={false} className="h-50">
                <div>
                    {userInfo && <Avatar src={userInfo.avatar} alt="avatar" variant="rounded" size="xxl"/>}
                </div>
            </CardHeader>
            <CardBody className="text-center">
                        <Typography variant="h4" color="blue-gray" className="mb-2">
                            欢迎，{userInfo?.nickname}!
                        </Typography>
                        {userInfo && <Typography color="blue-gray" className="font-medium" textGradient>
                            用户名: {userInfo.username}。
                        </Typography>}
                        {userInfo && <Typography color="blue-gray" className="font-medium" textGradient>
                            电子邮件: {userInfo.email}。
                        </Typography>}
                        {userInfo && <Typography color="blue-gray" className="font-medium" textGradient>
                            个人简介: {userInfo.bio}。
                        </Typography>}
                        <Typography color="blue-gray" className="font-medium" textGradient>
                            当前时间是 {currentTime.toLocaleTimeString()}。
                        </Typography>

                <Button onClick={handleEdit} color="lightBlue">修改</Button>
            </CardBody>
        </Card>
    );
}

export default UserInfo;