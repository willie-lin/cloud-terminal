import {
    List,
    ListItem,
    ListItemPrefix,
    ListItemSuffix,
    Chip, Typography, Accordion, AccordionHeader, AccordionBody, Alert
} from "@material-tailwind/react";
import {
    PresentationChartBarIcon,
    UserCircleIcon,
    Cog6ToothIcon,
    InboxIcon,
    PowerIcon,} from "@heroicons/react/24/solid";
import {Link} from "react-router-dom";
import React from "react";
import {
    ChevronDownIcon,
    ChevronRightIcon,
    CubeTransparentIcon,
    UserIcon
} from "@heroicons/react/16/solid";
import { useTheme } from './ThemeContext';




function Sidebar({ email, onLogout }) {
    const [open, setOpen] = React.useState(0);
    const [openAlert, setOpenAlert] = React.useState(true);
    const { isDarkMode } = useTheme();


    const handleOpen = (value) => {
        setOpen(open === value ? 0 : value);
    };


    return (
        <div className={`h-full w-full max-w-[20rem] p-4 ${isDarkMode ? 'bg-gray-800 text-white' : 'bg-white text-gray-800'}`}>
            <div className="mb-2 flex items-center gap-4 p-4">
                {/*<Link to="/" className="-m-1.5 p-1.5">*/}
                <Link to="/">
                    <img src="https://docs.material-tailwind.com/img/logo-ct-dark.png" alt="brand" className="h-8 w-8"/>
                </Link>
                <Typography variant="h5" className={`${isDarkMode ? 'text-white' : 'text-blue-gray-900'}`}>
                    <Link to="/">
                        Sidebar
                    </Link>
                </Typography>
            </div>
            {/*<div className="p-2">*/}
            {/*    <Input icon={<MagnifyingGlassIcon className="h-5 w-5"/>} label="Search"/>*/}
            {/*</div>*/}
            <List>
                <Accordion
                    open={open === 1}
                    icon={
                        <ChevronDownIcon
                            strokeWidth={2.5}
                            className={`mx-auto h-4 w-4 transition-transform ${open === 1 ? "rotate-180" : ""}`}
                        />
                    }
                >
                    <ListItem
                        className="p-0"

                        selected={open === 1}>
                        <AccordionHeader onClick={() => handleOpen(1)} className="border-b-0 p-3">
                            <ListItemPrefix>
                                <PresentationChartBarIcon className={`h-5 w-5 ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-500'}`}/>
                            </ListItemPrefix>
                            <Typography
                                className={`${isDarkMode ? 'text-white mr-auto font-normal' : 'text-blue-gray-900 mr-auto font-normal'}`}>
                                Dashboard
                            </Typography>
                        </AccordionHeader>
                    </ListItem>
                    <AccordionBody className="py-1">
                        <List className="p-0">
                            <ListItem
                                className={`${isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-blue-gray-50'} transition-colors duration-200`}>
                                <ListItemPrefix>
                                    <ChevronRightIcon strokeWidth={3} className="h-3 w-5"/>
                                </ListItemPrefix>
                                Analytics
                            </ListItem>
                            <ListItem
                                className={`${isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-blue-gray-50'} transition-colors duration-200`}>
                                <ListItemPrefix>
                                    <ChevronRightIcon strokeWidth={3} className="h-3 w-5"/>
                                </ListItemPrefix>
                                Reporting
                            </ListItem>
                            <ListItem
                                className={`${isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-blue-gray-50'} transition-colors duration-200`}>
                                <ListItemPrefix>
                                    <ChevronRightIcon strokeWidth={3} className="h-3 w-5"/>
                                </ListItemPrefix>
                                Projects
                            </ListItem>
                        </List>
                    </AccordionBody>
                </Accordion>
                <Accordion
                    open={open === 2}
                    icon={
                        <ChevronDownIcon
                            strokeWidth={2.5}
                            className={`mx-auto h-4 w-4 transition-transform ${open === 2 ? "rotate-180" : ""}`}
                        />
                    }
                >
                    {/*{user.role === 'Super Admin' && (*/}
                    <ListItem className="p-0" selected={open === 2}>
                        <AccordionHeader onClick={() => handleOpen(2)} className="border-b-0 p-3">
                            <ListItemPrefix>
                                <UserIcon className={`h-5 w-5 ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-500'}`}/>
                            </ListItemPrefix>
                            <Typography
                                className={`${isDarkMode ? 'text-white mr-auto font-normal' : 'text-blue-gray-900 mr-auto font-normal'}`}>
                                Management
                            </Typography>
                        </AccordionHeader>
                    </ListItem>
                    <AccordionBody className="py-1">
                        <List className="p-0">
                            <ListItem
                                className={`${isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-blue-gray-50'} transition-colors duration-200`}>
                                <Link to="/users" style={{display: 'flex', alignItems: 'center', width: '100%'}}>
                                    <ListItemPrefix>
                                        <ChevronRightIcon strokeWidth={3} className="h-3 w-5"/>
                                    </ListItemPrefix>
                                    Users
                                </Link>
                            </ListItem>
                            <ListItem
                                className={`${isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-blue-gray-50'} transition-colors duration-200`}>
                                <Link to="/roles" style={{display: 'flex', alignItems: 'center', width: '100%'}}>
                                    <ListItemPrefix>
                                        <ChevronRightIcon strokeWidth={3} className="h-3 w-5"/>
                                    </ListItemPrefix>
                                    Roles
                                </Link>
                            </ListItem>
                            <ListItem
                                className={`${isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-blue-gray-50'} transition-colors duration-200`}>
                                <Link to="/permissions" style={{display: 'flex', alignItems: 'center', width: '100%'}}>
                                    <ListItemPrefix>
                                        <ChevronRightIcon strokeWidth={3} className="h-3 w-5"/>
                                    </ListItemPrefix>
                                    Permissions
                                </Link>
                            </ListItem>
                            <ListItem
                                className={`${isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-blue-gray-50'} transition-colors duration-200`}>
                                <Link to="/authorization"
                                      style={{display: 'flex', alignItems: 'center', width: '100%'}}>
                                    <ListItemPrefix>
                                        <ChevronRightIcon strokeWidth={3} className="h-3 w-5"/>
                                    </ListItemPrefix>
                                    Authorization
                                </Link>
                            </ListItem>
                        </List>
                    </AccordionBody>
                </Accordion>
                {/*</ListItem>*/}
                {/*)}*/}
                {/*    {user.role === 'Super Admin' && (*/}
                {/*        <ListItem className="p-0" selected={open === 2}>*/}
                {/*            <AccordionHeader onClick={() => handleOpen(2)} className="border-b-0 p-3">*/}
                {/*                <ListItemPrefix>*/}
                {/*                    <UserIcon className={`h-5 w-5 ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-500'}`}/>*/}
                {/*                </ListItemPrefix>*/}
                {/*                <Typography className={`${isDarkMode ? 'text-white mr-auto font-normal' : 'text-blue-gray-900 mr-auto font-normal'}`}>
                {/*                    User Management*/}
                {/*                </Typography>*/}
                {/*            </AccordionHeader>*/}
                {/*            <AccordionBody className="py-1">*/}
                {/*                <List className="p-0">*/}
                {/*                    <ListItem>*/}
                {/*                        <ListItemPrefix>*/}
                {/*                            <ChevronRightIcon strokeWidth={3} className="h-3 w-5"/>*/}
                {/*                        </ListItemPrefix>*/}
                {/*                        User Permissions*/}
                {/*                    </ListItem>*/}
                {/*                    <ListItem>*/}
                {/*                        <ListItemPrefix>*/}
                {/*                            <ChevronRightIcon strokeWidth={3} className="h-3 w-5"/>*/}
                {/*                        </ListItemPrefix>*/}
                {/*                        User Roles*/}
                {/*                    </ListItem>*/}
                {/*                </List>*/}
                {/*            </AccordionBody>*/}
                {/*        </ListItem>*/}
                {/*    )}*/}

                <hr className={`my-2 ${isDarkMode ? 'border-gray-600' : 'border-blue-gray-50'}`}/>
                <ListItem
                    className={`${isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-blue-gray-50'} transition-colors duration-200`}
                >
                    <ListItemPrefix>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5}
                             stroke="currentColor" className={`h-5 w-5 ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-500'}`}>
                            <path strokeLinecap="round" strokeLinejoin="round"
                                  d="m6.75 7.5 3 2.25-3 2.25m4.5 0h3m-9 8.25h13.5A2.25 2.25 0 0 0 21 18V6a2.25 2.25 0 0 0-2.25-2.25H5.25A2.25 2.25 0 0 0 3 6v12a2.25 2.25 0 0 0 2.25 2.25Z"/>
                        </svg>
                    </ListItemPrefix>
                    Cloud Terminal
                </ListItem>

                <ListItem
                    className={`${isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-blue-gray-50'} transition-colors duration-200`}
                >
                    <ListItemPrefix>
                        <InboxIcon className={`h-5 w-5 ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-500'}`}/>
                    </ListItemPrefix>
                    Inbox
                    <ListItemSuffix>
                        <Chip value="14" size="sm" variant="ghost" color="blue-gray" className="rounded-full"/>
                    </ListItemSuffix>
                </ListItem>

                <ListItem
                    className={`${isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-blue-gray-50'} transition-colors duration-200`}
                >
                    <Link to="/userinfo" style={{display: 'flex', alignItems: 'center', width: '100%'}}>
                        <ListItemPrefix>
                            <UserCircleIcon className={`h-5 w-5 ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-500'}`}/>
                        </ListItemPrefix>
                        Profile</Link>
                </ListItem>

                <ListItem
                    className={`${isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-blue-gray-50'} transition-colors duration-200`}
                >
                    <Link to="/open-user-2fa" style={{display: 'flex', alignItems: 'center', width: '100%'}}>
                        <ListItemPrefix>
                            <Cog6ToothIcon className={`h-5 w-5 ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-500'}`}/>
                        </ListItemPrefix>
                        Settings
                    </Link>
                </ListItem>

                <ListItem
                    // className="hover:bg-blue-gray-100 transition-colors duration-200"
                    className={`${isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-blue-gray-50'} transition-colors duration-200`}
                >
                    <Link to="/login" style={{display: 'flex', alignItems: 'center', width: '100%'}} onClick={onLogout}>
                        <ListItemPrefix>
                            <PowerIcon className={`h-5 w-5 ${isDarkMode ? 'text-gray-400' : 'text-blue-gray-500'}`}/>
                        </ListItemPrefix>
                        Logout
                    </Link>
                </ListItem>
            </List>
            <Alert open={openAlert} className={`mt-auto ${isDarkMode ? 'bg-gray-700 text-white' : 'bg-white text-gray-800'}`} onClose={() => setOpenAlert(false)}>
                <CubeTransparentIcon className="mb-4 h-12 w-12"/>
                <Typography variant="h6" className={`${isDarkMode ? 'text-white mr-auto font-normal' : 'text-blue-gray-900 mr-auto font-normal'}`}>
                    Upgrade to PRO
                </Typography>
                <Typography variant="small" className={`${isDarkMode ? 'text-white mr-auto font-normal' : 'text-blue-gray-900 mr-auto font-normal'}`}>
                    Upgrade to Material Tailwind PRO and get even more components, plugins, advanced features
                    and premium.
                </Typography>
                <div className="mt-4 flex gap-3">
                    <Typography
                        as="a"
                        href="#"
                        variant="small"
                        // className="font-medium opacity-80"
                        className={`${isDarkMode ? 'text-white mr-auto font-medium opacity-80' : 'text-blue-gray-900 mr-auto font-medium opacity-80'}`}
                        onClick={() => setOpenAlert(false)}
                    >
                        Dismiss
                    </Typography>
                    <Typography as="a" href="#" variant="small"
                                // className="font-medium"
                                className={`${isDarkMode ? 'text-white mr-auto font-medium opacity-80' : 'text-blue-gray-900 mr-auto font-medium opacity-80'}`}
                    >
                        Upgrade Now
                    </Typography>
                </div>
            </Alert>
        </div>


    );
}

export default Sidebar;