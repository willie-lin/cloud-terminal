// import {useContext, useState} from "react";
// import {
//     ArrowDownTrayIcon,
//     MagnifyingGlassIcon,
// } from "@heroicons/react/24/outline";
// import {
//     Card,
//     CardHeader,
//     Typography,
//     Button,
//     CardBody,
//     CardFooter,
//     Input,
// } from "@material-tailwind/react";
// import {UserPlusIcon} from "@heroicons/react/16/solid";
// import AddPermission from "./AddPermission";
// import {useFetchPermissions} from "./PermissionHook";
// import RenderPermission from "./RenderPermission";
// import {AuthContext} from "../../../App";
// import { saveAs } from 'file-saver';
//
// function PolicyList() {
//     const { currentUser } = useContext(AuthContext);
//     // 判断当前用户是否具有删除权限
//     const canDelete = currentUser?.roleName === 'Admin' || currentUser?.roleName === 'SuperAdmin'
//
//
//     const TABLE_HEAD = ["ID", "NAME", "DESCRIPTION", "CREATED", "LASTMODIFIED", ""];
//
//
//     const permissions = useFetchPermissions() || [];
//
//     // 创建一个函数来生成 CSV 文件的内容
//     const createCsvContent = (permissions) => {
//         // CSV 文件的头部
//         const header = ["ID", "NAME", "ACTIONS", "RESOURCE_TYPE", "DESCRIPTION", "DEFAULT",  "CREATED", "LASTMODIFIED", ""];
//         // CSV 文件的数据
//         const data = permissions.map(permission => [
//             permission.id,
//             permission.name,
//             permission.actions,
//             permission.description,
//             permission.resource_type,
//             permission.is_disabled,
//             permission.created_at,
//             permission.updated_at,
//         ]);
//         // 将头部和数据合并，然后转换为 CSV 格式
//         return [header, ...data].map(row => row.join(',')).join('\n');
//     };
//
//     // 创建一个函数来处理 "Download" 按钮的点击事件
//     const handleDownload = () => {
//         // 生成 CSV 文件的内容
//         const csvContent = createCsvContent(permissions);
//         // 创建一个 Blob 对象
//         const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
//         // 使用 file-saver 库来保存文件
//         saveAs(blob, 'permissions.csv');
//     };
//
//
//     const [page, setPage] = useState(1);
//     const permissionsPerPage = 10;
//     const totalPages = Math.ceil(permissions.length / permissionsPerPage);
//
//     const indexOfLastPermission = page * permissionsPerPage;
//     const indexOfFirstPermission = indexOfLastPermission - permissionsPerPage;
//
//     let currentPermissions = permissions.slice(indexOfFirstPermission, indexOfLastPermission);
//     if (!Array.isArray(currentPermissions)) {
//         currentPermissions = [];
//     }
//     const handlePrevious = () => {
//         if (page > 1) {
//             setPage(page - 1);
//         }
//     };
//     const handleNext = () => {
//         if (page < totalPages) {
//             setPage(page + 1);
//         }
//     };
//
//
//
//     // 添加role
//     const [isAddPermissionOpen, setIsAddPermissionOpen] = useState(false);
//
//     const handleAddPermission = () => {
//         setIsAddPermissionOpen(false);
//     };
//
//     const openAddPermission = () => {
//         setIsAddPermissionOpen(true);
//     };
//
//     const closeAddPermission = () => {
//         setIsAddPermissionOpen(false);
//     };
//     return (
//         <Card className="h-full w-full">
//             <CardHeader floated={false} shadow={false} className="rounded-none">
//                 <div className="mb-4 flex flex-col justify-between gap-8 md:flex-row md:items-center">
//                     <div>
//                         <Typography variant="h5" color="blue-gray">
//                             Permission List
//                         </Typography>
//                         <Typography color="gray" className="mt-1 font-normal">
//                             These are details about the last Permissions.
//                         </Typography>
//                     </div>
//                     {canDelete && (
//                         <div className="flex w-full shrink-0 gap-2 md:w-max">
//                             <div className="w-full md:w-72">
//                                 <Input
//                                     label="Search"
//                                     icon={<MagnifyingGlassIcon className="h-5 w-5" />}
//                                 />
//                             </div>
//                             <Button className="flex items-center gap-3" size="sm"  onClick={handleDownload}>
//                                 {/*<Button className="flex items-center gap-3" size="sm"  >*/}
//                                 <ArrowDownTrayIcon strokeWidth={2} className="h-4 w-4" /> Download
//                             </Button>
//                             <Button className="flex items-center gap-3" size="sm" onClick={openAddPermission}>
//                                 <UserPlusIcon strokeWidth={2} className="h-4 w-4"/> Add Permission
//                             </Button>
//                         </div>
//                     )}
//                     {isAddPermissionOpen && (
//                         <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50"
//                              onClick={(e) => {
//                                  // 如果事件的目标是这个容器本身，那么关闭模态窗口
//                                  if (e.target === e.currentTarget) {
//                                      closeAddPermission();
//                                  }
//                              }}
//                         >
//                             <AddPermission onAddPermission={handleAddPermission} onClose={closeAddPermission}/>
//                         </div>
//                     )}
//                 </div>
//             </CardHeader>
//             <CardBody className="overflow-scroll px-0">
//                 <table className="w-full min-w-max table-auto text-left">
//                     <thead>
//                     <tr>
//                         {TABLE_HEAD.map((head) => (
//                             <th
//                                 key={head}
//                                 className="border-y border-blue-gray-100 bg-blue-gray-50/50 p-4"
//                             >
//                                 <Typography
//                                     variant="small"
//                                     color="blue-gray"
//                                     className="font-normal leading-none opacity-70"
//                                 >
//                                     {head}
//                                 </Typography>
//                             </th>
//                         ))}
//                     </tr>
//                     </thead>
//                     <tbody>
//                     {currentPermissions.map(( permission, index) =>
//                         <RenderPermission permission={permission} isLast={index === currentPermissions.length - 1} />
//                     )}
//                     </tbody>
//                 </table>
//             </CardBody>
//             <CardFooter className="flex items-center justify-between border-t border-blue-gray-50 p-0">
//                 <Typography variant="small" color="blue-gray" className="font-normal mx-2">
//                     Page {page} of {totalPages}
//                 </Typography>
//                 <div className="flex gap-2">
//                     <Button variant="outlined" size="sm" onClick={handlePrevious}>
//                         Previous
//                     </Button>
//                     <Button variant="outlined" size="sm" onClick={handleNext}>
//                         Next
//                     </Button>
//                 </div>
//             </CardFooter>
//         </Card>
//     );
// }
//
// export default PolicyList;