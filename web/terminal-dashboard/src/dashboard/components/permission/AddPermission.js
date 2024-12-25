import {Button, Card, CardBody, IconButton, Input, Option, Select, Tooltip, Typography} from "@material-tailwind/react";
import React, {useState} from "react";
import {TrashIcon} from "@heroicons/react/16/solid";
import {addPermission, checkPermissionName} from "../../../api/api";
import {useNavigate} from "react-router-dom";


function AddPermission({ onAddPermission, onClose }) {
    const [permissions, setPermissions] = useState([{ name: '', description: '', actions: [], resource_type: '' }]); // resourceType 初始为空
    const [nameError, setNameError] = useState('');

    const navigate = useNavigate();

    const availableActions = [
        { value: "create", label: "Create" },
        { value: "create,read", label: "Create, Read" },
        { value: "create,read,update", label: "Create, Read, Update" },
        { value: "create,read,update,delete", label: "Create, Read, Update, Delete" },
        { value: "create,read,update,delete,deploy", label: "Create, Read, Update, Delete, Deploy" },
    ];

    const handlePermissionChange = async (index, field, value) => {
        const newPermissions = [...permissions];
        newPermissions[index][field] = value;
        setPermissions(newPermissions);

        if (field === 'name') {
            try {
                const exists = await checkPermissionName({ name: value });
                setNameError(exists ? 'PermissionName already exists' : '');
            } catch (error) {
                console.error("检查权限名时出错:", error);
                setNameError('检查权限名时出错');
            }
        }
    };

    const handleActionChange = (index, selectedValue) => {
        const newPermissions = [...permissions];
        const finalActions = selectedValue.split(',');
        newPermissions[index].actions = finalActions;
        setPermissions(newPermissions);
    };




    const handleAddPermission = () => {
        setPermissions([...permissions, { name: '', description: '', actions: [], resource_type: '' }]);
    };

    const handleRemovePermission = (index) => {
        const newPermissions = [...permissions];
        newPermissions.splice(index, 1);
        setPermissions(newPermissions);
    };


    const handleSubmit = async (event) => {
        event.preventDefault();
        console.log(permissions);
        // 检查每个角色的'name'和'description'是否为空

        for (let permission of permissions) {
            if (!permission.name.trim() || !permission.description.trim() || !permission.resource_type.trim()) {
                alert('权限名、描述和资源类型不能为空');
                return;
            }
            // if (permission.actions.length === 0) { // 检查 actions 数组是否为空
            //     alert('每个权限至少需要选择一个操作');
            //     return;
            // }
        }
        try {
            // const formattedPermissions = permissions.map(p => ({
            //     Name: p.name,
            //     Description: p.description,
            //     Actions: p.actions,
            //     Module: p.resource_type,
            // }));
            await addPermission(permissions)
            onAddPermission(permissions);
            onClose();
            navigate("/")
        } catch (error) {
            console.error("添加权限时出错:", error);
            alert('添加权限失败');
        }
    };

    return (
        <Card className="w-256">
            <CardBody className="px-4 py-8">
                <div className="flex justify-between items-center mb-4">
                    <Typography variant="h4" color="gray">
                        Create Permission
                    </Typography>
                    <Button color="gray" buttonType="link"   onClick={onClose}>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5}
                             stroke="currentColor" className="w-4 h-4">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12"/>
                        </svg>
                    </Button>
                </div>
                <Typography variant="body2" color="blueGray" className="mb-4">
                    Enter the data for the new Permission.
                </Typography>
                <form onSubmit={handleSubmit}>
                    {permissions.map((permission, index) => (
                        <div className="mb-1 flex gap-6 items-end" key={index}>
                            <div className="flex flex-col">
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    {/*Name*/}
                                </Typography>
                                <Input
                                    variant="outlined"
                                    label="name"
                                    type="name"
                                    color="lightBlue"
                                    size="regular"
                                    outline={true}
                                    value={permission.name}
                                    name="name"
                                    onChange={(e) => handlePermissionChange(index, 'name', e.target.value)}
                                    className="mb-4"
                                    error={!!nameError}
                                />
                            </div>
                            <div className="flex flex-col">
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    {/*Description*/}
                                </Typography>
                                <Input
                                    variant="outlined"
                                    label="description"
                                    type="description"
                                    color="lightBlue"
                                    size="regular"
                                    outline={true}
                                    value={permission.description}
                                    name="description"
                                    onChange={(e) => handlePermissionChange(index, 'description', e.target.value)}
                                    className="mb-4"

                                />
                            </div>

                            <div className="flex flex-col">
                                <Typography variant="small" color="blue-gray" className="mb-2 font-medium">
                                    {/*Actions*/}
                                </Typography>
                                <Select
                                    size="md"
                                    label="Actions"
                                    multiple
                                    value={permission.actions.join(',')}
                                    onChange={(value) => handleActionChange(index, value)}
                                >
                                    {availableActions.map((action) => (
                                        <Option key={action.value} value={action.value}>
                                            {action.label}
                                        </Option>
                                    ))}
                                </Select>
                            </div>
                            <div>
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    {/*ResourceType*/}
                                </Typography>
                                <Input
                                    variant="outlined"
                                    label="ResourceType"
                                    type="text"
                                    color="lightBlue"
                                    size="regular"
                                    outline={true}
                                    value={permission.resource_type}
                                    name="description"
                                    onChange={(e) => handlePermissionChange(index, 'resource_type', e.target.value)}
                                    className="mb-2"
                                />
                            </div>

                            {permissions.length > 1 && ( // 只有当有多于一个权限时才显示删除按钮
                                <div className="self-end"> {/* 使用 self-end 将按钮放在右下角 */}
                                    <Tooltip content="删除权限" placement="top">
                                        <IconButton variant="text" color="red"
                                                    onClick={() => handleRemovePermission(index)}>
                                            <TrashIcon className="h-4 w-4"/>
                                        </IconButton>
                                    </Tooltip>
                                </div>
                            )}
                        </div>
                    ))}
                    <Button
                        color="lightBlue"
                        buttonType="filled"
                        size="regular"
                        rounded={false}
                        block={false}
                        iconOnly={false}
                        ripple="light"
                        onClick={handleAddPermission}
                    >
                        + Add additional role
                    </Button>
                    <Button fullWidth
                            type="submit"
                            color="lightBlue"
                            buttonType="filled"
                            size="regular"
                            rounded={false}
                            block={false}
                            iconOnly={false}
                            ripple="light"
                            className="mt-6"
                    >
                        Save
                    </Button>
                </form>
            </CardBody>
        </Card>
    )
}

export default AddPermission;