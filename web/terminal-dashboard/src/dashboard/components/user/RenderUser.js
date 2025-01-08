import {Avatar, Chip, IconButton, Tooltip, Typography} from "@material-tailwind/react";
import {PencilIcon, TrashIcon} from "@heroicons/react/16/solid";
import UpdateUser from "./UpdateUser";
import {useContext, useState} from "react";
import DeleteUserForm from "./DeleteUserForm";
import {AuthContext} from "../../../App";


function RenderUser({ user, isLast, isDarkMode }) {
    const { currentUser } = useContext(AuthContext);
    // 判断当前用户是否具有删除权限
    const canDelete = currentUser?.roleName === 'admin' || currentUser?.roleName === 'super_admin'


    // const classes = isLast ? "p-2" : "p-2 border-b border-blue-gray-50";
    const classes = `p-2 ${isLast ? '' : `border-b ${isDarkMode ? 'border-gray-700' : 'border-blue-gray-50'}`}`;


    const [editingUser, setEditingUser] = useState(null);
    const [isUpdateUserOpen, setIsUpdateUserOpen] = useState(false);

    const handleUpdateUser = () => {
        setIsUpdateUserOpen(false);
    };
    const openUpdateUser = (user) => {
        setEditingUser(user)
        setIsUpdateUserOpen(true);
    };
    const closeUpdateUser = () => {
        setIsUpdateUserOpen(false);
    };

    const [deleteUser, setDeleteUser] = useState(null);
    const [isDeleteUserOpen, setIsDeleteUserOpen] = useState(false);

    const handleDeleteUser = () => {
        setIsDeleteUserOpen(false);
    };

    function openDeleteUser(user) {
        setDeleteUser(user)
        setIsDeleteUserOpen(true)
    }

    function closeDeleteUser() {
        setIsDeleteUserOpen(false);
    }


    return (
        <>
            <tr key={user.id} className={isDarkMode ? 'text-gray-100' : 'text-gray-900'}>
                <td className={classes}>
                    <div className="flex flex-col">
                        {/*<Typography variant="small" color="blue-gray" className="font-normal">*/}
                        <Typography variant="small"
                                    className={`font-normal ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-600'}`}>
                            {user.id}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex items-center gap-3">
                        {/* 这里可以添加Avatar组件 */}
                        <Avatar src={user.avatar} alt={user.nickname} size="sm"/>
                        <div className="flex flex-col">
                            {/*<Typography variant="small" color="blue-gray" className="font-normal">*/}
                            <Typography variant="small"
                                        className={`font-normal ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-600'}`}>
                                {user.nickname}
                            </Typography>
                            {/*<Typography variant="small" color="blue-gray" className="font-normal">*/}
                            {/*<Typography variant="small" className={`font-normal ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-600'}`}>*/}
                            {/*    {user.email}*/}
                            {/*</Typography>*/}
                        </div>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        {/*<Typography variant="small" color="blue-gray" className="font-normal">*/}
                        <Typography variant="small"
                                    className={`font-normal ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-600'}`}>
                            {user.username}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        {/*<Typography variant="small" color="blue-gray" className="font-normal">*/}
                        <Typography variant="small"
                                    className={`font-normal ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-600'}`}>
                            {user.email}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        {/*<Typography variant="small" color="blue-gray" className="font-normal">*/}
                        <Typography variant="small"
                                    className={`font-normal ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-600'}`}>
                            {user.phone_number}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        {/*<Typography variant="small" color="blue-gray" className="font-normal">*/}
                        <Typography variant="small"
                                    className={`font-normal ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-600'}`}>
                            {user.bio}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        {/*<Typography variant="small" color="blue-gray" className="font-normal">*/}
                        <Typography variant="small"
                                    className={`font-normal ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-600'}`}>
                            {user.totp_secret}
                        </Typography>
                    </div>
                </td>
                {/*<td className={classes}>*/}
                {/*    <div className="flex flex-col">*/}
                {/*        /!*<Typography variant="small" color="blue-gray" className="font-normal">*!/*/}
                {/*        <Typography variant="small"*/}
                {/*                    className={`font-normal ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-600'}`}>*/}
                {/*            {(() => {*/}
                {/*                switch (user.status) {*/}
                {/*                    case 'active':*/}
                {/*                        return '启用';*/}
                {/*                    case 'inactive':*/}
                {/*                        return '禁用';*/}
                {/*                    case 'blocked':*/}
                {/*                        return '已阻止';*/}
                {/*                    default:*/}
                {/*                        return '未知状态'; // 处理未知状态的情况*/}
                {/*                }*/}
                {/*            })()}*/}
                {/*        </Typography>*/}
                {/*    </div>*/}
                {/*</td>*/}
                <td className={classes}>
                    <div className="w-max">
                        <Chip variant="ghost" size="sm" value={
                            (() => {
                            switch (user.status) {
                                case 'active':
                                    return '启用';
                                case 'inactive':
                                    return '禁用';
                                case 'blocked':
                                    return '已阻止';
                                default:
                                    return '未知状态'; // 处理未知状态的情况
                            }
                        })()}
                              color={user.status ? "green" : "blue-gray"}
                        />
                    </div>
                </td>
                <td className={classes}>
                    <div className="w-max">
                        <Chip variant="ghost" size="sm" value={user.online ? "在线" : "离线"}
                              color={user.online ? "green" : "blue-gray"}
                        />
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        {/*<Typography variant="small" color="blue-gray" className="font-normal">*/}
                        <Typography variant="small"
                                    className={`font-normal ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-600'}`}>
                            {user.created_at}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        {/*<Typography variant="small" color="blue-gray" className="font-normal">*/}
                        <Typography variant="small"
                                    className={`font-normal ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-600'}`}>
                            {user.updated_at}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        {/*<Typography variant="small" color="blue-gray" className="font-normal">*/}
                        <Typography variant="small"
                                    className={`font-normal ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-600'}`}>
                            {user.last_login_time}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <Tooltip content="Edit User" placement="top">
                        <IconButton variant="text"
                                    color="lightBlue"
                                    buttonType="filled"
                                    size="regular"
                                    rounded={false}
                                    block={false}
                                    iconOnly={true}
                                    ripple="light"  // dark
                                    onClick={() => openUpdateUser(user)}
                        >
                            <PencilIcon className="h-4 w-4"/>
                        </IconButton>
                    </Tooltip>
                    {isUpdateUserOpen && (
                        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
                             onClick={(e) => {
                                 // 如果事件的目标是这个容器本身，那么关闭模态窗口
                                 if (e.target === e.currentTarget) {
                                     closeUpdateUser();
                                 }
                             }}
                        >
                            <UpdateUser user={editingUser} onUpdateUser={handleUpdateUser} onClose={closeUpdateUser}/>
                        </div>
                    )}
                </td>
                {canDelete && (
                    <td className={classes}>
                        <Tooltip content="Delete User" placement="top">
                            <IconButton
                                variant="text"
                                color="red"
                                buttonType="filled"
                                size="regular"
                                rounded={false}
                                block={false}
                                iconOnly={true}
                                ripple="light"
                                onClick={() => openDeleteUser(user)}
                            >
                                <TrashIcon className="h-4 w-4"/>
                            </IconButton>
                        </Tooltip>
                        {isDeleteUserOpen && (
                            <div
                                className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
                                onClick={(e) => {
                                    if (e.target === e.currentTarget) {
                                        closeDeleteUser();
                                    }
                                }}
                            >
                                <DeleteUserForm user={deleteUser} onDeleteUser={handleDeleteUser}
                                                onClose={closeDeleteUser}/>
                            </div>
                        )}
                    </td>
                )}
            </tr>
        </>
    );
}

export default RenderUser;

// 然后在你的代码中使用这个组件：
