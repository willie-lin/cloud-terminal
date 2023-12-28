import {useLocation, useNavigate} from "react-router-dom";
import {Breadcrumbs} from "@material-tailwind/react";
import React from "react";

function BreadCrumbNavigation() {
    const navigate = useNavigate();
    const location = useLocation();
    const pathname = location.pathname.split('/').filter((x) => x);

    return (
        <Breadcrumbs>
            <a href="/" onClick={(e) => {
                e.preventDefault();
                navigate("/");
            }}>Dashboard</a>
            {pathname.map((value, index) => {
                const last = index === pathname.length - 1;
                const to = `/${pathname.slice(0, index + 1).join('/')}`;
                return last ? (
                    <a href={to} className="opacity-60" key={to}>{value}</a>
            ) : (
                    <a href={to} className="opacity-60" key={to} onClick={() => navigate(to)}>
                        {value}
                    </a>
                );
            })}
        </Breadcrumbs>
    );
}

export default BreadCrumbNavigation;