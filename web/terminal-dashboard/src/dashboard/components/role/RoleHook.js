
// 自定义钩子函数，用于获取所有用户
import {useEffect, useState} from "react";
import { getAllRoles } from "../../../api/api";

export const useFetchRoles = () => {
    const [roles, setRoles] = useState([]);

    useEffect(() => {
        getAllRoles()
            .then(data => setRoles(data))
            .catch(error => console.error('Error:', error));
    }, []);
    return roles;
};

