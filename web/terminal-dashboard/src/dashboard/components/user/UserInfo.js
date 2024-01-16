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
import {
    MdAccessTimeFilled,
    MdDescription,
    MdEmail,
    MdOutlineAccessTime,
    MdOutlineAccessTimeFilled
} from "react-icons/md";
import {FaAudioDescription, FaUser} from "react-icons/fa6";

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
                <div className="flex items-center">
                    <FaUser/>
                    {userInfo && (
                        <Typography color="blue-gray" className="font-medium ml-2" textGradient>
                            {userInfo.username}
                        </Typography>
                    )}
                </div>
                <div className="flex items-center">
                    <MdEmail/>
                    {userInfo && (
                        <Typography color="blue-gray" className="font-medium ml-2" textGradient>
                            {userInfo.email}
                        </Typography>
                    )}
                </div>
                <div className="flex items-start">
                    <FaAudioDescription size="1.5em" />
                    {userInfo && (
                        <Typography color="blue-gray" className="font-medium ml-2" textGradient>
                            {userInfo.bio}
                        </Typography>
                    )}
                </div>

                <div className="flex items-center">
                    <MdOutlineAccessTimeFilled size="1.2em" />
                    <Typography color="blue-gray" className="font-medium ml-2" textGradient>
                        {currentTime.toLocaleTimeString()}
                    </Typography>
                </div>


            </CardBody>
            <Button onClick={handleEdit} color="lightBlue">Edit User</Button>
        </Card>

        // <Card className="w-96 mx-auto my-10">
        //     <CardHeader floated={false} className="h-50">
        //         <div className="flex justify-center">
        //             {userInfo && <Avatar src={ userInfo.avatar } alt="avatar" variant="rounded" size="xxl"/>}
        //         </div>
        //     </CardHeader>
        //     <CardBody className="text-center">
        //         <Typography variant="h4" color="blue-gray" className="mb-2">
        //             欢迎，{userInfo?.nickname}!
        //         </Typography>
        //         {userInfo && <Typography color="blue-gray" className="font-medium mb-2">
        //             <HeartIcon name="account_circle" size="xl" /> USERNAME: {userInfo.username}
        //         </Typography>}
        //         {userInfo && <Typography color="blue-gray" className="font-medium mb-2">
        //             <HeartIcon name="email" size="xl" /> EMAIL: {userInfo.email}
        //         </Typography>}
        //         {userInfo && <Typography color="blue-gray" className="font-medium mb-2">
        //             <HeartIcon name="info" size="xl" /> BIO: {userInfo.bio}
        //         </Typography>}
        //         <Typography color="blue-gray" className="font-medium mb-2">
        //             <HeartIcon name="access_time" size="xl" /> TIME: {currentTime.toLocaleTimeString()}
        //         </Typography>
        //     </CardBody>
        //     <CardFooter>
        //         <Button onClick={handleEdit} color="lightBlue" fullWidth>Edit User</Button>
        //     </CardFooter>
        // </Card>

    );
}

export default UserInfo;