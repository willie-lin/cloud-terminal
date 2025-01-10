import {Chip, IconButton, Tooltip, Typography} from "@material-tailwind/react";
import {PencilIcon, TrashIcon} from "@heroicons/react/16/solid";
import React, {useContext, useState} from "react";
import DeleteUserForm from "../user/DeleteUserForm";
import DeleteRole from "./DeleteRole";
import {AuthContext} from "../../../App";

function RenderRole({ role, isLast }) {

    const { currentUser } = useContext(AuthContext);
    // 判断当前用户是否具有删除权限
    // const canDelete = currentUser?.roleName === 'Admin' || currentUser?.roleName === 'SuperAdmin'
    const canDelete = (currentUser?.isTenantAdmin  || currentUser?.roleName === 'super_admin') && !role.is_default

    const classes = isLast ? "p-4" : "p-4 border-b border-blue-gray-50";

    const [deleteRole, setDeleteRole] = useState(null);
    const [isDeleteRoleOpen, setIsDeleteRoleOpen] = useState(false);

    const handleDeleteRole = () => {
        setIsDeleteRoleOpen(false);
    };

    function openDeleteRole(user) {
        setDeleteRole(user)
        setIsDeleteRoleOpen(true)
    }

    function closeDeleteRole() {
        setIsDeleteRoleOpen(false);
    }

    const getColor = roleName => {
        if (roleName.toLowerCase().includes("admin")) return "green";
        return "gray";
    };


    return (
        <tr key={role.id}>
            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                    {role.id}
                </Typography>
            </td>
            <td className={classes}>
                <div className="w-max">
                    <Chip
                        size="sm"
                        variant="ghost"
                        value={role.name}
                        color={getColor(role.name)
                        }
                    />
                </div>
            </td>
            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                    {role.description}
                </Typography>
            </td>
            <td className={classes}>
                <div className="w-max">
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        <Chip variant="ghost" size="sm" value={role.is_disabled ? "TRUE" : "FALSE"}
                              color={role.is_disabled ? "green" : "blue-gray"}
                        />
                    </Typography>
                </div>
            </td>
            <td className={classes}>
                <div className="w-max">
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        <Chip variant="ghost" size="sm" value={role.is_default ? "TRUE" : "FALSE"}
                              color={role.is_default ? "green" : "blue-gray"}
                        />
                    </Typography>
                </div>
            </td>

            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                    {role.created_at}
                </Typography>
            </td>
            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                    {role.updated_at}
                </Typography>
            </td>
            {canDelete && (
                <td className={classes}>
                    <Tooltip content="Delete Role" placement="top">
                        <IconButton variant="text"
                                    color="red"
                                    buttonType="filled"
                                    size="regular"
                                    rounded={false}
                                    block={false}
                                    iconOnly={true}
                                    ripple="light"
                                    onClick={() => openDeleteRole(role)}
                            //         onClick={() => handleRemoveRole(index)}
                        >
                            <TrashIcon className="h-4 w-4"/>
                        </IconButton>
                    </Tooltip>
                    {isDeleteRoleOpen && (
                        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
                             onClick={(e) => {
                                 // 如果事件的目标是这个容器本身，那么关闭模态窗口
                                 if (e.target === e.currentTarget) {
                                     closeDeleteRole();
                                 }
                             }}
                        >
                            <DeleteRole role={deleteRole} onDeleteRole={handleDeleteRole} onClose={closeDeleteRole}/>
                        </div>
                    )}
                </td>
            )}
        </tr>
    );
}

export default RenderRole;