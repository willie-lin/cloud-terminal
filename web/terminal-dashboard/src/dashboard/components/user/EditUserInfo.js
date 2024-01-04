import React, {useEffect, useRef, useState} from "react";
import {Button, Card, CardBody, CardHeader, Input, Typography} from "@material-tailwind/react";
import {getUserByEmail, uploadFile} from "../../../api/api";

function EditUserInfo({ email, onUpdate }) {
    const [userInfo, setUserInfo] = useState(null);
    const [newInfo, setNewInfo] = useState(userInfo);
    const [file, setFile] = useState(null); // 添加一个新的 state 来存储文件
    const [preview, setPreview] = useState(null); // 用于存储预览图片的 URL
    const fileInputRef = useRef(); // 用于访问文件输入元素

    // 获取用户信息
    useEffect(() => {
        if (email) {  // 添加这一行来检查email是否存在
            getUserByEmail(email)
                .then(data => setUserInfo(data))
                .catch(error => console.error('Error:', error));
        }
    }, [email]);
    // 当 userInfo 改变时，更新 newInfo
    useEffect(() => {
        setNewInfo(userInfo);
    }, [userInfo]);

    const handleChange = (event) => {
        setNewInfo({
            ...newInfo,
            [event.target.name]: event.target.value,
        });
    };
    const handleFileChange = (event) => { // 当文件改变时，更新 file state
        setFile(event.target.files[0]);
        if (file && file.type.startsWith('image/')) {
            // 创建 FileReader 对象来读取此文件
            const reader = new FileReader();
            reader.onloadend = () => {
                // 使用 canvas 来裁剪图片
                const img = document.createElement('img');
                img.onload = () => {
                    const canvas = document.createElement('canvas');
                    const ctx = canvas.getContext('2d');
                    // 设置 canvas 的宽度和高度
                    canvas.width = 90;
                    canvas.height = 90;
                    // 在 canvas 上绘制图片
                    ctx.drawImage(img, 0, 0, 90, 90);
                    // 获取裁剪后的图片
                    const croppedImage = canvas.toDataURL();
                    // 更新预览图片
                    setPreview(croppedImage);
                };
                img.src = reader.result;
            };
            reader.readAsDataURL(file);
        }
    };

    const handleSubmit = async (event) => { // 注意这里变成了异步函数
        event.preventDefault();
        if (file) {
            const filePath = await uploadFile(file); // 上传文件并获取文件路径
            onUpdate({ ...newInfo, avatar: filePath }); // 更新用户信息，包括新的头像路径
        } else {
            onUpdate(newInfo);
        }
    };

    return (
        <div className="flex justify-center items-center h-screen">
            <Card className="w-1/2">
                <CardHeader color="lightBlue">
                    <Typography color="black" style={{marginBottom: '0'}} className="font-bold text-center">
                        修改用户信息
                    </Typography>
                </CardHeader>


                <CardBody className="px-4 py-8">
                    <form onSubmit={handleSubmit}>
                        <label>
                            <img
                                className="w-24 h-24 rounded-full object-cover mb-6 mx-auto cursor-pointer"
                                src={preview || newInfo?.avatar || "https://i1.pngguru.com/preview/137/834/449/cartoon-cartoon-character-avatar-drawing-film-ecommerce-facial-expression-png-clipart.jpg"}
                                alt="Avatar Upload"
                            />
                            <input
                                type="file"
                                className="hidden"
                                multiple={false}
                                accept="image/*"
                                onChange={handleFileChange} // 当文件改变时，调用 handleFileChange 函数
                                ref={fileInputRef} // 添加 ref
                            />
                        </label>
                        {/*<div className="mb-10">*/}
                        {/*    <Input*/}
                        {/*        type="text"*/}
                        {/*        color="lightBlue"*/}
                        {/*        size="regular"*/}
                        {/*        outline={true}*/}
                        {/*        placeholder="头像URL"*/}
                        {/*        onChange={handleChange}*/}
                        {/*    />*/}
                        {/*</div>*/}

                        <div className="mb-1 flex flex-col gap-6">
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Nickname
                            </Typography>
                            <Input
                                type="text"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                placeholder="昵称"
                                onChange={handleChange}
                            />
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Email
                            </Typography>
                            <Input
                                type="email"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                placeholder="邮箱"
                                onChange={handleChange}
                            />
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Bio
                            </Typography>
                            <Input
                                type="text"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                placeholder="个人简介"
                                onChange={handleChange}
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
                                Submit
                            </Button>
                    </form>
                </CardBody>

            </Card>
        </div>
    );
}

export default EditUserInfo;
