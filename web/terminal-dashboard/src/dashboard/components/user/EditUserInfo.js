import React, {useEffect, useRef, useState} from "react";
import {Button, Card, CardBody, CardHeader, Input, Typography} from "@material-tailwind/react";
import {confirm2FA, editUserInfo, uploadFile} from "../../../api/api";
import {useLocation} from "react-router-dom";

function EditUserInfo({ email }) {
    // 获取用户信息
    const location = useLocation();
    const userInfo = location.state.userInfo;
    const [newInfo, setNewInfo] = useState(userInfo);
    const [nickname, setNickname] = useState(null)
    const [phone, setPhone] = useState(null)
    const [bio, setBio] = useState(null)


    const [file, setFile] = useState(null); // 添加一个新的 state 来存储文件
    const [preview, setPreview] = useState(null); // 用于存储预览图片的 URL
    const fileInputRef = useRef(); // 用于访问文件输入元素

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

    const handleFileChange = (event) => {
        const file = event.target.files[0];
        setFile(file);
        if (file && file.type.startsWith('image/')) {
            const reader = new FileReader();
            reader.onloadend = () => {
                const img = document.createElement('img');
                img.onload = () => {
                    const canvas = document.createElement('canvas');
                    const ctx = canvas.getContext('2d');
                    // 设置 canvas 的宽度和高度
                    canvas.width = 90;
                    canvas.height = 90;
                    // 计算裁剪的起始位置
                    const startX = img.width > img.height ? (img.width - img.height) / 2 : 0;
                    const startY = img.height > img.width ? (img.height - img.width) / 2 : 0;
                    // 计算裁剪的宽度和高度
                    const sideLength = Math.min(img.width, img.height);
                    // 在 canvas 上绘制图片
                    ctx.drawImage(img, startX, startY, sideLength, sideLength, 0, 0, canvas.width, canvas.height);
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

    const handleSubmit = async (event) => {
        event.preventDefault();

        let updatedInfo = {...newInfo }; // 这里将newInfo的属性添加到userInfo中
        if (file) {
            const filePath = await uploadFile(file);
            updatedInfo = { ...updatedInfo, avatar: filePath};
        }
        try {
            const data = {
                email: email,  // 使用传递过来的email
                avatar: updatedInfo.avatar,
                nickname:updatedInfo.nickname,
                phone: updatedInfo.phone,
                bio: updatedInfo.bio
            };
            // console.log(data);  // 打印出data对象
            await editUserInfo(data);
        } catch (error) {
            console.error(error);
        }
    };

    return (
        <div className="flex justify-center items-center h-screen">
            <Card className="w-1/2">
                <CardHeader color="lightBlue">
                    <Typography color="black" style={{marginBottom: '0'}} className="font-bold text-center">
                        Edit UserInfo
                    </Typography>
                </CardHeader>
                <CardBody className="px-4 py-8">
                    <form onSubmit={handleSubmit}>
                        <label>
                            <img
                                className="w-24 h-24 rounded-full object-cover mb-6 mx-auto cursor-pointer"
                                src={preview || newInfo?.avatar || "https://i1.pngguru.com/preview/137/834/449/cartoon-cartoon-character-avatar-drawing-film-ecommerce-facial-expression-png-clipart.jpg"}
                                alt="Avatar Upload"
                                style={{width: '120px', height: '120px', border: '1px solid'}} // Added border here
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
                        <div className="mb-1 flex flex-col gap-6">
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Nickname
                            </Typography>
                            <Input
                                type="text"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                placeholder="Nickname"
                                value={nickname}
                                name="nickname"  // 添加name属性
                                onChange={handleChange}
                            />
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Phone
                            </Typography>
                            <Input
                                type="phone"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                placeholder="Phone"
                                value={phone}
                                name="phone"  // 添加name属性
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
                                placeholder="Bio"
                                value={bio}
                                name="bio"  // 添加name属性
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
