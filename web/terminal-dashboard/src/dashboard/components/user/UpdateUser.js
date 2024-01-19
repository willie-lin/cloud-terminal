import React, { useState } from 'react';
import {
    Button,
    Card,
    CardBody,
    CardHeader,
    Input,
    Textarea,
    Typography
} from '@material-tailwind/react';

function UpdateUser({ onUpdateUser }) {
    const [nickname, setNickname] = useState(null)
    const [phone, setPhone] = useState(null)
    const [bio, setBio] = useState(null)

    function handleSubmit() {

    }

    return (
        <Card className="w-96">
            <CardHeader variant="gradient" color="gray" className="mb-4 grid h-20 place-items-center">
                <Typography variant="h4" color="white">
                        Update User
                    </Typography>
                </CardHeader>
                <CardBody className="px-4 py-4">
                    <form onSubmit={handleSubmit}>
                        <div className="mb-1 flex flex-col gap-6">
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Nickname
                            </Typography>
                            <Input
                                variant="outlined"
                                label="Nickname"
                                type="nickname"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                // placeholder="Nickname"
                                value={nickname}
                                name="nickname"  // 添加name属性
                                // onChange={handleChange}
                            />
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Phone
                            </Typography>
                            <Input
                                variant="outlined"
                                label="Phone"
                                type="phone"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                // placeholder="Phone"
                                value={phone}
                                name="phone"  // 添加name属性
                                // onChange={handleChange}
                            />
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Bio
                            </Typography>
                            <Textarea
                                variant="outlined"
                                label="Bio"
                                type="text"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                // className={`border-${inputError ? 'red-500' : 'blue-500'}`}
                                // placeholder="Bio"
                                value={bio}
                                name="bio"  // 添加name属性
                                // onChange={handleChange}
                            />
                        </div>
                        <Button fullWidth
                                type="submit"
                                color="lightBlue"
                                buttonType="filled"
                                size="regular"
                                rounded={false}
                                block={false}
                                iconOnly={false}
                                ripple="light"
                                className="mt-6" // 添加边距
                        >
                            Submit
                        </Button>
                    </form>
                </CardBody>
            </Card>
    );
}

export default UpdateUser;
