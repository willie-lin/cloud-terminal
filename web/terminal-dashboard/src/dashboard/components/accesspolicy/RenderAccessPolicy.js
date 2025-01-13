import React, {useContext, useState} from "react";
import {AuthContext} from "../../../App";
import {
    Button,
    Chip,
    Dialog,
    DialogBody,
    DialogFooter,
    DialogHeader,
    Typography
} from "@material-tailwind/react";
import {TrashIcon} from "@heroicons/react/16/solid";


function RenderAccessPolicy({ accessPolicy, isLast }) {

    const { currentUser } = useContext(AuthContext);
    // 判断当前用户是否具有删除权限
    // const canDelete = currentUser?.roleName === 'Admin' || currentUser?.roleName === 'SuperAdmin'


    const classes = isLast ? "p-4" : "p-4 border-b border-blue-gray-50";

    const [deleteAccessPolicy, setDeleteAccessPolicy] = useState(null);
    const [isDeleteAccessPolicyOpen, setIsDeleteAccessPolicyOpen] = useState(false);
    const [isStatementsOpen, setIsStatementsOpen] = useState(false);
    const [isDialogOpen, setIsDialogOpen] = useState(false);


    const handleDeleteAccessPolicy = () => {
        setIsDeleteAccessPolicyOpen(false);
    };

    function openDeleteAccessPolicy(user) {
        setDeleteAccessPolicy(user)
        setIsDeleteAccessPolicyOpen(true)
    }

    function closeDeleteAccessPolicy() {
        setIsDeleteAccessPolicyOpen(false);
    }

    const toggleStatements = () => {
        setIsStatementsOpen(!isStatementsOpen);
    };

    const handleDialogOpen = () => {
        setIsDialogOpen(true);
    };

    const handleDialogClose = () => {
        setIsDialogOpen(false);
    };

    const getColor = roleName => {
        if (roleName.toLowerCase().includes("admin")) return "green";
        return "gray";
    }

    return (
        <>
            <tr key={accessPolicy.id}>
                <td className={classes}>
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {accessPolicy.id}
                    </Typography>
                </td>
                <td className={classes}>
                    {/*<Typography variant="small" color="blue-gray" className="font-normal">*/}
                    {/*    {accessPolicy.name}*/}
                    {/*</Typography>*/}
                    <div className="w-max">
                        <Chip
                            size="sm"
                            variant="ghost"
                            value={accessPolicy.name}
                            color={getColor(accessPolicy.name)
                            }
                        />
                    </div>
                </td>
                <td className={classes}>
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {accessPolicy.description || 'No description'}
                    </Typography>
                </td>
                <td className={classes}>
                    <Button variant="outlined" color="blue" onClick={handleDialogOpen}>
                        View Statements
                    </Button>
                </td>
                <td className={classes}>
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {new Date(accessPolicy.created_at).toLocaleString()}
                    </Typography>
                </td>
                <td className={classes}>
                    <Typography variant="small" color="blue-gray" className="font-normal">
                        {new Date(accessPolicy.updated_at).toLocaleString()}
                    </Typography>
                </td>
            </tr>

            <Dialog open={isDialogOpen} handler={handleDialogClose} size="xl">
                <DialogHeader>Statements</DialogHeader>
                <DialogBody divider>
                    {accessPolicy.statements.map((statement, index) => (
                        <div key={index} className="mt-2 p-2 bg-gray-100 rounded">
                            <Typography variant="h6" color="blue-gray" className="font-normal">
                                <strong>Effect:</strong> {statement.Effect}
                            </Typography>
                            <Typography variant="h6" color="blue-gray" className="font-normal">
                                <strong>Actions:</strong> {statement.Action.join(', ')}
                            </Typography>
                            <Typography variant="h6" color="blue-gray" className="font-normal">
                                <strong>Resources:</strong> {statement.Resource.join(', ')}
                            </Typography>
                        </div>
                    ))}
                </DialogBody>
                <DialogFooter>
                    <Button variant="outlined" color="blue" onClick={handleDialogClose}>
                        Close
                    </Button>
                </DialogFooter>
            </Dialog>
        </>
    );
}

export default RenderAccessPolicy;