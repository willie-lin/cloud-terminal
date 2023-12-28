import UserInfo from "../dashboard/components/UserInfo";
import UserList from "../dashboard/components/UserList";
import {Input, ListItemPrefix} from "@material-tailwind/react";
import {Cog6ToothIcon} from "@heroicons/react/24/solid";
import BreadCrumbNavigation from "./BreadCrumbNavigation";
import React from "react";

function HomePage({ email }) {
    return (
        <div className="flex flex-col items-center justify-center h-screen bg-blue-50">
            <div className="flex justify-end items-center w-full p-4">
                <BreadCrumbNavigation />
                <div className="flex items-center">
                <div className="group relative">
                    <Input
                        type="email"
                        placeholder="Search"
                        className="focus:!border-t-gray-900 group-hover:border-2 group-hover:!border-gray-900"
                        labelProps={{
                            className: "hidden",
                        }}
                        readOnly
                    />
                    <div className="absolute top-[calc(50%-1px)] right-2.5 -translate-y-2/4">
                        <kbd
                            className="rounded border border-blue-gray-100 bg-white px-1 pt-px pb-0 text-xs font-medium text-gray-900 shadow shadow-black/5">
                            <span className="mr-0.5 inline-block translate-y-[1.5px] text-base">⌘</span>
                            K
                        </kbd>
                    </div>
                </div>
                <ListItemPrefix>
                    <Cog6ToothIcon className="h-5 w-5"/>
                </ListItemPrefix>
                </div>
            </div>
            <UserInfo email={email}/>
            <UserList email={email}/>
        </div>
    );
}

export default HomePage;
