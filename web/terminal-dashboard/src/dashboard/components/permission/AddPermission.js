
import {Button, Card, CardBody, IconButton, Input, Tooltip, Typography} from "@material-tailwind/react";
import React, {useState} from "react";
import {TrashIcon} from "@heroicons/react/16/solid";
import {addPermission} from "../../../api/api";


function AddPermission({ onAddPermission, onClose }) {

    const [permissions, setPermissions] = useState([{ name: '', description: '' }]);

    const handlePermissionChange = (index, field, value) => {
        const newPermissions = [...permissions];
        newPermissions[index][field] = value;
        setPermissions(newPermissions);
    };


    const handleAddPermission = () => {
        setPermissions([...permissions, { name: '', description: '' }]);
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
            if (!permission.name.trim() || !permission.description.trim()) {
                console.error('Error: Permission name and description cannot be empty');
                return;  // 如果'name'或'description'为空，阻止提交
            }
        }
        try {
            await addPermission(permissions)
            onAddPermission(permissions);
        } catch (error) {
            console.log(error)
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
                    Enter the data for the new permission.
                </Typography>
                <form onSubmit={handleSubmit}>
                    {permissions.map((permission, index) => (
                        <div className="mb-1 flex gap-6 items-end" key={index}>
                            <div className="flex flex-col">
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    {/*Key*/}
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
                                />
                            </div>
                            <div className="flex flex-col">
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    {/*Display Name*/}
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
                            <Tooltip content="Delete Role" placement="top">
                                <IconButton variant="text"
                                            color="red"
                                            buttonType="filled"
                                            size="regular"
                                            rounded={false}
                                            block={false}
                                            iconOnly={true}
                                            ripple="light"
                                    // onClick={() => openDeleteUser(user)}
                                            onClick={() => handleRemovePermission(index)}
                                >
                                    <TrashIcon className="h-4 w-4"/>
                                </IconButton>
                            </Tooltip>
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