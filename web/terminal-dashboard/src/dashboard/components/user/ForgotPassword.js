import {Button, Card, CardBody, CardHeader, Input, Typography} from "@material-tailwind/react";
import React from "react";


function ForgotPassword(email) {
    return (
        <div className="flex justify-center items-center h-screen">
            <Card className="w-1/2">
                <CardHeader color="lightBlue">
                    <Typography color="black" style={{marginBottom: '0'}} className="font-bold text-center">
                        修改用户信息
                    </Typography>
                </CardHeader>

                <CardBody className="px-4 py-8">
                    <form onSubmit={handleSubmit}>
                        <label>
                            <img
                                className="w-24 h-24 rounded-full object-cover mb-6 mx-auto cursor-pointer"
                                src={preview || newInfo?.avatar || "https://i1.pngguru.com/preview/137/834/449/cartoon-cartoon-character-avatar-drawing-film-ecommerce-facial-expression-png-clipart.jpg"}
                                alt="Avatar Upload"
                                style={{width: '120px', height: '120px'}}
                            />
                            <input
                                type="file"
                                className="hidden"
                                multiple={false}
                                accept="image/*"
                                onChange={handleFileChange} // 当文件改变时，调用 handleFileChange 函数
                                ref={fileInputRef} // 添加 ref
                            />
                        </label>
                        {/*<div className="mb-10">*/}
                        {/*    <Input*/}
                        {/*        type="text"*/}
                        {/*        color="lightBlue"*/}
                        {/*        size="regular"*/}
                        {/*        outline={true}*/}
                        {/*        placeholder="头像URL"*/}
                        {/*        onChange={handleChange}*/}
                        {/*    />*/}
                        {/*</div>*/}

                        <div className="mb-1 flex flex-col gap-6">
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Nickname
                            </Typography>
                            <Input
                                type="text"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                placeholder="昵称"
                                onChange={handleChange}
                            />
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Email
                            </Typography>
                            <Input
                                type="email"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                placeholder="邮箱"
                                onChange={handleChange}
                            />
                            <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                Bio
                            </Typography>
                            <Input
                                type="text"
                                color="lightBlue"
                                size="regular"
                                outline={true}
                                placeholder="个人简介"
                                onChange={handleChange}
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
        </div>
    )
}