import {Avatar, Chip, IconButton, Tooltip, Typography} from "@material-tailwind/react";
import {PencilIcon, TrashIcon} from "@heroicons/react/16/solid";
import UpdateUser from "./UpdateUser";
import {useState} from "react";
import DeleteUserForm from "./DeleteUserForm";


function UserRow({ user, isLast }) {
    const classes = isLast ? "p-2" : "p-2 border-b border-blue-gray-50";

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
            <tr key={user.id}>
                <td className={classes}>
                    <div className="flex flex-col">
                        <Typography variant="small" color="blue-gray" className="font-normal">
                            {user.id}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex items-center gap-3">
                        {/* 这里可以添加Avatar组件 */}
                        <Avatar src={user.avatar} alt={user.nickname} size="sm"/>
                        <div className="flex flex-col">
                            <Typography variant="small" color="blue-gray"
                                        className="font-normal">
                                {user.nickname}
                            </Typography>
                            <Typography variant="small" color="blue-gray"
                                        className="font-normal">
                                {user.email}
                            </Typography>
                        </div>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        <Typography variant="small" color="blue-gray" className="font-normal">
                            {user.username}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        <Typography variant="small" color="blue-gray" className="font-normal">
                            {user.email}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        <Typography variant="small" color="blue-gray" className="font-normal">
                            {user.phone}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        <Typography variant="small" color="blue-gray" className="font-normal">
                            {user.bio}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        <Typography variant="small" color="blue-gray" className="font-normal">
                            {user.totp_secret}
                        </Typography>
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
                        <Typography variant="small" color="blue-gray" className="font-normal">
                            {user.enable_type ? '启用' : '禁用'}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        <Typography variant="small" color="blue-gray" className="font-normal">
                            {user.created_at}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        <Typography variant="small" color="blue-gray" className="font-normal">
                            {user.updated_at}
                        </Typography>
                    </div>
                </td>
                <td className={classes}>
                    <div className="flex flex-col">
                        <Typography variant="small" color="blue-gray" className="font-normal">
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
                            <UpdateUser user={editingUser} onUpdateUser={handleUpdateUser}  onClose={closeUpdateUser}/>
                        </div>
                    )}
                </td>
                <td className={classes}>
                    <Tooltip content="Delete User" placement="top">
                        <IconButton variant="text"
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
                        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
                             onClick={(e) => {
                                 // 如果事件的目标是这个容器本身，那么关闭模态窗口
                                 if (e.target === e.currentTarget) {
                                     closeDeleteUser();
                                 }}}
                        >
                        <DeleteUserForm user={deleteUser} onDeleteUser={handleDeleteUser} onclose={closeDeleteUser}/>
                        </div>
                    )}
                </td>
            </tr>

    </>
    );
}
export default UserRow;

// 然后在你的代码中使用这个组件：
