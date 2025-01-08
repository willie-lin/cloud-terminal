import {
    Button,
    Card,
    CardBody, CardFooter,
    CardHeader,
    Input,
    Tab,
    Tabs,
    TabsHeader,
    Typography
} from "@material-tailwind/react";
import {ChevronUpDownIcon, MagnifyingGlassIcon, UserPlusIcon} from "@heroicons/react/16/solid";
import {useFetchUsers} from "./UserHook";
import {useContext, useState} from "react";
import RenderUser from "./RenderUser";
import AddUserForm from "./AddUser";
import {ArrowDownTrayIcon} from "@heroicons/react/24/outline";
import { saveAs } from 'file-saver';
import {useTheme}  from '../../../layout/ThemeContext';
import {AuthContext} from "../../../App";


function UserList() {

    const { currentUser } = useContext(AuthContext);
    // 判断当前用户是否具有删除权限
    const canDelete = currentUser?.roleName === 'admin' || currentUser?.roleName === 'super_admin'

    const { isDarkMode } = useTheme();

    const TABS = [
        {
            label: "All",
            value: "all",
        },
        {
            label: "Monitored",
            value: "monitored",
        },
        {
            label: "Unmonitored",
            value: "unmonitored",
        },
    ];

    const [search, setSearch] = useState('');
    const [sortField, setSortField] = useState(null);
    const [sortDirection, setSortDirection] = useState('asc');
    const headers = ['ID', 'NICKNAME', 'USERNAME', 'EMAIL', 'PHONE', 'BIO', '2FA','STATUS', 'ONLINE', 'CREATED', 'UPDATED', 'LAST MODIFIED', 'EDIT USER', ''];

    const handleSort = (field) => {
        if (sortField === field) {
            // 如果已经在按这个字段排序，那么改变排序方向
            setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
        } else {
            // 否则，按这个字段进行升序排序
            setSortField(field);
            setSortDirection('asc');
        }
        // 在这里，你应该根据sortField和sortDirection来更新你的数据
    };

    const users = useFetchUsers() || [];
    const [page, setPage] = useState(1);
    const usersPerPage = 10;
    const totalPages = Math.ceil(users.length / usersPerPage);

    const indexOfLastUser = page * usersPerPage;
    const indexOfFirstUser = indexOfLastUser - usersPerPage;

    // const filteredUsers = users.filter(user => user.nickname.includes(search) || user.email.includes(search));
    const filteredUsers = users.filter(user =>
        (user.nickname && user.nickname.includes(search)) ||
        (user.email && user.email.includes(search)) ||
        (user.username && user.username.includes(search))
    );
    // 添加新的状态
    const [isViewAll, setIsViewAll] = useState(false);

    // 更新 currentUsers 的定义
    let currentUsers;
    if (isViewAll) {
        currentUsers = users;
    } else {
        currentUsers = users.slice(indexOfFirstUser, indexOfLastUser);
    }

    // 在 "view all" 按钮的点击事件处理函数中
    const handleViewAll = () => {
        setIsViewAll(true);
    };

    const handlePrevious = () => {
        if (page > 1) {
            setPage(page - 1);
        }
    };
    const handleNext = () => {
        if (page < totalPages) {
            setPage(page + 1);
        }
    };


    // 创建一个函数来生成 CSV 文件的内容
    const createCsvContent = (users) => {
        // CSV 文件的头部
        const header = ['ID', 'NICKNAME', 'USERNAME', 'EMAIL', 'PHONE', 'BIO', '2FA','STATUS', 'ONLINE', 'CREATED', 'UPDATED', 'LAST MODIFIED'];
        // CSV 文件的数据
        const data = users.map(user => [
            user.id,
            user.nickname,
            user.username,
            user.email,
            user.phone_number,
            user.bio,
            user.totp_secret,
            user.status,
            user.online,
            user.enable_type,
            user.created_at,
            user.updated_at,
            user.last_login_time,
        ]);
        // 将头部和数据合并，然后转换为 CSV 格式
        return [header, ...data].map(row => row.join(',')).join('\n');
    };

// 创建一个函数来处理 "Download" 按钮的点击事件
    const handleDownload = () => {
        // 生成 CSV 文件的内容
        const csvContent = createCsvContent(users);
        // 创建一个 Blob 对象
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        // 使用 file-saver 库来保存文件
        saveAs(blob, 'users.csv');
    };

    // 添加用户
    const [isAddUserOpen, setIsAddUserOpen] = useState(false);

    const handleAddUser = () => {
        // 这里是处理新用户的代码，例如发送一个请求到后端服务
        setIsAddUserOpen(false);
        // 添加用户成功后，将新用户添加到用户列表中
    };

    const openAddUser = () => {
        setIsAddUserOpen(true);
    };

    const closeAddUser = () => {
        setIsAddUserOpen(false);
    };

//

    return (
        <Card className={`h-full w-full ${isDarkMode ? 'dark bg-gray-900 text-gray-100' : 'light bg-gray-100 text-gray-900'}`}>
            <CardHeader floated={false} shadow={false} className={`rounded-none ${isDarkMode ? 'dark bg-gray-900 text-gray-100' : 'light bg-gray-100 text-gray-900'}`}>
                <div className="mb-4 flex flex-col justify-between gap-8 md:flex-row md:items-center">
                    <div>
                        <Typography variant="h4" className={isDarkMode ? 'text-gray-100' : 'text-gray-900'}>
                            Users List
                        </Typography>
                        <Typography className={`mt-1 font-normal ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                            See information about all Users
                        </Typography>
                    </div>
                    {canDelete && (
                    <div className="flex shrink-0 flex-col gap-2 sm:flex-row">
                        <Button
                            className={`flex items-center gap-3 ${isDarkMode ? 'bg-blue-500 hover:bg-blue-600' : ''}`}
                            size="sm"
                            onClick={handleDownload}
                        >
                            <ArrowDownTrayIcon strokeWidth={2} className="h-4 w-4" /> Download
                        </Button>
                        <Button
                            variant="outlined"
                            size="sm"
                            onClick={handleViewAll}
                            className={isDarkMode ? 'text-gray-100 border-gray-100' : ''}
                        >
                            view all
                        </Button>
                        <Button
                            className={`flex items-center gap-3 ${isDarkMode ? 'bg-blue-500 hover:bg-blue-600' : ''}`}
                            size="sm"
                            onClick={openAddUser}
                        >
                            <UserPlusIcon strokeWidth={2} className="h-4 w-4"/> Add User
                        </Button>
                    </div>
                        )}
                </div>
                <div className="flex flex-col items-center justify-between gap-4 md:flex-row">
                    <Tabs
                        value="all"
                        className={`w-full md:w-max ${isDarkMode ? 'text-white-100' : 'text-gray-900'}`}
                    >
                        <TabsHeader
                            className={`${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}
                            indicatorProps={{
                                className: `${isDarkMode ? 'bg-blue-500' : 'bg-black'} shadow-none rounded-md transition-colors duration-300`,
                            }}
                        >
                            {TABS.map(({label, value}) => (
                                <Tab
                                    key={value}
                                    value={value}
                                    className={`
                            ${isDarkMode
                                        ? 'text-gray-300 hover:text-blue-400'
                                        : 'text-gray-700 hover:text-black'
                                    }
                            px-4 py-2 rounded-md transition-colors duration-300
                            relative z-0
                        `}
                                >
                                    <span className="relative z-10">{label}</span>
                                </Tab>
                            ))}
                        </TabsHeader>
                    </Tabs>

                    <div className="w-full md:w-72">
                        <Input
                            label="Search"
                            icon={<MagnifyingGlassIcon className={`h-4 w-4 ${isDarkMode ? 'text-gray-400' : ''}`}/>}
                            value={search}
                            onChange={e => setSearch(e.target.value)}
                            className={isDarkMode ? 'text-gray-100 bg-gray-800' : ''}
                        />
                    </div>
                </div>
            </CardHeader>
            <CardBody className={`overflow-auto justify-between px-4 flex-grow mt-0 ${isDarkMode ? 'bg-gray-900' : ''}`}>
                <div className="overflow-auto" style={{maxHeight: 'calc(100vh - 200px)'}}>
                    <table className={`mt-1 w-full min-w-max table-auto text-left ${isDarkMode ? 'text-gray-100' : ''}`}>
                        <thead className="sticky top-0">
                        <tr>
                            {headers.map((head, index) => (
                                <th key={head}
                                    className={`cursor-pointer border-y p-2 transition-colors ${
                                        isDarkMode
                                            ? 'border-gray-700 bg-gray-800 hover:bg-gray-700'
                                            : 'border-blue-gray-100 bg-blue-gray-50/50 hover:bg-blue-gray-50'
                                    }`}
                                    onClick={() => handleSort(head)}
                                >
                                    <Typography
                                        variant="small"
                                        className={`flex items-center justify-between gap-2 font-normal leading-none opacity-70 ${
                                            isDarkMode ? 'text-gray-100' : 'text-gray-900'
                                        }`}
                                    >
                                        {head}
                                        {index !== 7 && (
                                            <ChevronUpDownIcon strokeWidth={2} className="h-4 w-4"/>
                                        )}
                                    </Typography>
                                </th>
                            ))}
                        </tr>
                        </thead>
                        <tbody>
                        {(search ? filteredUsers : currentUsers).map((user, index) => (
                            <RenderUser
                                key={user.id}
                                user={user}
                                isLast={index === (search ? filteredUsers : currentUsers).length - 1}
                                isDarkMode={isDarkMode}
                            />
                        ))}
                        </tbody>
                    </table>
                </div>
            </CardBody>
            {!isViewAll && (
                <CardFooter className={`flex items-center justify-between border-t p-0 ${
                    isDarkMode ? 'border-gray-700 bg-gray-900' : 'border-blue-gray-50'
                }`}>
                    <Typography variant="small" className={`font-normal mx-2 ${isDarkMode ? 'text-gray-100' : 'text-gray-900'}`}>
                        Page {page} of {totalPages}
                    </Typography>
                    <div className="flex gap-2">
                        <Button
                            variant="outlined"
                            size="sm"
                            onClick={handlePrevious}
                            className={isDarkMode ? 'text-gray-100 border-gray-100' : ''}
                        >
                            Previous
                        </Button>
                        <Button
                            variant="outlined"
                            size="sm"
                            onClick={handleNext}
                            className={isDarkMode ? 'text-gray-100 border-gray-100' : ''}
                        >
                            Next
                        </Button>
                    </div>
                </CardFooter>
            )}
            {isAddUserOpen && (
                <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
                     onClick={(e) => {
                         if (e.target === e.currentTarget) {
                             closeAddUser();
                         }
                     }}
                >
                    <AddUserForm onAddUser={handleAddUser} onClose={closeAddUser} isDarkMode={isDarkMode} />
                </div>
            )}
        </Card>
    );
}
export default UserList;