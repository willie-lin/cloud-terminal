import {Button, Card, CardBody, CardFooter, Dialog, Input, Typography} from "@material-tailwind/react";

function DeleteDialogPermission() {
    const [open, setOpen] = React.useState(false);
    const handleOpen = () => setOpen((cur) => !cur);

    function setName(value) {

    }

    return (
        <>
            <Button onClick={handleOpen}>Delete Permission</Button>
            <Dialog size="xs" open={open} handler={handleOpen} className="bg-transparent shadow-none">
                <Card className="mx-auto w-full max-w-[24rem]">
                    <CardBody className="flex flex-col gap-4">
                        <Typography variant="h4" color="blue-gray">
                            Delete Permission
                        </Typography>
                        <Typography
                            className="mb-3 font-normal"
                            variant="paragraph"
                            color="gray"
                        >
                            Delete the data for the Permission.
                        </Typography>
                        <Typography variant="small" color="blue-gray" className="-mb-3 font-medium">
                            Name
                        </Typography>
                        <Input
                            variant="outlined"
                            label="Disabled"
                            disabled
                            type="name"
                            color="lightBlue"
                            size="regular"
                            outline={true}
                            value={name}
                            name="username"  // 添加name属性
                            onChange={e => setName(e.target.value)}
                        />
                    </CardBody>
                    <CardFooter className="pt-0">
                        <Button variant="gradient" onClick={handleOpen} fullWidth>
                            Submit
                        </Button>
                    </CardFooter>
                </Card>
            </Dialog>
            </>
    );

}



