import {Chip, IconButton, Tooltip, Typography} from "@material-tailwind/react";
import {PencilIcon, TrashIcon} from "@heroicons/react/16/solid";
import React, {useState} from "react";
import DeleteUserForm from "../user/DeleteUserForm";
import DeleteRole from "./DeleteRole";

function RenderRole({ role, isLast }) {
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
                        color={
                            role.name === "管理员"
                                ? "green"
                                : role.name === "普通用户"
                                    ? "blue"
                                    : "gray"
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
                <Typography variant="small" color="blue-gray" className="font-normal">
                {role.created_at}
                </Typography>
            </td>
            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                    {role.updated_at}
                </Typography>
            </td>
            {/*<td className={classes}>*/}
            {/*    <Tooltip content="Edit Role">*/}
            {/*        <IconButton variant="text">*/}
            {/*            <PencilIcon className="h-4 w-4"/>*/}
            {/*        </IconButton>*/}
            {/*    </Tooltip>*/}
            {/*</td>*/}
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
                             }}}
                    >
                        <DeleteRole role={deleteRole} onDeleteRole={handleDeleteRole} onClose={closeDeleteRole}/>
                    </div>
                )}
            </td>
        </tr>
    );
}

export default RenderRole;