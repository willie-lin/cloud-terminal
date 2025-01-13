// 用户信息组件
import React, { useState } from "react";
import {
    Avatar, Button,
    Card,
    CardBody, CardFooter,
    CardHeader,
    Typography
} from "@material-tailwind/react";
import {useCurrentDateTime, useCurrentTime, useFetchUserInfo} from "./UserHook";
import {
    MdEmail,
    MdOutlineAccessTimeFilled, MdOutlinePhoneIphone
} from "react-icons/md";
import {FaAudioDescription, FaUser} from "react-icons/fa";
import EditUserInfo from "./EditUserInfo";

function UserInfo({ user }) {

    const currentTime = useCurrentTime();
    const currentDateTime = useCurrentDateTime();

    // 在 UserInfo 组件中
    const userInfo = useFetchUserInfo(user.email);

    const [editUser, setEditUser] = useState(null);
    const [isEditUserOpen, setIsEditUserOpen] = useState(false);
    const handleEditUser = () => {
        setIsEditUserOpen(false);
    };
    const openEditUser = (userInfo) => {
        setEditUser(userInfo)
        setIsEditUserOpen(true);
    };
    const closeEditUser = () => {
        setIsEditUserOpen(false);
    };


    return (
        <div className="w-full ">
            <Card className="w-full shadow-lg rounded-lg">
            {/*<Card className="w-full max-w-lg p-6 mx-auto mt-12 bg-white shadow-lg rounded-lg relative">*/}
                <div className="flex justify-center mb-6">
                    {userInfo && <Avatar src={userInfo.avatar} alt="avatar" withBorder={true} color={"blue"}
                                         className="p-0.5"
                                         // className="border border-green-500 shadow-xl shadow-green-900/20 ring-4 ring-green-500/30"
                                         variant="rounded" size="xxl"/>}
                </div>
                <CardBody className="text-left">
                    <Typography variant="h4" color="blue-gray" className="mb-4 text-center">
                        欢迎，{userInfo?.nickname}!
                    </Typography>
                    <div className="flex items-center mb-4">
                        <FaUser className="mr-2 text-lg"/>
                        {userInfo && (
                            <Typography color="blue-gray" className="font-medium" textGradient>
                                {userInfo.username}
                            </Typography>
                        )}
                    </div>
                    <div className="flex items-center mb-4">
                        <MdEmail className="mr-2 text-lg"/>
                        {userInfo && (
                            <Typography color="blue-gray" className="font-medium" textGradient>
                                {userInfo.email}
                            </Typography>
                        )}
                    </div>
                    <div className="flex items-center mb-4">
                        <MdOutlinePhoneIphone className="mr-2 text-lg"/>
                        {userInfo && (
                            <Typography color="blue-gray" className="font-medium" textGradient>
                                {userInfo.phone_number}
                            </Typography>
                        )}
                    </div>
                    <div className="flex items-center mb-4">
                        <FaAudioDescription className="mr-2 text-lg"/>
                        {userInfo && (
                            <Typography color="blue-gray" className="font-medium" textGradient>
                                {userInfo.bio}
                            </Typography>
                        )}
                    </div>
                    <div className="flex items-center mb-4">
                        <MdOutlineAccessTimeFilled className="mr-2 text-lg"/>
                        <Typography color="blue-gray" className="font-medium" textGradient>
                            {/*{currentTime.toLocaleTimeString()}*/}
                            { currentDateTime.toLocaleUpperCase()}

                        </Typography>
                    </div>
                </CardBody>
                <div className="text-center mt-4">
                    <Button
                        color="lightBlue"
                        onClick={() => openEditUser(userInfo)}
                        className="w-full"
                    >
                        Edit User
                    </Button>
                </div>
                {isEditUserOpen && (
                    <div
                        className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
                        onClick={(e) => {
                            if (e.target === e.currentTarget) {
                                closeEditUser();
                            }
                        }}
                    >
                        <EditUserInfo user={editUser} onEditUser={handleEditUser} onClose={closeEditUser}/>
                    </div>
                )}
                <CardFooter className="text-center pt-0">
                    <Typography variant="small" className="mt-6 flex justify-center">
                        {/*Need Help?*/}
                        <Typography
                            as="a"
                            href="#"
                            variant="small"
                            color="lightBlue"
                            className="ml-1 font-bold"
                        >
                            {/*Contact Support.*/}
                        </Typography>
                    </Typography>
                </CardFooter>
            </Card>
        </div>

    );
}


            export default UserInfo;