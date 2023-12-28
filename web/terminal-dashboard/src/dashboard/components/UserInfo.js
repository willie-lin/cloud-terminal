// 用户信息组件
function UserInfo({ userInfo, currentTime}) {


    return (
        <div className="w-full max-w-2xl p-4 bg-white rounded shadow-md">
            {/*<div className="flex-grow w-full max-w-2xl p-4 bg-white rounded shadow-md">*/}
            <h1 className="text-4xl font-bold text-blue-600 mb-4">欢迎，{userInfo?.nickname}!</h1>
            {userInfo && <p className="text-xl text-blue-500">你的用户名是 {userInfo.username}。</p>}
            {userInfo && <p className="text-xl text-blue-500">你的电子邮件是 {userInfo.email}。</p>}
            <p className="text-xl text-blue-500">当前时间是 {currentTime.toLocaleTimeString()}。</p>
        </div>
    );
}

export default UserInfo;