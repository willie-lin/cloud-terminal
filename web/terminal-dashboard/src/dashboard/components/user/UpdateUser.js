import React, { useState } from 'react';
import {
    Button,
    Card,
    CardBody,
    CardHeader,
    Input, Option, Select,
    Textarea,
    Typography
} from '@material-tailwind/react';
import {updateUser} from "../../../api/api";

function UpdateUser({ user, onUpdateUser, onClose}) {

    // 使用user的值来初始化你的状态
    const [username, setUsername] = useState(user ? user.username : '');
    const [nickname, setNickname] = useState(user ? user.nickname : '');
    const [phone, setPhone] = useState(user ? user.phone : '');
    const [bio, setBio] = useState(user ? user.bio : '');
    const [onlineStatus, setOnlineStatus] = useState(user ? user.online : '');
    const [enableType, setEnableType] =  useState(user ? user.enable_type : '');
    const [inputError, setInputError] = useState(false);
    const MAX_LENGTH = 180; // 设置最大长度为200

    async function handleSubmit(e) {
        e.preventDefault();  // 阻止表单的默认提交行为
        try {
            // 创建一个对象来存储你的表单数据
            const data = {
                username: username,
                nickname: nickname,
                phone: phone,
                bio: bio,
                onlineStatus: onlineStatus,
                enableType: enableType,
            };
            await updateUser(data);
        } catch (error) {
            console.error(error);
    }
    }

    return (
        <Card className="w-96">
            <CardHeader variant="gradient" color="gray" className="mb-4 grid h-20 place-items-center">
                <Typography variant="h4" color="white">
                        Update User
                    </Typography>
                </CardHeader>
                <CardBody className="px-4 py-4">
                    <form onSubmit={handleSubmit}>
                        <div className="mb-1 flex flex-col gap-6">
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                ID
                            </Typography>
                            <Input
                                variant="outlined"
                                label="Disabled"
                                disabled
                                type="username"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                value={username}
                                name="username"  // 添加name属性
                                onChange={e => setUsername(e.target.value)}
                            />
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Nickname
                            </Typography>
                            <Input
                                variant="outlined"
                                label="Nickname"
                                type="nickname"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                value={nickname}
                                name="nickname"  // 添加name属性
                                onChange={e => setNickname(e.target.value)}
                            />
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Phone
                            </Typography>
                            <Input
                                variant="outlined"
                                label="Phone"
                                type="phone"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                value={phone}
                                name="phone"  // 添加name属性
                                onChange={e => setPhone(e.target.value)}

                            />
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Bio
                            </Typography>
                            <Textarea
                                variant="outlined"
                                label="Bio"
                                type="text"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                className={`border-${inputError ? 'red-500' : 'blue-500'}`}
                                value={bio}
                                name="bio"  // 添加name属性
                                // onChange={e => setBio(e.target.value)}
                                onChange={e =>{
                                    const value = e.target.value;
                                    if (value.length > MAX_LENGTH) {
                                        setInputError(true);
                                    } else {
                                        setInputError(false);
                                    }
                                    setBio(e.target.value)}
                                }
                            />
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Online Status
                            </Typography>
                            <Select size="md"
                                    label="OnlineStatus"
                                    value={onlineStatus ? 'true' : 'false'}
                                    onChange={value => setOnlineStatus(value === 'true')} >
                                <Option value='true'>True</Option>
                                <Option value='false'>False</Option>
                            </Select>
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Enable Type
                            </Typography>
                            <Select size="md"
                                    label="EnableType"
                                    value={enableType ? 'true' : 'false'}
                                    onChange={value => setEnableType(value === 'true')} >
                                <Option value='true'>True</Option>
                                <Option value='false'>False</Option>
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

export default UpdateUser;
