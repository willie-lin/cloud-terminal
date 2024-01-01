import {useEffect, useState} from "react";
import {Button, Card, CardBody, CardHeader, Input, Typography} from "@material-tailwind/react";
import {getUserByEmail} from "../../../api/api";

function EditUserInfo({ email, onUpdate }) {
    const [userInfo, setUserInfo] = useState(null);

    // 获取用户信息
    useEffect(() => {
        if (email) {  // 添加这一行来检查email是否存在
            getUserByEmail(email)
                .then(data => setUserInfo(data))
                .catch(error => console.error('Error:', error));
        }
    }, [email]);


    const [newInfo, setNewInfo] = useState(userInfo);

    const handleChange = (event) => {
        setNewInfo({
            ...newInfo,
            [event.target.name]: event.target.value,
        });
    };

    const handleSubmit = (event) => {
        event.preventDefault();
        onUpdate(newInfo);
    };

    return (
        <div className="flex justify-center items-center h-screen">
            <Card className="w-1/2">
                <CardHeader color="lightBlue">
                    <Typography color="black" style={{ marginBottom: '0' }} className="font-bold text-center">
                        修改用户信息
                    </Typography>
                </CardHeader>

                <CardBody>
                    <div className="mb-8 px-4">
                        <Input
                            type="text"
                            color="lightBlue"
                            size="regular"
                            outline={true}
                            placeholder="头像URL"
                            onChange={handleChange}
                        />
                    </div>
                    <div className="mb-8 px-4">
                        <Input
                            type="text"
                            color="lightBlue"
                            size="regular"
                            outline={true}
                            placeholder="昵称"
                            onChange={handleChange}
                        />
                    </div>
                    <div className="mb-8 px-4">
                        <Input
                            type="email"
                            color="lightBlue"
                            size="regular"
                            outline={true}
                            placeholder="邮箱"
                            onChange={handleChange}
                        />
                    </div>
                    <div className="mb-8 px-4">
                        <Input
                            type="text"
                            color="lightBlue"
                            size="regular"
                            outline={true}
                            placeholder="个人简介"
                            onChange={handleChange}
                        />
                    </div>
                    <div className="flex justify-center px-4 mb-4">
                        <Button color="lightBlue" ripple="light">
                            提交
                        </Button>
                    </div>
                </CardBody>
            </Card>
        </div>
    );
}
export default EditUserInfo;
