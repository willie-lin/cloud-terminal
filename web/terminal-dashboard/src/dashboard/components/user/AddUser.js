import React, {useEffect, useState} from 'react';
import {Alert, Button, Card, CardBody, CardHeader, Input, Select, Typography, Option} from "@material-tailwind/react";
import {addUser, checkEmail, getAllRoles} from "../../../api/api";
import CryptoJS from "crypto-js";
import {useNavigate} from "react-router-dom";
import {useFetchRoles} from "../role/RoleHook";

function AddUserForm({ onAddUser, onClose }) {
    const [addUserError, setAddUserError] = useState('');
    const [password, setPassword] = useState('');
    const [emailError, setEmailError] = useState('');
    const [email, setEmail] = useState('');
    const [enableType, setEnableType] = useState('');
    const [onlineStatus, setOnlineStatus] = useState('');
    const [selectedRole, setSelectedRole] = useState(''); // State for selected role

    const roles = useFetchRoles(); // 使用自定义的 hook 获取角色列表

    const navigate = useNavigate();
    const handleEmailChange = async (e) => {
        const email = e.target.value;
        setEmail(email);
        try {
            const exists = await checkEmail(email);
            setEmailError(exists ? 'Email already registered' : '');
        } catch (error) {
            console.error(error);
        }
    };


    // const CryptoJS = require("crypto-js");
    const handleSubmit = async (e) => {
        e.preventDefault();
        // 验证电子邮件和密码是否已填写
        if (!email || !password) {
            setAddUserError('请填写所有必填字段');
            setTimeout(() => setAddUserError(''), 1000); // 1秒后清除错误信息
            return;
        }
        try {
            // 对密码进行哈希处理
            const hashedPassword = CryptoJS.SHA256(password).toString();
            const data = {
                email: email,  // 使用传递过来的email
                password: hashedPassword,
                roleID: selectedRole, // 添加角色到提交数据
                online: onlineStatus,
                enableType: enableType
            }
            await addUser(data); // 使用 register 函数
            // console.log(datas);
            onAddUser(email);
            navigate("/")
        } catch (error) {
            console.error(error);
        }
    };

    return (
            <Card className="w-1/3">
                <CardBody className="px-4 py-8">
                    <div className="flex justify-between items-center mb-4">
                        <Typography variant="h4" color="gray">
                            Create User
                        </Typography>
                        <Button color="gray" buttonType="link" onClick={onClose}>
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5}
                                 stroke="currentColor" className="w-4 h-4">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12"/>
                            </svg>
                        </Button>
                    </div>
                    <Typography variant="body2" color="blueGray" className="mb-4">
                        Enter the data for the new User.
                    </Typography>
                    <form onSubmit={handleSubmit}>
                        <div className="mb-1 flex flex-col gap-6">
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Email
                            </Typography>
                            <Input
                                variant="outlined"
                                label="Email"
                                type="email"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                value={email}
                                name="email"  // 添加name属性
                                onChange={handleEmailChange}
                                className="mb-4  w-full" // 添加边距
                                error={!!emailError}
                            />
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Password
                            </Typography>
                            <Input
                                variant="outlined"
                                label="Password"
                                type="password"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                value={password}
                                name="password"  // 添加name属性
                                // onChange={handleChange}
                                className="mb-4" // 添加边距
                                onChange={(e) => setPassword(e.target.value)}
                            />
                            {addUserError && (
                                <Alert color="red" className="mb-4">
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center">
                                            <i className="fas fa-info-circle mr-2"></i>
                                            <span className="text-sm">{addUserError}</span>
                                        </div>
                                    </div>
                                </Alert>
                            )}
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Role
                            </Typography>
                            <Select size="lg" label="Role" onChange={value => setSelectedRole(value)}>
                                {roles && roles.map((role) =>
                                    ( <Option key={role.id} value={role.id}>{role.name}</Option> ))}
                            </Select>
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Online Status
                            </Typography>
                            <Select size="md"
                                    label="OnlineStatus"
                                    onChange={value => setOnlineStatus(value)}>
                                <Option value={true}>True</Option>
                                <Option value={false}>False</Option>
                            </Select>
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Enable Type
                            </Typography>
                            <Select size="lg" label="EnableType" onChange={value => setEnableType(value)}>
                                <Option value={true}>True</Option>
                                <Option value={false}>False</Option>
                            </Select>
                        </div>
                        <Button fullWidth
                                type="submit"
                                color="lightBlue"
                                buttonType="filled"
                                size="regular"
                                rounded={false}
                                block={false}
                                iconOnly={false}
                                ripple="light"
                                className="mt-6" // 添加边距
                        >
                            Submit
                        </Button>
                    </form>
                </CardBody>
            </Card>
    );
}

export default AddUserForm;

