// 用户信息组件
import {useEffect, useState} from "react";
import {getUserByEmail} from "../../../api/api";
import {
    Avatar,
    Card,
    CardBody,
    CardFooter,
    CardHeader,
    Tooltip,
    Typography
} from "@material-tailwind/react";

function UserInfo({ email, onUpdate }) {
    const [currentTime, setCurrentTime] = useState(new Date());
    const [userInfo, setUserInfo] = useState(null);

    // 获取用户信息
    useEffect(() => {
        if (email) {  // 添加这一行来检查email是否存在
            getUserByEmail(email)
                .then(data => setUserInfo(data))
                .catch(error => console.error('Error:', error));
        }
    }, [email]);

    // 更新当前时间
    useEffect(() => {
        const timer = setInterval(() => {
            setCurrentTime(new Date());
        }, 1000); // 每秒更新一次
        return () => {
            clearInterval(timer);
        };
    }, []);


    const [isEditing, setIsEditing] = useState(false);
    const [newInfo, setNewInfo] = useState({});

    const handleEdit = () => {
        setIsEditing(true);
        setNewInfo(userInfo);
    };

    const handleChange = (event) => {
        setNewInfo({
            ...newInfo,
            [event.target.name]: event.target.value,
        });
    };

    const handleSubmit = (event) => {
        event.preventDefault();
        onUpdate(newInfo);
        setIsEditing(false);
    };
    return (
        <Card className="w-full max-w-2xl p-4 bg-white rounded shadow-md">
        {/*<div className="flex-grow w-full max-w-2xl p-4 bg-white rounded shadow-md">*/}
        {/*    <Avatar src={userInfo.avatar} alt="avatar" />*/}
        {/*    <h1 className="text-blueGray-500">欢迎，{userInfo?.nickname}!</h1>*/}
        {/*    {userInfo && <p className="text-blueGray-500">你的用户名是 {userInfo.username}。</p>}*/}
        {/*    {userInfo && <p className="text-blueGray-500">你的电子邮件是 {userInfo.email}。</p>}*/}
        {/*    <p className="text-blueGray-500">当前时间是 {currentTime.toLocaleTimeString()}。</p>*/}
        {/*</div>*/}
            <CardHeader floated={false} className="h-80">
                <div>
                    {userInfo && <Avatar src={userInfo.avatar} alt="avatar" variant="rounded" size="xxl"/>}

                    {/*<img src="https://docs.material-tailwind.com/img/team-3.jpg" alt="profile-picture"/>*/}
                </div>
            </CardHeader>
            <CardBody className="text-center">
                <Typography variant="h4" color="blue-gray" className="mb-2">
                    Natalie Paisley
                </Typography>
                <Typography color="blue-gray" className="font-medium" textGradient>
                    CEO / Co-Founder
                </Typography>
            </CardBody>
            <CardFooter className="flex justify-center gap-7 pt-2">
                <Tooltip content="Like">
                    <Typography
                        as="a"
                        href="#facebook"
                        variant="lead"
                        color="blue"
                        textGradient
                    >
                        <i className="fab fa-facebook" />
                    </Typography>
                </Tooltip>
                <Tooltip content="Follow">
                    <Typography
                        as="a"
                        href="#twitter"
                        variant="lead"
                        color="light-blue"
                        textGradient
                    >
                        <i className="fab fa-twitter" />
                    </Typography>
                </Tooltip>
                <Tooltip content="Follow">
                    <Typography
                        as="a"
                        href="#instagram"
                        variant="lead"
                        color="purple"
                        textGradient
                    >
                        <i className="fab fa-instagram" />
                    </Typography>
                </Tooltip>
            </CardFooter>


            {/*{isEditing ? (*/}
            {/*    <Form onSubmit={handleSubmit}>*/}
        {/*        <p>昵称：</p>*/}
        {/*        <Input name="nickname" value={newInfo.nickname} onChange={handleChange}/>*/}
        {/*        <p>邮箱：</p>*/}
        {/*        <Input name="email" value={newInfo.email} onChange={handleChange}/>*/}
        {/*        /!* 其他字段 *!/*/}
        {/*        <Button type="submit" color="lightBlue">提交</Button>*/}
        {/*    </Form>*/}
        {/*) : (*/}
        {/*    <Button onClick={handleEdit} color="lightBlue">修改</Button>*/}
        {/*)}*/}
    </Card>

);
}

export default UserInfo;