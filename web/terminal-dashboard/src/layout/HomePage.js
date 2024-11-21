import UserList from "../dashboard/components/user/UserList";
import {useTheme} from "./ThemeContext";

function HomePage({ email }) {

    const { isDarkMode } = useTheme();
    return (
        <div className={`w-full ${isDarkMode ? 'text-white' : 'text-gray-800'}`}>
                <UserList />
        </div>
    );
}

export default HomePage;
