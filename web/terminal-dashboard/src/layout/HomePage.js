import UserList from "../dashboard/components/user/UserList";
import DefaultPagination from "../dashboard/components/user/Pagination";

function HomePage({ email }) {
    return (
        <div className="flex flex-col items-center justify-center h-screen bg-blue-50">
            <UserList email={email}/>
        </div>
    );
}

export default HomePage;
