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
import {useCurrentTime, useFetchUserInfo} from "./UserHook";

function UserInfo({ email }) {
    const userInfo = useFetchUserInfo(email);
    const currentTime = useCurrentTime();
    const navigate = useNavigate();

    const handleEdit = () => {
        navigate('/edit-user-info',  { state: { userInfo } });  // 跳转到用户信息编辑页面
        // navigate(0);  // 强制刷新页面
    };
    return (
        <Card className="w-96">
            <CardHeader floated={false} className="h-50">
                <div>
                    {userInfo && <Avatar src={ userInfo.avatar } alt="avatar" variant="rounded" size="xxl"/>}
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