import {Input, ListItemPrefix} from "@material-tailwind/react";
import BreadCrumbNavigation from "./BreadCrumbNavigation";
import {Cog6ToothIcon} from "@heroicons/react/24/solid";
import { useTheme } from './ThemeContext';
import { MoonIcon, SunIcon } from '@heroicons/react/24/solid';

function TopNavbar() {
    // // 添加一个状态来追踪是否启用了夜间模式
    // const [isDarkMode, setIsDarkMode] = useState(false);
    //
    // // 创建一个函数来处理按钮的点击事件
    // const handleDarkModeToggle = () => {
    //     setIsDarkMode(!isDarkMode);
    //     if (!isDarkMode) {
    //         document.documentElement.classList.add('dark');
    //     } else {
    //         document.documentElement.classList.remove('dark');
    //     }
    // };
    const { isDarkMode, toggleDarkMode } = useTheme();
    return (
        // <div className="flex justify-start items-center w-full p-2 ">
        <div className="flex justify-between items-center w-full p-4 overflow-x-auto bg-white dark:bg-gray-800">
            <BreadCrumbNavigation/>
            {/*<div className="group relative ml-auto">*/}
            <div className="flex items-center space-x-4">
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
                <button
                    onClick={toggleDarkMode}
                    className="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700"
                >
                    {isDarkMode ? (
                        <SunIcon className="h-5 w-5 text-yellow-500"/>
                    ) : (
                        <MoonIcon className="h-5 w-5 text-gray-500"/>
                    )}
                </button>
                <ListItemPrefix>
                    <Cog6ToothIcon className="h-5 w-5"/>
                </ListItemPrefix>
            </div>
        </div>
    );
}

export default TopNavbar;
