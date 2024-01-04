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
        // <Card className="w-full max-w-2xl p-4 bg-white rounded shadow-md">

        <Card className="w-96">
            <CardHeader floated={false} className="h-50">
                <div>
                    {userInfo && <Avatar src={userInfo.avatar} alt="avatar" variant="rounded" size="xxl"/>}
                    {/*<img src="https://docs.material-tailwind.com/img/team-3.jpg" alt="profile-picture"/>*/}
                </div>
            </CardHeader>
            {/*<CardBody className="text-center">*/}
            <CardBody className="text-left">
            <Typography variant="h4" color="blue-gray" className="mb-2">
                            欢迎，{userInfo?.nickname}!
            </Typography>

                {userInfo && <Typography color="blue-gray" className="font-medium" textGradient>
                            USERNAME: {userInfo.username}。
                </Typography>}
                {userInfo && <Typography color="blue-gray" className="font-medium" textGradient>
                            EMAIL: {userInfo.email}。
                </Typography>}
                {userInfo && <Typography color="blue-gray" className="font-medium" textGradient>
                            BIO: {userInfo.bio}。
                </Typography>}
                <Typography color="blue-gray" className="font-medium" textGradient>
                            TiME: {currentTime.toLocaleTimeString()}。
                </Typography>
            </CardBody>
            <Button onClick={handleEdit} color="lightBlue">修改</Button>
        </Card>
    );
}

export default UserInfo;