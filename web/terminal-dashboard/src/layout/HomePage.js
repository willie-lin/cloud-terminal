import UserList from "../dashboard/components/user/UserList";

function HomePage({ email }) {
    return (
        <div className="flex flex-col items-center justify-center h-screen bg-blue-50">
           {/*<TopNavbar />*/}
            {/*<Route path="/userinfo" render={() => <UserInfo email={email}/>}/>*/}

            <UserList email={email}/>
        </div>
    );
}

export default HomePage;
