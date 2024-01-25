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
            <CardBody className="px-4 py-4">
                <div className="flex justify-between items-center mb-2">
                    <Typography variant="h4" color="gray">
                        Delete User
                    </Typography>
                    <Button color="gray" buttonType="link" onClick={onClose}>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5}
                             stroke="currentColor" className="w-4 h-4">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12"/>
                        </svg>
                    </Button>
                </div>
                <Typography variant="body2" color="blueGray" className="mb-2">
                    Delete the data for the user.
                </Typography>
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
                            确定
                        </Button>
                </form>
            </CardBody>
        </Card>
    )
        ;
}

export default DeleteUserForm;
