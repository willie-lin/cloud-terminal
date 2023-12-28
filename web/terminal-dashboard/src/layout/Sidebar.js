import {
    Card,
    List,
    ListItem,
    ListItemPrefix,
    ListItemSuffix,
    Chip,
} from "@material-tailwind/react";
import {
    PresentationChartBarIcon,
    ShoppingBagIcon,
    UserCircleIcon,
    Cog6ToothIcon,
    InboxIcon,
    PowerIcon,
} from "@heroicons/react/24/solid";
import {Link} from "react-router-dom";
import React from "react";

function Sidebar({ onLogout }) {
    return (
        <Card className="h-[calc(100vh-2rem)] w-full max-w-[20rem] p-4 shadow-xl shadow-blue-gray-900/5">
            <div className="mb-4 p-4">
                <Link to="/" className="-m-1.5 p-1.5">
                    <span className="sr-only">Your Company</span>
                    <img className="h-8 w-auto"
                         src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=600" alt=""/>
                </Link>
            </div>
            <div className="relative">
                <img className="w-10 h-10 rounded"
                     src="../../../assets/soft-ui-flowbite/images/people/profile-picture-5.jpg" alt=""/>
                <span
                    className="absolute bottom-0 left-8 transform translate-y-1/4 w-3.5 h-3.5 bg-green-400 border-2 border-white dark:border-gray-800 rounded-full"></span>
            </div>
            <List className="space-y-2">
                <ListItem className="hover:bg-blue-gray-100 transition-colors duration-200">
                    <ListItemPrefix>
                        <PresentationChartBarIcon className="h-5 w-5"/>
                    </ListItemPrefix>
                    <Link to="/dashboard">Dashboard</Link>
                </ListItem>
                <ListItem>
                    <ListItemPrefix>
                        <ShoppingBagIcon className="h-5 w-5"/>
                    </ListItemPrefix>
                    E-Commerce
                </ListItem>
                <ListItem>
                    <ListItemPrefix>
                        <InboxIcon className="h-5 w-5"/>
                    </ListItemPrefix>
                    Inbox
                    <ListItemSuffix>
                        <Chip value="14" size="sm" variant="ghost" color="blue-gray" className="rounded-full"/>
                    </ListItemSuffix>
                </ListItem>
                <ListItem>
                    <ListItemPrefix>
                        <UserCircleIcon className="h-5 w-5"/>
                    </ListItemPrefix>
                    <Link to="/userinfo">Profile</Link>
                </ListItem>
                <ListItem>
                    <ListItemPrefix>
                        <Cog6ToothIcon className="h-5 w-5"/>
                    </ListItemPrefix>
                    Settings
                </ListItem>
                <ListItem className="hover:bg-blue-gray-100 transition-colors duration-200">
                    <ListItemPrefix>
                        <PowerIcon className="h-5 w-5"/>
                    </ListItemPrefix>
                    <Link to="/login" onClick={onLogout} >Logout</Link>
                </ListItem>
            </List>
        </Card>
    );
}

export default Sidebar;