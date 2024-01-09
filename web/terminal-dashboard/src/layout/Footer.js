import {Typography} from "@material-tailwind/react";
import {HeartIcon} from "@heroicons/react/16/solid";

function Footer() {
    return (
        <footer className="py-2">
            <div className="flex w-full flex-wrap items-center justify-center gap-6 px-2 md:justify-between">
                <Typography variant="small" className="font-normal text-inherit">
                    &copy; {2024}, made with{" "}
                    <HeartIcon className="-mt-0.5 inline-block h-3.5 w-3.5 text-red-600"/> by{" "}
                    for a better web.
                </Typography>
            </div>
        </footer>
    )
}
export default Footer;