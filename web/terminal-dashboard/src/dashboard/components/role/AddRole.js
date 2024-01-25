import {Button, Card, CardBody, CardHeader, IconButton, Input, Tooltip, Typography} from "@material-tailwind/react";
import React, {useState} from "react";
import {TrashIcon} from "@heroicons/react/16/solid";


function AddRole({ onAddUser, onClose }) {

    const [roles, setRoles] = useState([{ key: '', displayName: '', group: '' }]);

    const handleRoleChange = (index, field, value) => {
        const newRoles = [...roles];
        newRoles[index][field] = value;
        setRoles(newRoles);
    };

    const handleAddRole = () => {
        setRoles([...roles, { key: '', displayName: '', group: '' }]);
    };

    const handleRemoveRole = (index) => {
        const newRoles = [...roles];
        newRoles.splice(index, 1);
        setRoles(newRoles);
    };


    const handleSubmit = (event) => {
        event.preventDefault();
        console.log(roles);
    };

    return (
        <Card className="w-256">
            <CardBody className="px-4 py-8">
                <div className="flex justify-between items-center mb-4">
                    <Typography variant="h4" color="gray">
                        Create Role
                    </Typography>
                    <Button color="gray" buttonType="link"   onClick={onClose}>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5}
                             stroke="currentColor" className="w-4 h-4">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12"/>
                        </svg>
                    </Button>
                </div>
                <Typography variant="body2" color="blueGray" className="mb-4">
                    Enter the data for the new role.
                </Typography>
                <form onSubmit={handleSubmit}>
                    {roles.map((role, index) => (
                        <div className="mb-1 flex gap-6 items-end" key={index}>
                            <div className="flex flex-col">
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    {/*Key*/}
                                </Typography>
                                <Input
                                    variant="outlined"
                                    label="KEY"
                                    type="KEY"
                                    color="lightBlue"
                                    size="regular"
                                    outline={true}
                                    value={role.key}
                                    name="KEY"
                                    onChange={(e) => handleRoleChange(index, 'key', e.target.value)}
                                    className="mb-4"
                                />
                            </div>
                            <div className="flex flex-col">
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    {/*Display Name*/}
                                </Typography>
                                <Input
                                    variant="outlined"
                                    label="DisplayName"
                                    type="DisplayName"
                                    color="lightBlue"
                                    size="regular"
                                    outline={true}
                                    value={role.displayName}
                                    name="DisplayName"
                                    onChange={(e) => handleRoleChange(index, 'displayName', e.target.value)}
                                    className="mb-4"
                                />
                            </div>
                            <Tooltip content="Delete User" placement="top">
                                <IconButton variant="text"
                                            color="red"
                                            buttonType="filled"
                                            size="regular"
                                            rounded={false}
                                            block={false}
                                            iconOnly={true}
                                            ripple="light"
                                    // onClick={() => openDeleteUser(user)}
                                            onClick={() => handleRemoveRole(index)}
                                >
                                    <TrashIcon className="h-4 w-4"/>
                                </IconButton>
                            </Tooltip>
                        </div>
                    ))}
                    <Button
                        color="lightBlue"
                        buttonType="filled"
                        size="regular"
                        rounded={false}
                        block={false}
                        iconOnly={false}
                        ripple="light"
                        onClick={handleAddRole}
                    >
                        + Add additional role
                    </Button>
                    <Button fullWidth
                            type="submit"
                            color="lightBlue"
                            buttonType="filled"
                            size="regular"
                            rounded={false}
                            block={false}
                            iconOnly={false}
                            ripple="light"
                            className="mt-6"
                    >
                        Save
                    </Button>
                </form>
            </CardBody>
        </Card>






        // const handleSubmit = () => {
        //
        // };
        // return (
        //
        //     <Card className="w-96">
        //         <CardHeader variant="gradient" color="gray" className="mb-4 grid h-20 place-items-center">
        //             <Typography variant="h4" color="white">
        //                 Add Role
        //             </Typography>
        //         </CardHeader>
        //         <CardBody className="px-4 py-8">
        //             <form onSubmit={handleSubmit}>
        //                 <div className="mb-1 flex flex-col gap-6">
        //                     <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
    //                         Key
    //                     </Typography>
    //                     <Input
    //                         variant="outlined"
    //                         label="KEY"
    //                         type="KEY"
    //                         color="lightBlue"
    //                         size="regular"
    //                         outline={true}
    //                         // value={KEY}
    //                         name="KEY"  // 添加name属性
    //                         // onChange={handleEmailChange}
    //                         className="mb-4" // 添加边距
    //                         // error={!!emailError}
    //                     />
    //                     <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
    //                         Display Name
    //                     </Typography>
    //                     <Input
    //                         variant="outlined"
    //                         label="DisplayName"
    //                         type="DisplayName"
    //                         color="lightBlue"
    //                         size="regular"
    //                         outline={true}
    //                         // value={DisplayName}
    //                         name="DisplayName"  // 添加name属性
    //                         // onChange={handleChange}
    //                         className="mb-4" // 添加边距
    //                         // onChange={(e) => setPassword(e.target.value)}
    //                     />
    //                 </div>
    //                 <Button fullWidth
    //                         type="submit"
    //                         color="lightBlue"
    //                         buttonType="filled"
    //                         size="regular"
    //                         rounded={false}
    //                         block={false}
    //                         iconOnly={false}
    //                         ripple="light"
    //                         className="mt-6" // 添加边距
    //                 >
    //                     Submit
    //                 </Button>
    //             </form>
    //         </CardBody>
    //     </Card>
    )
}

export default AddRole