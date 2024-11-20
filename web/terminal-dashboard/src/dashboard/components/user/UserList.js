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
import {useState} from "react";
import UserRow from "./UserRow";
import AddUserForm from "./AddUser";
import {ArrowDownTrayIcon} from "@heroicons/react/24/outline";
import { saveAs } from 'file-saver';

function UserList() {

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
    const headers = ['ID', 'NICKNAME', 'USERNAME', 'EMAIL', 'PHONE', 'BIO', '2FA','STATUS', 'ONLINE', 'CREATED', 'UPDATED', 'LAST MODIFIED', "", ""];

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



    // let currentUsers = users.slice(indexOfFirstUser, indexOfLastUser);
    // if (!Array.isArray(currentUsers)) {
    //     currentUsers = [];
    // }

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
            user.phone,
            user.bio,
            user.totp_secret,
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

    return (
            // <Card className="h-full w-full">
            //     <CardHeader floated={true} shadow={false} className="rounded-none">
            //         <div className="flex items-center justify-between gap-4 mb-1">
        <Card className="h-full w-full">
            <CardHeader floated={false} shadow={false} className="rounded-none">
                <div className="mb-4 flex flex-col justify-between gap-8 md:flex-row md:items-center">
                        <div>
                            <Typography variant="h4" color="blue-gray">
                                Users List
                            </Typography>
                            <Typography color="gray" className="mt-1 font-normal">
                                See information about all Users
                            </Typography>
                        </div>
                        <div className="flex shrink-0 flex-col gap-2 sm:flex-row">
                            <Button className="flex items-center gap-3" size="sm" onClick={handleDownload}>
                                <ArrowDownTrayIcon strokeWidth={2} className="h-4 w-4" /> Download
                            </Button>
                            <Button variant="outlined" size="sm" onClick={handleViewAll}>
                                view all
                            </Button>
                            <Button className="flex items-center gap-3" size="sm" onClick={openAddUser}>
                                <UserPlusIcon strokeWidth={2} className="h-4 w-4"/> Add User
                            </Button>
                        </div>
                        {isAddUserOpen && (
                            <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
                                 onClick={(e) => {
                                     // 如果事件的目标是这个容器本身，那么关闭模态窗口
                                     if (e.target === e.currentTarget) {
                                         closeAddUser();
                                     }
                                 }}
                            >
                                <AddUserForm onAddUser={handleAddUser} onClose={closeAddUser}/>
                            </div>
                        )}
                    </div>
                    <div className="flex flex-col items-center justify-between gap-4 md:flex-row">
                        <Tabs value="all" className="w-full md:w-max">
                            <TabsHeader>
                                {TABS.map(({label, value}) => (
                                    <Tab key={value} value={value}>
                                        &nbsp;&nbsp;{label}&nbsp;&nbsp;
                                    </Tab>
                                ))}
                            </TabsHeader>
                        </Tabs>
                        <div className="w-full md:w-72">
                            <Input
                                label="Search"
                                icon={<MagnifyingGlassIcon className="h-4 w-4"/>}
                                value={search}
                                onChange={e => setSearch(e.target.value)}
                            />
                        </div>
                    </div>
                </CardHeader>
                <CardBody className="overflow-auto justify-between px-4 flex-grow mt-0">
                    <div className="overflow-auto" style={{maxHeight: 'calc(100vh - 200px)'}}>
                        <table className="mt-1 w-full min-w-max table-auto text-left">
                            <thead className="sticky top-0">
                            <tr>
                                {headers.map((head, index) => (
                                    // <th key={head} className="p-4 border-b border-blue-gray-50">
                                    <th key={head}
                                        className="cursor-pointer border-y border-blue-gray-100 bg-blue-gray-50/50 p-2 transition-colors hover:bg-blue-gray-50"
                                        onClick={() => handleSort(head)}>
                                        <Typography variant="small" color="blue-gray"
                                                    className="flex items-center justify-between gap-2 font-normal leading-none opacity-70">
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
                                <UserRow user={user} isLast={index === (search ? filteredUsers : currentUsers).length - 1} />
                            ))}
                            </tbody>
                        </table>
                    </div>
                </CardBody>
                {!isViewAll && (
                    <CardFooter className="flex items-center justify-between border-t border-blue-gray-50 p-0">
                        <Typography variant="small" color="blue-gray" className="font-normal mx-2">
                            Page {page} of {totalPages}
                        </Typography>
                        <div className="flex gap-2">
                            <Button variant="outlined" size="sm" onClick={handlePrevious}>
                                Previous
                            </Button>
                            <Button variant="outlined" size="sm" onClick={handleNext}>
                                Next
                            </Button>
                        </div>
                    </CardFooter>
                )}
            </Card>
);
}
export default UserList;