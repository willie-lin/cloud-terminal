// 用户信息组件
import { useState } from "react";
import {
    Avatar, Button,
    Card,
    CardBody,
    CardHeader,
    Typography
} from "@material-tailwind/react";
import {useCurrentTime, useFetchUserInfo} from "./UserHook";
import {
    MdEmail,
    MdOutlineAccessTimeFilled, MdOutlinePhoneIphone
} from "react-icons/md";
import {FaAudioDescription, FaUser} from "react-icons/fa";
import EditUserInfo from "./EditUserInfo";

function UserInfo({ user }) {

    const currentTime = useCurrentTime();

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

    // return (
//         <Card className="w-1/2">
//             <CardHeader floated={false} className="h-50">
//                 <div>
//                     {userInfo && <Avatar src={ userInfo.avatar } alt="avatar" variant="rounded" size="xxl"/>}
//                 </div>
//             </CardHeader>
//             {/*<CardHeader floated={false} className="h-50 flex items-center justify-center">*/}
//             {/*    <div>*/}
//             {/*        {userInfo && <Avatar src={userInfo.avatar} alt="avatar" variant="rounded" size="xxl"/>}*/}
//             {/*    </div>*/}
//             {/*</CardHeader>*/}
//             <CardBody className="text-left">
//                 <Typography variant="h4" color="blue-gray" className="mb-2">
//                     欢迎，{userInfo?.nickname}!
//                 </Typography>
//                 <div className="flex items-center">
//                     <FaUser/>
//                     {userInfo && (
//                         <Typography color="blue-gray" className="font-medium ml-2" textGradient>
//                             {userInfo.username}
//                         </Typography>
//                     )}
//                 </div>
//                 <div className="flex items-center">
//                     <MdEmail/>
//                     {userInfo && (
//                         <Typography color="blue-gray" className="font-medium ml-2" textGradient>
//                             {userInfo.email}
//                         </Typography>
//                     )}
//                 </div>
//
//                 <div className="flex items-center">
//                     <MdOutlinePhoneIphone size="1.3em"/>
//                     {userInfo && (
//                         <Typography color="blue-gray" className="font-medium ml-2" textGradient>
//                             {userInfo.phone_number}
//                         </Typography>
//                     )}
//                 </div>
//
//                 <div className="flex items-center">
//                     <FaAudioDescription/>
//                     {userInfo && (
//                         <Typography color="blue-gray" className="font-medium ml-2" textGradient>
//                             {userInfo.bio}
//                         </Typography>
//                     )}
//                 </div>
//
//
//                 <div className="flex items-center">
//                     <MdOutlineAccessTimeFilled size="1.2em"/>
//                     <Typography color="blue-gray" className="font-medium ml-2" textGradient>
//                         {currentTime.toLocaleTimeString()}
//                     </Typography>
//                 </div>
//
//
//             </CardBody>
//             <Button
//                 // onClick={handleEdit}
//                 color="lightBlue"
//                 onClick={() => openEditUser(userInfo)}
//             >
//                 Edit User
//             </Button>
//             {isEditUserOpen && (
//                 <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
//                      onClick={(e) => {
//                          // 如果事件的目标是这个容器本身，那么关闭模态窗口
//                          if (e.target === e.currentTarget) {
//                              closeEditUser();
//                          }
//                      }}
//                 >
//                     <EditUserInfo user={editUser} onEditUser={handleEditUser} onClose={closeEditUser}/>
//                 </div>
//             )}
//         </Card>
//     );
// }

    return (
        <Card className="w-full max-w-lg p-6 mx-auto mt-12 bg-white shadow-lg rounded-lg relative">
            <div className="flex justify-center mb-6">
                {userInfo && <Avatar src={userInfo.avatar} alt="avatar" variant="rounded" size="xxl" />}
            </div>
            <CardBody className="text-left">
                <Typography variant="h4" color="blue-gray" className="mb-4 text-center">
                    欢迎，{userInfo?.nickname}!
                </Typography>
                <div className="flex items-center mb-4">
                    <FaUser className="mr-2 text-lg" />
                    {userInfo && (
                        <Typography color="blue-gray" className="font-medium" textGradient>
                            {userInfo.username}
                        </Typography>
                    )}
                </div>
                <div className="flex items-center mb-4">
                    <MdEmail className="mr-2 text-lg" />
                    {userInfo && (
                        <Typography color="blue-gray" className="font-medium" textGradient>
                            {userInfo.email}
                        </Typography>
                    )}
                </div>
                <div className="flex items-center mb-4">
                    <MdOutlinePhoneIphone className="mr-2 text-lg" />
                    {userInfo && (
                        <Typography color="blue-gray" className="font-medium" textGradient>
                            {userInfo.phone_number}
                        </Typography>
                    )}
                </div>
                <div className="flex items-center mb-4">
                    <FaAudioDescription className="mr-2 text-lg" />
                    {userInfo && (
                        <Typography color="blue-gray" className="font-medium" textGradient>
                            {userInfo.bio}
                        </Typography>
                    )}
                </div>
                <div className="flex items-center mb-4">
                    <MdOutlineAccessTimeFilled className="mr-2 text-lg" />
                    <Typography color="blue-gray" className="font-medium" textGradient>
                        {currentTime.toLocaleTimeString()}
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
                    <EditUserInfo user={editUser} onEditUser={handleEditUser} onClose={closeEditUser} />
                </div>
            )}
        </Card>
    );
}
export default UserInfo;