import React from "react";
import {Button, Card, CardBody, CardHeader, Input, Typography} from "@material-tailwind/react";

function ResetPasswordForm({ onResetPassword }) {
    return (
        <div className="flex justify-center items-center h-screen">
            <Card className="w-1/2">
                <div className="text-center">
                    <Typography variant="h2" className="font-bold mb-4">Reset Password</Typography>
                    <Typography variant="paragraph" color="blue-gray" className="text-lg font-normal">Enter your email reset password to Sign In.</Typography>
                </div>
                <form onSubmit="">
                {/*<form onSubmit={handleSubmit}>*/}
                <div className="mb-1 flex flex-col gap-6">
                    <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                        Email
                    </Typography>
                    <Input
                        type="email"
                        color="lightBlue"
                        size="regular"
                        outline={true}
                        placeholder="Email"
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
            </Card>
        </div>
    )
}

export default ResetPasswordForm;