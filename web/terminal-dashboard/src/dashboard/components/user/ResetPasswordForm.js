import React from "react";
import {Button, Card, CardBody, CardHeader, Input, Typography} from "@material-tailwind/react";

function ResetPasswordForm({ onResetPassword }) {
    return (
        <section className="m-8 flex">
            <div className="w-2/5 h-full hidden lg:block">
                <img src="/img/pattern.png" className="h-full w-full object-cover rounded-3xl" alt="/"/>
            </div>
            <div className="w-full lg:w-3/5 flex flex-col items-center justify-center">
                <div className="text-center">
                    <Typography variant="h2" className="font-bold mb-4">Join Us Today</Typography>
                    <Typography variant="paragraph" color="blue-gray" className="text-lg font-normal">Enter your email
                        and password to register.</Typography>
                </div>
            </div>
            <div className="flex justify-center items-center h-screen">
                <Card className="w-1/2">
                    <CardHeader color="lightBlue">
                        <Typography color="black" style={{marginBottom: '0'}} className="font-bold text-center">
                            Reset User Password
                        </Typography>
                    </CardHeader>
                    <CardBody className="px-4 py-8">
                        <form onSubmit="">
                        {/*<form onSubmit={handleSubmit}>*/}
                            <div className="mb-1 flex flex-col gap-6">
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    Email
                                </Typography>
                                <Input
                                    size="lg"
                                    type="email"
                                    color="lightBlue"
                                    outline={true}
                                    placeholder="name@mail.com"
                                    // value={email}
                                    // onChange={ handleEmailChange }
                                    // error={!!emailError}
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
                                    // onChange={handleChange}
                                />
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    Password
                                </Typography>
                                <Input
                                    type="password"
                                    color="lightBlue"
                                    size="regular"
                                    outline={true}
                                    placeholder="Password"
                                    // value={password}
                                    // onChange={(e) => setPassword(e.target.value)}
                                />
                                <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                                    Confirm Password
                                </Typography>
                                <Input
                                    type="password"
                                    color="lightBlue"
                                    size="regular"
                                    outline={true}
                                    placeholder="Confirm Password"
                                    // value={confirmPassword}
                                    // onChange={(e) => setConfirmPassword(e.target.value)}
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
        </section>
    )
}

export default ResetPasswordForm;