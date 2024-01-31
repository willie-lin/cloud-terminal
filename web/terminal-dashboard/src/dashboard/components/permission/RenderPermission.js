import {Chip, IconButton, Tooltip, Typography} from "@material-tailwind/react";
import {TrashIcon} from "@heroicons/react/16/solid";
import React, {useState} from "react";
import DeletePermission from "./DeletePermission";

function RenderPermission({ permission, isLast }) {

    const classes = isLast ? "p-4" : "p-4 border-b border-blue-gray-50";

    const [deletePermission, setDeletePermission] = useState(null);
    const [isDeletePermissionOpen, setIsDeletePermissionOpen] = useState(false);

    const handleDeletePermission = () => {
        setIsDeletePermissionOpen(false);
    };

    function openDeletePermission(user) {
        setDeletePermission(user)
        setIsDeletePermissionOpen(true)
    }

    function closeDeletePermission() {
        setIsDeletePermissionOpen(false);
    }
    
    

    return (
        <tr key={permission.id}>
            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                    {permission.id}
                </Typography>
            </td>
            <td className={classes}>
                <div className="w-max">
                    <Chip
                        size="sm"
                        variant="ghost"
                        value={permission.name}
                        color={
                            permission.name === "管理员"
                                ? "green"
                                : permission.name === "普通用户"
                                    ? "blue"
                                    : "gray"
                        }
                    />
                </div>
            </td>
            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                    {permission.description}
                </Typography>
            </td>
            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                    {permission.created_at}
                </Typography>
            </td>
            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                    {permission.updated_at}
                </Typography>
            </td>
            {/*<td className={classes}>*/}
            {/*    <Tooltip content="Edit permission">*/}
            {/*        <IconButton variant="text">*/}
            {/*            <PencilIcon className="h-4 w-4"/>*/}
            {/*        </IconButton>*/}
            {/*    </Tooltip>*/}
            {/*</td>*/}
            <td className={classes}>
                <Tooltip content="Delete permission" placement="top">
                    <IconButton variant="text"
                                color="red"
                                buttonType="filled"
                                size="regular"
                                rounded={false}
                                block={false}
                                iconOnly={true}
                                ripple="light"
                                onClick={() => openDeletePermission(permission)}
                        //         onClick={() => handleRemovepermission(index)}
                    >
                        <TrashIcon className="h-4 w-4"/>
                    </IconButton>
                </Tooltip>
                {isDeletePermissionOpen && (
                    <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
                         onClick={(e) => {
                             // 如果事件的目标是这个容器本身，那么关闭模态窗口
                             if (e.target === e.currentTarget) {
                                 closeDeletePermission();
                             }}}
                    >
                        <DeletePermission permission={deletePermission} onDeletepermission={handleDeletePermission} onClose={closeDeletePermission}/>
                    </div>
                )}
            </td>
        </tr>
    );
}
export default RenderPermission;