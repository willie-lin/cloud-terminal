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
    const headers = ['ID', '昵称', '用户名', '邮箱', '手机号', '个人简介', '2FA','在线状态', '启用类型', '创建时间', '更新时间', '最后登录时间'];

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

    const users = useFetchUsers();
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

    let currentUsers = users.slice(indexOfFirstUser, indexOfLastUser);
    if (!Array.isArray(currentUsers)) {
        currentUsers = [];
    }
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
            <Card className="h-full w-full">
                <CardHeader floated={true} shadow={false} className="rounded-none">
                    <div className="flex items-center justify-between gap-4 mb-1">
                        <div>
                            <Typography variant="h4" color="blue-gray">
                                Users list
                            </Typography>
                            <Typography color="gray" className="mt-1 font-normal">
                                See information about all users
                            </Typography>
                        </div>
                        <div className="flex shrink-0 flex-col gap-2 sm:flex-row">
                            <Button variant="outlined" size="sm">
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
            </Card>
);
}
export default UserList;