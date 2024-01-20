import React, {useState} from 'react';
import {Button, Card, CardBody, CardHeader, Input, Typography} from "@material-tailwind/react";
import {deleteUser, } from "../../../api/api";

function DeleteUserForm({ user, onDeleteUser, onClose }) {
    // 使用user的值来初始化你的状态
    const [username, setUsername] = useState(user ? user.username : '');
    async function handleSubmit(e) {
        e.preventDefault();  // 阻止表单的默认提交行为
        if (window.confirm('Are you sure you want to delete this user?')) {
            try {
                const data = {
                    username: username,
                };
                await deleteUser(data);
                onDeleteUser();
            } catch (error) {
                console.error(error);
                alert('An error occurred while deleting the user.');
            }
        }
    }

    return (
        <Card className="w-96">
            <CardHeader variant="gradient" color="gray" className="mb-4 grid h-20 place-items-center">
                <Typography variant="h4" color="white">
                    Delete User
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
                    </div>
                    <div className="flex justify-between">
                        <Button
                            color="lightBlue"
                            buttonType="filled"
                            size="regular"
                            rounded={false}
                            block={false}
                            iconOnly={false}
                            ripple="light"
                            type="submit"
                        >
                            确定
                        </Button>
                        <Button
                            color="red"
                            buttonType="filled"
                            size="regular"
                            rounded={false}
                            block={false}
                            iconOnly={false}
                            ripple="light"
                            type="button"
                            onClick={onClose}
                        >
                            取消
                        </Button>

                    </div>
                </form>
            </CardBody>
        </Card>
    )
        ;
}

export default DeleteUserForm;
