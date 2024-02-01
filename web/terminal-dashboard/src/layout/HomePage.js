import UserList from "../dashboard/components/user/UserList";

function HomePage({ email }) {
    return (
        <div className="flex flex-col items-center justify-center bg-blue-50">
            <UserList email={email}/>
        </div>
    );
}

export default HomePage;
