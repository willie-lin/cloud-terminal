import React, {useRef, useState} from "react";
import {Button, Card, CardBody, Input, Textarea, Typography} from "@material-tailwind/react";
import {editUserInfo, uploadFile} from "../../../api/api";
import {useNavigate} from "react-router-dom";
import {useFetchUserInfo} from "./UserHook";

function EditUserInfo({ user, onEditUser,  onUserChange, onClose }) {

    // 使用user的值来初始化你的状态
    const [email, setEmail] = useState(user ? user.email : '');
    const [avatar, setAvatar] = useState(user ? user.avatar : '');

    const [nickname, setNickname] = useState(user ? user.nickname : '');
    const [phone, setPhone] = useState(user ? user.phone_number : '');
    const [bio, setBio] = useState(user ? user.bio : '');

    const [file, setFile] = useState(null); // 添加一个新的 state 来存储文件
    const [preview, setPreview] = useState(null); // 用于存储预览图片的 URL
    const fileInputRef = useRef(); // 用于访问文件输入元素

    const MAX_LENGTH = 180; // 设置最大长度为180
    const [inputError, setInputError] = useState(false);

    // 在你的组件中
    const navigate = useNavigate();

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
    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            let avatar = user.avatar;  // 默认使用用户原来的头像
            if (file) {
                avatar = await uploadFile(file);  // 如果用户上传了新的头像，就使用新的头像
            }
                const data = {
                    email: email,  // 使用传递过来的email
                    avatar: avatar,
                    nickname: nickname,
                    phone: phone,
                    bio: bio
                };
                // console.log(data);  // 打印出data对象
                await editUserInfo(data);
               onEditUser()
            // navigate("/userinfo")
            navigate("/")
        }

        catch (error) {
            console.error(error);
        }
    };

    return (
            <Card className="w-1/3">
                <CardBody className="px-4 py-4">
                    <div className="flex justify-between items-center mb-2">
                        <Typography variant="h4" color="gray">
                            Edit User
                        </Typography>
                        <Button color="gray" buttonType="link" onClick={onClose}>
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5}
                                 stroke="currentColor" className="w-4 h-4">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12"/>
                            </svg>
                        </Button>
                    </div>
                    <Typography variant="body2" color="blueGray" className="mb-2">
                        Edit the data for the User.
                    </Typography>
                    <form onSubmit={handleSubmit}>
                        <label>
                            <img
                                className="w-24 h-24 rounded-full object-cover mb-6 mx-auto cursor-pointer"
                                src={preview || user?.avatar || "https://demos.creative-tim.com/test/corporate-ui-dashboard/assets/img/team-2.jpg"}
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
                                variant="outlined"
                                label="Nickname"
                                type="nickname"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                value={nickname}
                                name="nickname"  // 添加name属性
                                // onChange={handleChange}
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
                                // placeholder="Phone"
                                value={phone}
                                name="phone"  // 添加name属性
                                // onChange={handleChange}
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
                                // placeholder="Bio"
                                value={bio}
                                name="bio"  // 添加name属性
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

export default EditUserInfo;
