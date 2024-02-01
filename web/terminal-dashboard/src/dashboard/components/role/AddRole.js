import {Button, Card, CardBody, IconButton, Input, Tooltip, Typography} from "@material-tailwind/react";
import React, {useState} from "react";
import {TrashIcon} from "@heroicons/react/16/solid";
import {addRole, checkEmail, checkRoleName} from "../../../api/api";


function AddRole({ onAddrole, onClose }) {

    const [nameError, setNameError] = useState('');

    const [roles, setRoles] = useState([{ name: '', description: '' }]);

    const handleRoleChange = async (index, field, value) => {
        const newRoles = [...roles];
        newRoles[index][field] = value;
        setRoles(newRoles);
        if (field === 'name') {
            try {
                const date = { name: value };  // 创建一个包含名称的对象
                const exists = await checkRoleName(date); // 等待函数完成
                setNameError(exists ? 'Name already registered' : '');
            } catch (error) {
                console.error(error);
            }
        }
    };
    const handleAddRole = () => {
        setRoles([...roles, { name: '', description: '' }]);
    };

    const handleRemoveRole = (index) => {
        const newRoles = [...roles];
        newRoles.splice(index, 1);
        setRoles(newRoles);
    };


    const handleSubmit = async (event) => {
        event.preventDefault();
        // console.log(roles);
        // 检查每个角色的'name'和'description'是否为空
        for (let role of roles) {
            if (!role.name.trim() || !role.description.trim()) {
                console.error('Error: Role name and description cannot be empty');
                return;  // 如果'name'或'description'为空，阻止提交
            }
        }
        try {
            await addRole(roles)
            onAddrole(roles);
        } catch (error) {
            console.log(error)
        }
    };

    return (
        <Card className="w-256">
            <CardBody className="px-4 py-8">
                <div className="flex justify-between items-center mb-4">
                    <Typography variant="h4" color="gray">
                        Create Role
                    </Typography>
                    <Button color="gray" buttonType="link"   onClick={onClose}>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5}
                             stroke="currentColor" className="w-4 h-4">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12"/>
                        </svg>
                    </Button>
                </div>
                <Typography variant="body2" color="blueGray" className="mb-4">
                    Enter the data for the new Role.
                </Typography>
                <form onSubmit={handleSubmit}>
                    {roles.map((role, index) => (
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
                                    value={role.name}
                                    name="name"
                                    onChange={(e) => handleRoleChange(index, 'name', e.target.value)}
                                    className="mb-4"
                                    // onChange={handleRoleNameChange}
                                    // className="mb-4  w-full" // 添加边距
                                    error={!!nameError}
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
                                    value={role.description}
                                    name="description"
                                    onChange={(e) => handleRoleChange(index, 'description', e.target.value)}
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
                                            onClick={() => handleRemoveRole(index)}
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
                        onClick={handleAddRole}
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

export default AddRole;