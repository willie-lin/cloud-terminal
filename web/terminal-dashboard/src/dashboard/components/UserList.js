import {useEffect, useState} from "react";
import {getAllUsers} from "../../api/api";

function UserList({ email }) {
    const [users, setUsers] = useState([]);
    // 获取所有用户
    useEffect(() => {
        getAllUsers()
            .then(data => setUsers(data))
            .catch(error => console.error('Error:', error));
    }, []);


    return (
        <div className="flex-grow w-full mt-4 overflow-x-auto shadow-md sm:rounded-lg">
            <table className="min-w-full text-sm text-left text-gray-500 dark:text-gray-400">
                <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                <tr>
                    <th scope="col" className="px-6 py-3">ID</th>
                    <th scope="col" className="px-6 py-3">创建时间</th>
                    <th scope="col" className="px-6 py-3">更新时间</th>
                    <th scope="col" className="px-6 py-3">用户名</th>
                    <th scope="col" className="px-6 py-3">邮箱</th>
                    <th scope="col" className="px-6 py-3">在线状态</th>
                    <th scope="col" className="px-6 py-3">启用类型</th>
                    <th scope="col" className="px-6 py-3">最后登录时间</th>
                </tr>
                </thead>
                <tbody>
                {users.map(user => (
                    <tr key={user.id} className="bg-white border-b dark:bg-gray-800 dark:border-gray-700">
                        <td className="px-6 py-4">{user.id}</td>
                        <td className="px-6 py-4">{user.created_at}</td>
                        <td className="px-6 py-4">{user.updated_at}</td>
                        <td className="px-6 py-4">{user.username}</td>
                        <td className="px-6 py-4">{user.email}</td>
                        <td className="px-6 py-4">{user.online ? '在线' : '离线'}</td>
                        <td className="px-6 py-4">{user.enable_type ? '启用' : '禁用'}</td>
                        <td className="px-6 py-4">{user.last_login_time}</td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
}

export default UserList;