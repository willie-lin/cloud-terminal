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
import {FaAudioDescription, FaUser} from "react-icons/fa";
import {BiObjectsHorizontalRight} from "react-icons/bi";
import UpdateUser from "./UpdateUser";

function UserInfo({ email }) {
    const userInfo = useFetchUserInfo(email);
    const currentTime = useCurrentTime();
    const navigate = useNavigate();

    // const handleEdit = () => {
    //     navigate('/edit-user-info',  { state: { userInfo } });  // 跳转到用户信息编辑页面
    //     // navigate(0);  // 强制刷新页面
    // };


    const [editUser, setEditUser] = useState(null);
    const [isUpdateUserOpen, setIsUpdateUserOpen] = useState(false);
    const handleUpdateUser = () => {
        setIsUpdateUserOpen(false);
    };
    const openUpdateUser = (user) => {
        setEditUser(user)
        setIsUpdateUserOpen(true);
    };
    const closeUpdateUser = () => {
        setIsUpdateUserOpen(false);
    };


    return (
        <Card className="w-1/2">
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

                <div className="flex items-center">
                    <FaAudioDescription />
                    {userInfo && (
                        <Typography color="blue-gray" className="font-medium ml-2" textGradient>
                            {userInfo.bio}
                        </Typography>
                    )}
                </div>


                <div className="flex items-center">
                    <MdOutlineAccessTimeFilled size="1.2em"/>
                    <Typography color="blue-gray" className="font-medium ml-2" textGradient>
                        {currentTime.toLocaleTimeString()}
                    </Typography>
                </div>


            </CardBody>
            <Button
                // onClick={handleEdit}
                    color="lightBlue"
                    onClick={() => openUpdateUser(userInfo)}
            >
                Edit User
            </Button>
            {isUpdateUserOpen && (
                <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
                     onClick={(e) => {
                         // 如果事件的目标是这个容器本身，那么关闭模态窗口
                         if (e.target === e.currentTarget) {
                             closeUpdateUser();
                         }
                     }}
                >
                    <UpdateUser user={editUser} onUpdateUser={handleUpdateUser}  onClose={closeUpdateUser}/>
                </div>
            )}
        </Card>
    );
}

export default UserInfo;