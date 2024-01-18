import {Avatar, Chip, IconButton, Tooltip, Typography} from "@material-tailwind/react";
import {PencilIcon} from "@heroicons/react/16/solid";

function UserRow({ user, isLast }) {
    const classes = isLast
        ? "p-2"
        : "p-2 border-b border-blue-gray-50";

    return (
        <tr key={user.id}>
            <td className={classes}>
                <div className="flex flex-col">
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {user.id}
                    </Typography>
                </div>
            </td>
            <td className={classes}>
                <div className="flex items-center gap-3">
                    {/* 这里可以添加Avatar组件 */}
                    <Avatar src={user.avatar} alt={user.nickname} size="sm"/>
                    <div className="flex flex-col">
                        <Typography variant="small" color="blue-gray"
                                    className="font-normal">
                            {user.nickname}
                        </Typography>
                        <Typography variant="small" color="blue-gray"
                                    className="font-normal">
                            {user.email}
                        </Typography>
                    </div>
                </div>
            </td>
            <td className={classes}>
                <div className="flex flex-col">
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {user.username}
                    </Typography>
                </div>
            </td>
            <td className={classes}>
                <div className="flex flex-col">
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {user.email}
                    </Typography>
                </div>
            </td>
            <td className={classes}>
                <div className="flex flex-col">
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {user.phone}
                    </Typography>
                </div>
            </td>
            <td className={classes}>
                <div className="flex flex-col">
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {user.bio}
                    </Typography>
                </div>
            </td>
            <td className={classes}>
                <div className="flex flex-col">
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {user.totp_secret}
                    </Typography>
                </div>
            </td>
            <td className={classes}>
                <div className="w-max">
                    <Chip variant="ghost" size="sm" value={user.online ? "在线" : "离线"}
                          color={user.online ? "green" : "blue-gray"}
                    />
                </div>
            </td>
            <td className={classes}>
                <div className="flex flex-col">
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {user.enable_type ? '启用' : '禁用'}
                    </Typography>
                </div>
            </td>
            <td className={classes}>
                <div className="flex flex-col">
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {user.created_at}
                    </Typography>
                </div>
            </td>
            <td className={classes}>
                <div className="flex flex-col">
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {user.updated_at}
                    </Typography>
                </div>
            </td>
            <td className={classes}>
                <div className="flex flex-col">
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {user.last_login_time}
                    </Typography>
                </div>
            </td>
            <td className={classes}>
                <Tooltip content="修改用户">
                    <IconButton variant="text">
                        <PencilIcon className="h-4 w-4"/>
                    </IconButton>
                </Tooltip>
            </td>
        </tr>
    );
}
export default UserRow;

// 然后在你的代码中使用这个组件：
