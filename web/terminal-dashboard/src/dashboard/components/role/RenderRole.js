import {Chip, IconButton, Tooltip, Typography} from "@material-tailwind/react";
import {PencilIcon} from "@heroicons/react/16/solid";

function RenderRole({ role, isLast }) {
    const classes = isLast
        ? "p-4"
        : "p-4 border-b border-blue-gray-50";

    return (
        <tr key={role.id}>
            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                    {role.id}
                </Typography>
            </td>
            <td className={classes}>
                <div className="w-max">
                    <Chip
                        size="sm"
                        variant="ghost"
                        value={role.name}
                        color={
                            role.name === "管理员"
                                ? "green"
                                : role.name === "普通用户"
                                    ? "blue"
                                    : "gray"
                        }
                    />
                </div>
            </td>
            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                    {role.description}
                </Typography>
            </td>
            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                {role.created_at}
                </Typography>
            </td>
            <td className={classes}>
                <Typography variant="small" color="blue-gray" className="font-normal">
                    {role.updated_at}
                </Typography>
            </td>
            <td className={classes}>
                <Tooltip content="Edit Role">
                    <IconButton variant="text">
                        <PencilIcon className="h-4 w-4"/>
                    </IconButton>
                </Tooltip>
            </td>
        </tr>
    );
}

export default RenderRole;