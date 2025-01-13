import {useContext, useState} from "react";
import {
    ArrowDownTrayIcon,
    MagnifyingGlassIcon,
} from "@heroicons/react/24/outline";
import {
    Card,
    CardHeader,
    Typography,
    Button,
    CardBody,
    CardFooter,
    Input,
} from "@material-tailwind/react";
import {UserPlusIcon} from "@heroicons/react/16/solid";
import {AuthContext} from "../../../App";
import { saveAs } from 'file-saver';
import AddPermission from "../permission/AddPermission";
import {useFetchAccessPolicies} from "./AccessPolicyHook";
import RenderAccessPolicy from "./RenderAccessPolicy";

function AccessPolicyList() {
    const { currentUser } = useContext(AuthContext);
    // 判断当前用户是否具有删除权限
    // const canDelete = currentUser?.roleName === 'Admin' || currentUser?.roleName === 'SuperAdmin'
    const canDelete = currentUser?.isTenantAdmin  || currentUser?.roleName === 'super_admin'


    const TABLE_HEAD = ["ID", "NAME", "DESCRIPTION", "STATEMENTS", "IMMUTABLE", "CREATED", "UPDATED", ""];


    const accessPolicies = useFetchAccessPolicies() || [];

    // 创建一个函数来生成 CSV 文件的内容
    const createCsvContent = (accessPolicies) => {
        // CSV 文件的头部
        const header = ["ID", "NAME", "DESCRIPTION", "STATEMENTS", "IMMUTABLE", "CREATED", "UPDATED", ""];

        // CSV 文件的数据
        const data = accessPolicies.map(accessPolicy => [
            accessPolicy.id,
            accessPolicy.name,
            accessPolicy.description,
            accessPolicy.statements,
            accessPolicy.created_at,
            accessPolicy.updated_at,
        ]);
        // 将头部和数据合并，然后转换为 CSV 格式
        return [header, ...data].map(row => row.join(',')).join('\n');
    };

    // 创建一个函数来处理 "Download" 按钮的点击事件
    const handleDownload = () => {
        // 生成 CSV 文件的内容
        const csvContent = createCsvContent(accessPolicies);
        // 创建一个 Blob 对象
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        // 使用 file-saver 库来保存文件
        saveAs(blob, 'AccessPolicy.csv');
    };


    const [page, setPage] = useState(1);
    const accessPoliciesPerPage = 10;
    const totalPages = Math.ceil(accessPolicies.length / accessPoliciesPerPage);

    const indexOfLastAccessPolicy = page * accessPoliciesPerPage;
    const indexOfFirstAccessPolicy = indexOfLastAccessPolicy - accessPoliciesPerPage;

    let currentAccessPolicies = accessPolicies.slice(indexOfFirstAccessPolicy, indexOfLastAccessPolicy);
    if (!Array.isArray(currentAccessPolicies)) {
        currentAccessPolicies = [];
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



    // 添加role
    const [isAddAccessPolicyOpen, setIsAddAccessPolicyOpen] = useState(false);

    const handleAddAccessPolicy = () => {
        setIsAddAccessPolicyOpen(false);
    };

    const openAddAccessPolicy = () => {
        setIsAddAccessPolicyOpen(true);
    };

    const closeAddAccessPolicy = () => {
        setIsAddAccessPolicyOpen(false);
    };
    return (
        <Card className="h-full w-full">
            <CardHeader floated={false} shadow={false} className="rounded-none">
                <div className="mb-4 flex flex-col justify-between gap-8 md:flex-row md:items-center">
                    <div>
                        <Typography variant="h5" color="blue-gray">
                            AccessPolicy List
                        </Typography>
                        <Typography color="gray" className="mt-1 font-normal">
                            These are details about the last Access Policy.
                        </Typography>
                    </div>
                    {canDelete && (
                        <div className="flex w-full shrink-0 gap-2 md:w-max">
                            <div className="w-full md:w-72">
                                <Input
                                    label="Search"
                                    icon={<MagnifyingGlassIcon className="h-5 w-5" />}
                                />
                            </div>
                            <Button className="flex items-center gap-3" size="sm"  onClick={handleDownload}>
                                {/*<Button className="flex items-center gap-3" size="sm"  >*/}
                                <ArrowDownTrayIcon strokeWidth={2} className="h-4 w-4" /> Download
                            </Button>
                            <Button className="flex items-center gap-3" size="sm" onClick={openAddAccessPolicy}>
                                <UserPlusIcon strokeWidth={2} className="h-4 w-4"/> Add Policy
                            </Button>
                        </div>
                    )}
                    {isAddAccessPolicyOpen && (
                        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
                             onClick={(e) => {
                                 // 如果事件的目标是这个容器本身，那么关闭模态窗口
                                 if (e.target === e.currentTarget) {
                                     closeAddAccessPolicy();
                                 }
                             }}
                        >
                            <AddPermission onAddPermission={handleAddAccessPolicy} onClose={closeAddAccessPolicy}/>
                        </div>
                    )}
                </div>
            </CardHeader>
            <CardBody className="overflow-scroll px-0">
                <table className="w-full min-w-max table-auto text-left">
                    <thead>
                    <tr>
                        {TABLE_HEAD.map((head) => (
                            <th
                                key={head}
                                className="border-y border-blue-gray-100 bg-blue-gray-50/50 p-4"
                            >
                                <Typography
                                    variant="small"
                                    color="blue-gray"
                                    className="font-normal leading-none opacity-70"
                                >
                                    {head}
                                </Typography>
                            </th>
                        ))}
                    </tr>
                    </thead>
                    <tbody>
                    {currentAccessPolicies.map(( accessPolicy, index) =>
                        <RenderAccessPolicy accessPolicy={accessPolicy} isLast={index === currentAccessPolicies.length - 1} />
                    )}
                    </tbody>
                </table>
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

export default AccessPolicyList;