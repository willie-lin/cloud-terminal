import React from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Breadcrumbs } from "@material-tailwind/react";
import { useTheme } from './ThemeContext';

function BreadCrumbNavigation() {
    const navigate = useNavigate();
    const location = useLocation();
    const { isDarkMode } = useTheme();
    const pathname = location.pathname.split('/').filter((x) => x);

    const linkClass = `transition-colors duration-200 ${
        isDarkMode
            ? 'text-blue-300 hover:text-blue-100'
            : 'text-blue-900 hover:text-blue-700'
    }`;

    const lastLinkClass = `opacity-60 ${
        isDarkMode ? 'text-gray-300' : 'text-gray-600'
    }`;

    return (
        <Breadcrumbs className={isDarkMode ? 'bg-gray-800' : 'bg-white'}>
            <a
                href="/"
                onClick={(e) => {
                    e.preventDefault();
                    navigate("/");
                }}
                className={linkClass}
            >
                Dashboard
            </a>
            {pathname.map((value, index) => {
                const last = index === pathname.length - 1;
                const to = `/${pathname.slice(0, index + 1).join('/')}`;
                return last ? (
                    <a
                        href={to}
                        className={`${linkClass} ${lastLinkClass}`}
                        key={to}
                    >
                        {value}
                    </a>
                ) : (
                    <a
                        href={to}
                        className={linkClass}
                        key={to}
                        onClick={(e) => {
                            e.preventDefault();
                            navigate(to);
                        }}
                    >
                        {value}
                    </a>
                );
            })}
        </Breadcrumbs>
    );
}

export default BreadCrumbNavigation;