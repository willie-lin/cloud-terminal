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

function UpdateUser({ user, onUpdateUser, onClose }) {

    // 使用user的值来初始化你的状态
    const [username, setUsername] = useState(user ? user.username : '');
    const [nickname, setNickname] = useState(user ? user.nickname : '');
    const [phone, setPhone] = useState(user ? user.phone : '');
    const [bio, setBio] = useState(user ? user.bio : '');
    const [online, setOnline] = useState(user ? user.online : '');
    const [status, setStatus] =  useState(user ? user.status : '');
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
                online: online,
                status: status,
            };
            await updateUser(data);
            onUpdateUser()
        } catch (error) {
            console.error(error);
    }
    }

    return (
        <Card className="w-1/3">
            <CardBody className="px-4 py-4">
                <div className="flex justify-between items-center mb-2">
                    <Typography variant="h4" color="gray">
                        Update User
                    </Typography>
                    <Button color="gray" buttonType="link" onClick={onClose}>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5}
                             stroke="currentColor" className="w-4 h-4">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12"/>
                        </svg>
                    </Button>
                </div>
                <Typography variant="body2" color="blueGray" className="mb-2">
                    Update the data for the User.
                </Typography>
                <form onSubmit={handleSubmit}>
                    <div className="mb-1 flex flex-col gap-3.5">
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
                            onChange={e => {
                                const value = e.target.value;
                                if (value.length > MAX_LENGTH) {
                                    setInputError(true);
                                } else {
                                    setInputError(false);
                                }
                                setBio(e.target.value)
                            }
                            }
                        />
                        <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                            Online
                        </Typography>
                        <Select size="md"
                                label="Online"
                                onChange={value => setOnline(value)}>
                            <Option value={true}>True</Option>
                            <Option value={false}>False</Option>
                        </Select>
                        <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                            Status
                        </Typography>
                        {/*<Select size="md"*/}
                        {/*        label="Status"*/}
                        {/*        value={status ? 'true' : 'false'}*/}
                        {/*        onChange={value => setStatus(value === 'true')}>*/}
                        {/*    <Option value='true'>True</Option>*/}
                        {/*    <Option value='false'>False</Option>*/}
                        {/*</Select>*/}

                        <Select size="lg" label="Status" onChange={value => setStatus(value)}>
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

export default UpdateUser;
