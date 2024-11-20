import UserList from "../dashboard/components/user/UserList";
import {useTheme} from "./ThemeContext";

function HomePage({ email }) {
    return (

        <div className="w-full">
            <div>
                <UserList/>
            </div>
        </div>
    );

    // const { isDarkMode } = useTheme();
    // return (
    //     <div className={`w-full ${isDarkMode ? 'text-white' : 'text-gray-800'}`}>
    //         <h1 className={`text-2xl font-bold mb-6 ${isDarkMode ? 'text-white' : 'text-gray-800'}`}>Dashboard</h1>
    //         <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-white'} rounded-lg shadow p-6`}>
    //             <h2 className={`text-xl font-semibold mb-4 ${isDarkMode ? 'text-white' : 'text-gray-800'}`}>User List</h2>
    //             <UserList />
    //         </div>
    //     </div>
    // );
}

export default HomePage;
