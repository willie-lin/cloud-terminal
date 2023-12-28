import UserInfo from "../dashboard/components/user/UserInfo";
import UserList from "../dashboard/components/user/UserList";
import TopNavbar from "./TopNavbar";

function HomePage({ email }) {
    return (
        <div className="flex flex-col items-center justify-center h-screen bg-blue-50">
           <TopNavbar />
            <UserInfo email={email}/>
            <UserList email={email}/>
        </div>
    );
}

export default HomePage;
